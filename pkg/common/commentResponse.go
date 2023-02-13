package common

// CommentListResponse 评论表的响应结构体
type CommentListResponse struct {
	BaseResponse
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

// CommentActionResponse 评论操作的响应结构体
type CommentActionResponse struct {
	BaseResponse
	Comment CommentResponse `json:"comment,omitempty"`
}

// UserInfo 用户信息的响应结构体
// type UserInfo struct {
// 	ID            uint   `json:"id,omitempty"`
// 	Name          string `json:"name,omitempty"`
// 	FollowCount   uint   `json:"follow_count,omitempty"`
// 	FollowerCount uint   `json:"follower_count,omitempty"`
// 	IsFollow      bool   `json:"is_follow,omitempty"`
// }

type CommenterInfo UserInfoQueryResponse

// CommentResponse 评论信息的响应结构体
// type CommentResponse struct {
// 	ID         uint     `json:"id,omitempty"`
// 	Content    string   `json:"content,omitempty"`
// 	CreateDate string   `json:"create_date,omitempty"`
// 	User       UserInfo `json:"user,omitempty"`
// }

type CommentResponse struct {
	ID         uint          `json:"id"`
	Content    string        `json:"content"`
	CreateDate string        `json:"create_date"`
	User       CommenterInfo `json:"user"`
}
