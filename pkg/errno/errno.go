package errno

import "errors"

var (
	ErrorNullUsername      = errors.New("用户名为空")
	ErrorUsernameExtend    = errors.New("用户名长度不符合规范")
	ErrorNullPassword      = errors.New("密码为空")
	ErrorPasswordLength    = errors.New("密码长度不符合规范")
	ErrorRedundantUsername = errors.New("用户名已存在")
	ErrorFullPossibility   = errors.New("账号或密码出错")
	ErrorWrongPassword     = errors.New("密码错误")
	ErrorNullVideo         = errors.New("视频不存在")

	//ErrorRelationExit      = errors.New("关注已存在")
	//ErrorNullRelation      = errors.New("关注不存在")
	//ErrorWeakPassword      = errors.New("密码强度过低，应同时包含大小写字母和数字且长度为8~32位")
	//ErrorNullPointer       = errors.New("空指针异常")
)
