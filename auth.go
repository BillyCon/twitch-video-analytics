package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"net/url"
	"encoding/json"
)

// get authentication token from twitch api
func getAuthenticationToken(client_id string, client_secret string) (string, error) {

	authTokenQueryData := url.Values{
		"client_id": {client_id},
		"client_secret": {client_secret},
	}

	authTokenResponse, err := http.PostForm(fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", client_id, client_secret), authTokenQueryData)

	if err != nil {
		fmt.Printf("Error fetching access token: %v\n", err)
		return "", err
	}

	if (authTokenResponse.StatusCode != 200){
		fmt.Printf("Twich token response StatusCode: %v\n", authTokenResponse.StatusCode)
		return "", err
	}

	authTokenResponseAsBytes, err := io.ReadAll(authTokenResponse.Body)

	if err != nil {
		fmt.Printf("Error reading response body from access token request: %v\n", err)
		return "", err
	}

	var authTokenResponseAsJson twitchAuthRequest
	json.Unmarshal(authTokenResponseAsBytes, &authTokenResponseAsJson)

	return authTokenResponseAsJson.Access_token, nil

}

// check that current token is valid and grab a new one if not
func validateAccessToken() (bool, error) {

	validateAccessTokenClient := http.Client{}
	validateAccessTokenRequest, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)

	if err != nil {
		return false, err
	}

	validateAccessTokenRequest.Header.Set("Authorization", fmt.Sprintf("OAuth %s", os.Getenv("ACCESS_TOKEN")))
	validateAccessTokenResponse, err := validateAccessTokenClient.Do(validateAccessTokenRequest)

	if err != nil {
		return false, err
	}

	if validateAccessTokenResponse.StatusCode != 200 {

		fmt.Println("Found non valid token. Attempting to fetch new access token")
		authenticationAccessToken, err := getAuthenticationToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

		if err != nil {
			fmt.Printf("Error getting authentication token: %v", err)
			return false, err
		}

		os.Setenv("ACCESS_TOKEN", authenticationAccessToken)
	}

	return true, nil

}
