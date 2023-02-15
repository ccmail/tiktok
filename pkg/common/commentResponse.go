package common

// CommentListBaseResp 评论表的响应结构体
type CommentListBaseResp struct {
	BaseResponse
	CommentList []CommentResp `json:"comment_list,omitempty"`
}

// CommentActionBaseResp 评论操作的响应结构体
type CommentActionBaseResp struct {
	BaseResponse
	Comment CommentResp `json:"comment,omitempty"`
}

type CommentResp struct {
	ID         uint         `json:"id"`
	Content    string       `json:"content"`
	CreateDate string       `json:"create_date"`
	User       UserInfoResp `json:"user"`
}
