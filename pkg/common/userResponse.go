package common

// user related responses
type UserIdTokenResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserRegisterResponse struct {
	BaseResponse
	UserIdTokenResponse
}

type UserLoginResponse struct {
	BaseResponse
	UserIdTokenResponse
}

type UserInfoQueryResponse struct {
	UserId        uint   `json:"id"`
	Username      string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type UserInfoResponse struct {
	BaseResponse
	UserList UserInfoQueryResponse `json:"user"`
}
