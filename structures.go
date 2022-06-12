package main

type mostViewedVideo struct {
	ViewCount 	int		`json:"view_count"`
	Title 		string	`json:"title"`
}

type filteredData struct {
	TotalVideoViews 		int				`json:"total_video_views"`
	AverageVideoViews 		int				`json:"average_video_views"`
	TotalVideoDuration 		string			`json:"total_video_duration"`
	AverageViewsPerMinute 	float64			`json:"average_views_per_minute"`
	VideoCount				int				`json:"video_count"`
	MostViewedVideo 		mostViewedVideo	`json:"most_viewed_video"`
}

type videoData struct {
	Id 				string	`json:"id"`
	Stream_id 		string	`json:"stream_id"`
	User_login 		string	`json:"user_login"`
	User_name 		string	`json:"user_name"`
	Title 			string	`json:"title"`
	Description 	string	`json:"description"`
	Created_at 		string	`json:"created_at"`
	Published_at 	string	`json:"published_at"`
	Url 			string	`json:"url"`
	Thumbnail_url 	string	`json:"thumbnail_url"`
	View_count 		int		`json:"view_count"`
	Language 		string	`json:"language"`
	Duration 		string	`json:"duration"`
}

type videoDataResponse struct {
	Data []videoData	`json:"data"`
}

type Error struct {
	Error string `json:"error"`
}

type twitchUser struct {
	Id 				string 	`json:"id"`
	Login 			string	`json:"login"`
	Display_name 	string	`json:"display_name"`
	View_count 		int		`json:"view_count"`
}

type twitchUserResponse struct {
	Data []twitchUser `json:"data"`
}

type twitchAuthRequest struct {
	Access_token 	string `json:"access_token"`
	Expires_in 		int		`json:"expires_in"`
	Token_type 		string	`json:"token_type"`
}
