package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	authv1 "k8s.io/api/authorization/v1"
	"log"
	)



func main() {

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
