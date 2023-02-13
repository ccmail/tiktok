package common

type FavoriteVideo struct { // 从 video 中获取
	ID            uint       `json:"id"`
	Author        AuthorInfo `json:"author"`
	PlayUrl       string     `json:"play_url"`
	CoverUrl      string     `json:"cover_url"`
	FavoriteCount uint       `json:"favorite_count"`
	CommentCount  uint       `json:"comment_count"`
	IsFavorite    bool       `json:"is_favorite"`
	Title         string     `json:"title"`
}

type AuthorInfo UserInfoQueryResponse

type LikeListResponse struct {
	BaseResponse
	VideoList []FavoriteVideo `json:"video_list"`
}
