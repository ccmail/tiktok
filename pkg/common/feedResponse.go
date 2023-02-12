package common

// feed related responses

type FeedUser struct {
	Id            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint   `json:"follow_count,omitempty"`
	FollowerCount uint   `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	//TotalFavorited uint   `json:"total_favorited"`
	//FavoriteCount uint `json:"favorite_count"`
}

//type FeedResponse struct {
//	BaseResponse
//	VideoList []FeedVideo `json:"video_list,omitempty"`
//	NextTime  uint        `json:"next_time,omitempty"`
//}

type FeedNoVideoResponse struct {
	BaseResponse
	NextTime uint `json:"next_time"`
}
