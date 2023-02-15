package common

type VideoListBaseResp struct {
	BaseResponse
	VideoList []VideoResp `json:"video_list"`
}

// FeedVideoListBaseResp Feed流请求的响应头中多了一个NextTime, 因此单独开一个结构体
type FeedVideoListBaseResp struct {
	BaseResponse
	NextTime  int64       `json:"next_time"`
	VideoList []VideoResp `json:"video_list"`
}

type VideoResp struct {
	ID            uint         `json:"id,omitempty"`
	Author        UserInfoResp `json:"author,omitempty"`
	PlayUrl       string       `json:"play_url,omitempty"`
	CoverUrl      string       `json:"cover_url,omitempty"`
	FavoriteCount uint         `json:"favorite_count,omitempty"`
	CommentCount  uint         `json:"comment_count,omitempty"`
	IsFavorite    bool         `json:"is_favorite,omitempty"`
	Title         string       `json:"title,omitempty"`
}
