package common

// UserIdTokenResp user related responses
type UserIdTokenResp struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserSignBaseResp struct {
	BaseResponse
	UserIdTokenResp
}

type UserInfoResp struct {
	UserID        uint   `json:"id"`
	Username      string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserInfoBaseResp struct {
	BaseResponse
	UserInfo UserInfoResp `json:"user"`
}

type UserInfoListBaseResp struct {
	BaseResponse
	UserList []UserInfoResp `json:"user_list"`
}
