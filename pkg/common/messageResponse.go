package common

type MessageResp struct {
	ID uint `json:"id"`
	//客户端修改了接口, 详见https://bytedance.feishu.cn/docs/doccnKrCsU5Iac6eftnFBdsXTof
	ToUserID   uint   `json:"to_user_id"`
	FromUserID uint   `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

type MessageListBaseResp struct {
	BaseResponse
	MessageResponseList []MessageResp `json:"message_list"`
}
