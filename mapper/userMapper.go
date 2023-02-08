package mapper

// User 暂时的字段,后期按照提供的app详情可能会需要完善字段
type User struct {
	//用户的底层ID, 除用户名外唯一标识, 在其他功能中, 尽量使用ID来标志用户
	ID uint64
	//这里存储抖音号, 只允许小写字母,数字, 例如"douyin1234", 在数据表中唯一, 考虑到更改的可能性, 不能将其设置为主键
	Account string `gorm:"unique;not null;"`
	//userName这里应该存储昵称
	UserName string
	Password string
	//用户头像, 图片存储到图床, 这里存储图床地址
	UserHeadPicURL string
}
