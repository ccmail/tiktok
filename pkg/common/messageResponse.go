package common

type MessageResponse struct {
	ID         uint   `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

type MessageListResponse struct {
	BaseResponse
	MessageResponseList []MessageResponse `json:"message_list"`
}
