package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	authv1 "k8s.io/api/authorization/v1"
	"log"
	"os"
	"sync/atomic"
	)

var healthy int32

func main() {

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	router := http.NewServerMux()
	router.Handle("/", index())
	router.Handle("/healthz", healthz())
	router.Handle("/search", search(logger))

	server := &http.Server{
		Addr: ":8000",
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
		
		ctx, cancel := context.WithTimeout(context.Background() 30*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server %v\n", err)
		}
		close(done)
	}()

	logger.Println("Server is ready to handle requests at :8000)
	atomic.StoreInt32(&healthy, 1)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on :8000: %v", err)
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
		if atomic.loadInt32(&healthy) == 1 {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}

type Target struct {
	Username string
}
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	subject := "x"
	resource := authv1.ResourceAttributes{
			Verb: "view",
			Resource: "namespace",
		}
	review := authv1.SubjectAccessReview{
					Spec: authv1.SubjectAccessReviewSpec{
						User:             subject,
						ResourceAttributes: &resource,
						},
					}
	
	_, err = clientset.AuthorizationV1().SubjectAccessReviews().Create(&review)

	if err != nil {
		log.Printf("Err while performing sar review %v", err)
	}

	

}
