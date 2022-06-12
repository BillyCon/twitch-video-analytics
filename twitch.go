package main

import (
	"fmt"
	"net/http"
	"os"
	"io"
	"encoding/json"
	"strings"
	"time"
)

func getTwitchVideoAnalytics(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%v %v\n", time.Now(), r)

	w.Header().Set("Content-Type", "application/json")

	urlSplit := strings.Split(r.URL.Path, "/videos/")
	twitchUsername := strings.Split(urlSplit[1], "?")[0]
	videoCountAsString := r.URL.Query().Get("videoCount")

	requestIsValid := validateRequest(w, twitchUsername, videoCountAsString)

	if !requestIsValid {
		return
	}

	isTokenValid, err := validateAccessToken()

	if err != nil || !isTokenValid {
		fmt.Printf("Error validating access token and or isTokenValid is false: %v\n", err)
		os.Exit(1)
	}

	users, err := getTwitchIdViaName(twitchUsername)

	if err != nil {
		fmt.Printf("Error running getTwitchIdViaName: %v\n", err)
		sendBackError(w, 500, "server error")
		return
	}

	if len(users) == 0 {
		sendBackError(w, 400, "user does not exist")
		return
	}

	twitchVideosClient := http.Client{}
	twitchVideosRequest, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/videos?user_id=%s&first=%s", users[0].Id, r.URL.Query().Get("videoCount")), nil)

	if err != nil {
		fmt.Printf("Error with forming request to get videos from twitchAPI: %v\n", err)
		sendBackError(w, 500, "server error")
		return
	}

	twitchVideosRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ACCESS_TOKEN")))
	twitchVideosRequest.Header.Set("Client-Id", os.Getenv("CLIENT_ID"))

	twitchVideosResponse, err := twitchVideosClient.Do(twitchVideosRequest)

	if err != nil {
		fmt.Printf("Error with request to get videos from twitchAPI: %v\n", err)
		sendBackError(w, 500, "server error")
		return
	}

	if (twitchVideosResponse.StatusCode != 200) {
		fmt.Printf("Unexpected statusCode for fetching videos from twitchAPI: %v\n", twitchVideosResponse.StatusCode)
		fmt.Printf("Body response: %v\n", twitchVideosResponse.Body)
		sendBackError(w, 500, "server error")
		return
	}

	twitchVideosResponseAsBytes, err := io.ReadAll(twitchVideosResponse.Body)

	if err != nil {
		fmt.Printf("Error reading response body from user login request: %v\n", err)
		sendBackError(w, 500, "server error")
		return
	}

	var twitchVideosResponseAsJson videoDataResponse
	json.Unmarshal(twitchVideosResponseAsBytes, &twitchVideosResponseAsJson)

	reformattedTwitchVideoData, err := reformatTwitchVideoData(twitchVideosResponseAsJson.Data)

	if err != nil {
		fmt.Printf("Error running reformatData(): %v\n", err)
		sendBackError(w, 500, "server error")
		return
	}

	json.NewEncoder(w).Encode(reformattedTwitchVideoData)
}

// get user_id from a twitch username
func getTwitchIdViaName(login string) ([]twitchUser, error) {
	twitchUsersClient := http.Client{}
	twitchUsersRequest, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/users?login=%s", login), nil)

	if err != nil {
		fmt.Printf("Error with forming request to get user login data: %v\n", err)
		return []twitchUser{}, err
	}

	twitchUsersRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ACCESS_TOKEN")))
	twitchUsersRequest.Header.Set("Client-Id", os.Getenv("CLIENT_ID"))

	twitchUsersResponse, err := twitchUsersClient.Do(twitchUsersRequest)

	if err != nil {
		fmt.Printf("Error with request to get user login data: %v\n", err)
		return []twitchUser{}, err
	}

	if (twitchUsersResponse.StatusCode != 200) {
		fmt.Printf("Unexpected statusCode for fetching users from twitchAPI: %v\n", twitchUsersResponse.StatusCode)
		fmt.Printf("Body response: %v\n", twitchUsersResponse.Body)
		return []twitchUser{}, err
	}

	twitchUsersResponseAsBytes, err := io.ReadAll(twitchUsersResponse.Body)

	if err != nil {
		fmt.Printf("Error reading response body from user login request: %v\n", err)
		return []twitchUser{}, err
	}

	var twitchUsersResponseAsJson twitchUserResponse
	json.Unmarshal(twitchUsersResponseAsBytes, &twitchUsersResponseAsJson)

	return twitchUsersResponseAsJson.Data, nil
}

// restructure twitch video data and calculate intended values
func reformatTwitchVideoData(data []videoData) (filteredData, error) {

	//set default values
	totalViews := 0
	totalDurationInSeconds := 0
	topViewedVideo := 0
	avgViews := 0
	averageViewsPerMinute := 0.0

	var topViewedVideoObj mostViewedVideo = mostViewedVideo{
		ViewCount: 0,
		Title: "",
	}

	for twitchVideoDataIndex, _ := range data {
		totalViews += data[twitchVideoDataIndex].View_count
		durationOfVideoInSeconds, err := formatDurationToSecondsInt(data[twitchVideoDataIndex].Duration)

		if err != nil {
			fmt.Printf("Error running formatDurationToSecondsInt: %v\n", err)
			return filteredData{}, err
		}

		totalDurationInSeconds += durationOfVideoInSeconds

		if (topViewedVideo < data[twitchVideoDataIndex].View_count){
			topViewedVideo = data[twitchVideoDataIndex].View_count
			topViewedVideoObj.ViewCount = data[twitchVideoDataIndex].View_count
			topViewedVideoObj.Title = data[twitchVideoDataIndex].Title
		}
	}

	if (len(data) > 0 || totalDurationInSeconds != 0){ // this check is to make sure we do not divide by 0
		avgViews = totalViews / len(data)
		averageViewsPerMinute = float64(totalViews) / float64(totalDurationInSeconds)
	}

	return filteredData{
		TotalVideoViews: 		totalViews,
		AverageVideoViews: 		avgViews,
		TotalVideoDuration: 	formatSecondsToDurationString(totalDurationInSeconds),
		AverageViewsPerMinute:  averageViewsPerMinute,
		VideoCount:				len(data),
		MostViewedVideo: 		topViewedVideoObj,
	}, nil
}
