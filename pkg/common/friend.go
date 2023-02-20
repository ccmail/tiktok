package common

type FriendListBaseResp struct {
	BaseResponse
	FriendInfo []FriendInfo `json:"user_list"`
}

type FriendInfo struct {
	Message string `json:"message"`
	MsgType int64  `json:"msgType"`
	UserInfoResp
}
