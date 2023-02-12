package common

type FavoriteVideo struct { // 从 video 中获取
	ID            uint                  `json:"id,omitempty"`
	Author        UserInfoQueryResponse `json:"author,omitempty"`
	PlayUrl       string                `json:"play_url,omitempty"`
	CoverUrl      string                `json:"cover_url,omitempty"`
	FavoriteCount uint                  `json:"favorite_count,omitempty"`
	CommentCount  uint                  `json:"comment_count,omitempty"`
	IsFavorite    bool                  `json:"is_favorite,omitempty"`
	Title         string                `json:"title,omitempty"`
}

type LikeListResponse struct {
	BaseResponse
	VideoList []FavoriteVideo `json:"video_list,omitempty"`
}
