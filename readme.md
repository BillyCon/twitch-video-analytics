# Twitch video analytics

Twitch video analytics is a golang service which fetches video data for a particular Twitch.tv user. It receives a user and a number of videos then works out a number of values including total video views and average views per minute.

## API Reference

### GET endpoint

- `/videos/{username}` username cannot be empty.

### Query parameters

- `videoCount` the amount of videos to index. Has to be a number, cannot be `0` and cannot exceed `100`

### Response object

The response is a single object with six key pair values:

`total_video_views` type `number`.  The sum of each videos view count.
`average_video_views` type `number`. The average number of views per video.
`total_video_duration` type `string`. The sum of each videos duration.
`average_views_per_minute` type `number`. Average number of views per minute of video.
`video_count` type `number`. How many videos were indexed (a channel could have only 5 videos when you may be expecting a video count of 10)
`most_viewed_video` type `object`.

**most_viewed_video object**
`view_count` type `number`. Number of views of the video.
`title` type `string`. Title of the video.

### Example query

```
curl --location --request GET 'localhost:8080/videos/shroud?videoCount=10'
```

### Example response

```
{
	"total_video_views":  6420194,
	"average_video_views":  642019,
	"total_video_duration":  "93h14m58s",
	"average_views_per_minute":  19.124909889245693,
	"video_count":  10,
	"most_viewed_video":  {
		"view_count":  973662,
		"title":  "totally not e3 watchparty"
	}
}
```

## Running locally
Run the service locally set the `CLIENT_ID` and `CLIENT_SECRET` for a Twitch app in the `.env` file. Copy and use the `.env.example`.

To launch run:

```
go run .
```

Note this is running for version 1.18 Go

## Building docker container

To build the service run:

```
docker build docker build -t twitch-video-analytics:v1.0 .
```

To run the service run:

```
docker run -p 8081:8081 twich-video-analytics:v1.0
```

## Running in local k8s (minikube)

Start a minikube cluster:
```
minikube start
```

Use this command to allow for minikube local repository:

```
eval $(minikube docker-env) && docker build docker build -t twitch-video-analytics:v1.0 .
```

### Setting the secret

Edit `./k8s/secret.yaml` where `BASE64ENCODEDVARS` should be replaced with the base64 encoded version of the `.env`file.

### Applying the configuration

Apply the deployment, service, secret and ingress to the minikube cluster using:

```
kubectl apply -f ./k8s
```

### Check and test

Check a pod is running `kubectl get pods` and then get the IP of the ingress `kubectl get ingress twitch-video-analytics-ingress` and access the service at the address.