package mapper

type Message struct {
	//消息的ID
	ID               uint64
	UserID           uint64
	UserName         string
	UserHeadPicURL   string
	FriendID         uint64
	FriendName       string
	FriendHeadPicURL string
	MessageText      string
}
