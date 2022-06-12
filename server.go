package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	findEnvironmentVariables()

	fmt.Println("Starting server...")
	fmt.Println("Authenticating with Twitch...")

	authenticationAccessToken, err := getAuthenticationToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	if err != nil {
		fmt.Printf("Error fetching access token: %v\n", err)
		os.Exit(1)
	}

	os.Setenv("ACCESS_TOKEN", authenticationAccessToken)

	fmt.Println("Validating access token")
	isTokenValid, err := validateAccessToken()

	if err != nil || !isTokenValid {
		fmt.Printf("Error validating access token: %v\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/videos/", getTwitchVideoAnalytics)

	fmt.Println("Listening on 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
