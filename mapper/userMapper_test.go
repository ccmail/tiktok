package mapper

import (
	"fmt"
	"log"
	"testing"
)

func TestFindMultiUserInfo(t *testing.T) {
	err := InitDBConnectorSupportTest()
	if err != nil {
		log.Panicln("链接数据库失败")
	}
	//var res []model.User
	//tesrStr := []string{"1", "2", "3", "4"}
	//find := DBConn.Model(&model.User{}).Where("id IN ?", tesrStr).Find(&res)
	//if find.Error != nil {
	//	fmt.Println("用字符串查找失败")
	//}
	//for i := range res {
	//	fmt.Println(i)
	//}

	test := []uint{1, 2, 3, 4}
	user, err := FindMultiUserInfo(test)
	if err != nil {
		t.Error()
	}
	for u := range user {
		fmt.Println(u)
	}
	fmt.Println(user)
}

func TestFindUserInfo(t *testing.T) {
	_ = InitDBConnectorSupportTest()
	user, err := FindUserInfo(1)
	if err != nil {
		t.Error()
	}
	fmt.Println(user)
}
