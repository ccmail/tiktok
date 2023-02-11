package common

type BaseResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// feed related responses

type FeedUser struct {
	Id             uint   `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    uint   `json:"follow_count,omitempty"`
	FollowerCount  uint   `json:"follower_count,omitempty"`
	IsFollow       bool   `json:"is_follow,omitempty"`
	TotalFavorited uint   `json:"total_favorited"`
	FavoriteCount  uint   `json:"favorite_count"`
}

type FeedVideo struct {
	Id            uint     `json:"id,omitempty"`
	Author        FeedUser `json:"author,omitempty"`
	PlayUrl       string   `json:"play_url,omitempty"`
	CoverUrl      string   `json:"cover_url,omitempty"`
	FavoriteCount uint     `json:"favorite_count,omitempty"`
	CommentCount  uint     `json:"comment_count,omitempty"`
	IsFavorite    bool     `json:"is_favorite,omitempty"`
	Title         string   `json:"title,omitempty"`
}

type FeedResponse struct {
	BaseResponse
	VideoList []FeedVideo `json:"video_list,omitempty"`
	NextTime  uint        `json:"next_time,omitempty"`
}

type FeedNoVideoResponse struct {
	BaseResponse
	NextTime uint `json:"next_time"`
}

// video related responses
type ReturnAuthor struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
type ReturnMyself struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
}

type ReturnVideo struct {
	VideoId uint `json:"id"`
	//根据https://www.apifox.cn/apidoc/shared-09d88f32-0b6c-4157-9d07-a36d32d7a75c/api-50707525中的接口信息,对id映射的json进行了更改, 由video_id->id
	Author        ReturnAuthor `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount uint         `json:"favorite_count"`
	CommentCount  uint         `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}
type ReturnVideo2 struct {
	VideoId       uint         `json:"video_id"`
	Author        ReturnMyself `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount uint         `json:"favorite_count"`
	CommentCount  uint         `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

type VideoListResponse struct {
	BaseResponse
	VideoList []ReturnVideo `json:"video_list"`
}
type VideoListResponse2 struct {
	BaseResponse
	VideoList []ReturnVideo2 `json:"video_list"`
}
