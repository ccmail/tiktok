package common

type VideoListResponse struct {
	BaseResponse
	VideoList []VideoResp `json:"video_list"`
}

type FeedVideoListResp struct {
	BaseResponse
	NextTime  int64       `json:"next_time"`
	VideoList []VideoResp `json:"video_list"`
}

type VideoResp struct {
	Id            uint         `json:"id,omitempty"`
	Author        UserInfoResp `json:"author,omitempty"`
	PlayUrl       string       `json:"play_url,omitempty"`
	CoverUrl      string       `json:"cover_url,omitempty"`
	FavoriteCount uint         `json:"favorite_count,omitempty"`
	CommentCount  uint         `json:"comment_count,omitempty"`
	IsFavorite    bool         `json:"is_favorite,omitempty"`
	Title         string       `json:"title,omitempty"`
}
