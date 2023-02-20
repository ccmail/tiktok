package mapper

import (
	"testing"
)

func TestInitDBConnector(t *testing.T) {
	err := InitDBConnectorSupportTest()
	if err != nil {
		t.Error("连接数据库失败")
	}
	err = createAllTable()
	if err != nil {
		t.Error("数据库建表发生错误")
	}
}
