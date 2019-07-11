package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"sync/atomic"
	"net/http"
	"time"
	"encoding/json"
	"strings"
	"fmt"
	"os/signal"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	)

var healthy int32

func main() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	router := http.NewServeMux()
	router.Handle("/", index())
	router.Handle("/healthz", healthz())
	router.Handle("/search", search(logger))

	server := &http.Server{
		Addr: ":8080",
		Handler: (logging(logger)(router)),
		ErrorLog: logger,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,	
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at :8080")
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on :8080: %v", err)
	}

	<-done
	logger.Println("Server stopped")
}

func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})
	
}

func healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

type Target struct {
	Username string
}

func search(logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			target := new(Target)
			err := json.NewDecoder(r.Body).Decode(&target)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			logger.Printf("Getting projects for %s", target.Username)
			array, err := getProjects(logger, target.Username)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			js, err := json.Marshal(array)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	})

}

func getProjects(logger *log.Logger, username string) ([]string, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		log.Printf("Error setting up projects %v", err)
		return []string{}, err
	}

	kc, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Printf("Error creating kc object %v", err)
		return []string{}, err
	}

	roles, err := kc.RbacV1().RoleBindings(metav1.NamespaceAll).List(metav1.ListOptions{})

	if err != nil {
		log.Printf("Error getting projects %v",err)
		return []string{}, nil
	}

	var namespaces []string
	for _, role := range roles.Items {
		if role.RoleRef.Kind == "ClusterRole" && role.RoleRef.Name == "view" {
			for _, subject := range role.Subjects {
				if subject.Kind == "User" && strings.ToUpper(subject.Name) == strings.ToUpper(username) {
					logger.Printf("User allowed to view project %v", role.ObjectMeta.Namespace)
					namespaces = append(namespaces, role.ObjectMeta.Namespace)
				}
			}
		}

	}

	return namespaces, nil


}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
