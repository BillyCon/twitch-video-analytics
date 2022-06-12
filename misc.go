package main

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/joho/godotenv"
	"os"
)

// func getEnvVar(key string, envLocation string) string {

// 	// load .env file
// 	err := godotenv.Load(envLocation)

// 	if err != nil {
// 	  fmt.Println("Error reading environment variables")
// 	  os.Exit(1)
// 	}

// 	return os.Getenv(key)
// }

func findEnvironmentVariables() string {
	err := godotenv.Load(".env")
	if err != nil {
		err := godotenv.Load("../.env/.env")

		if err != nil {
			fmt.Println("Cannot find environment variables")
			os.Exit(1)
		}

		return ".env"
	}

	return ".env/.env"
}

func formatDurationToSecondsInt(duration string) (int, error) { // "1h2m3s" -> 3723

	hours := 0
	minutes := 0
	seconds := 0

	if strings.Contains(duration, "h") && strings.Contains(duration, "m") {
		hoursSplit := strings.Split(duration, "h")
		hoursAsInt, err := strconv.Atoi(hoursSplit[0])

		if err != nil {
			fmt.Println("Error converting string to int")
			return 0, err
		}

		hours = hoursAsInt * 3600

		minutesSplit := strings.Split(hoursSplit[1], "m")
		minutesAsInt, err := strconv.Atoi(minutesSplit[0])

		if err != nil {
			fmt.Println("Error converting string to int")
			return 0, err
		}

		minutes = minutesAsInt * 60

		secondsSplit := strings.Split(minutesSplit[1], "s")
		secondsAsInt, err := strconv.Atoi(secondsSplit[0])
		if err != nil {
			fmt.Println("Error converting string to int")
			return 0, err
		}

		seconds = secondsAsInt

	} else if strings.Contains(duration, "m") && !strings.Contains(duration, "h") {
		minutesSplit := strings.Split(duration, "m")
		minutesAsInt, err := strconv.Atoi(minutesSplit[0])

		if err != nil {
			fmt.Println("Error converting string to int")
			return 0, err
		}

		minutes = minutesAsInt * 60

		secondsSplit := strings.Split(minutesSplit[1], "s")
		secondsAsInt, err := strconv.Atoi(secondsSplit[0])
		if err != nil {
			fmt.Println("Error converting string to int")
			return 0, err
		}

		seconds = secondsAsInt

	} else if !strings.Contains(duration, "m") && strings.Contains(duration, "s") {
		secondsSplit := strings.Split(duration, "s")
		secondsAsInt, err := strconv.Atoi(secondsSplit[0])
		if err != nil {
			fmt.Println("Error converting string to int")
			return 0, err
		}

		seconds = secondsAsInt

	}

	totalDurationInSeconds := hours + minutes + seconds

	return totalDurationInSeconds, nil

}

func formatSecondsToDurationString(seconds int) string { // 3723 -> "1h2m3s"

	hoursFloor := seconds / 3600
	hoursRemainder := seconds % 3600

	minutesFloor := hoursRemainder / 60
	minutesRemainder := hoursRemainder % 60

	secondsFloor := minutesRemainder

	return fmt.Sprintf("%vh%vm%vs", hoursFloor, minutesFloor, secondsFloor)

}

func sendBackError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Error{
		Error: message,
	})
}

// check request data is valid
func validateRequest(w http.ResponseWriter, username string, videoCount string) bool {
	if (username == ""){
		sendBackError(w, 400, "username has cannot be empty")
		return false
	}

	if (videoCount == "" || videoCount == "0"){
		sendBackError(w, 400, "videoCount has to exist and cannot be 0")
		return false
	}

	videoCountAsInt, err := strconv.Atoi(videoCount)

	if err != nil {
		sendBackError(w, 400, "videoCount has to be a number")
		return false
	}

	if (videoCountAsInt > 100){
		sendBackError(w, 400, "videoCount cannot be greater than 100")
		return false
	}

	return true
}
