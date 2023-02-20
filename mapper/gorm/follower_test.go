package gorm

import (
	"log"
	"testing"
	"tiktok/config"
)

func TestExistFollowRecord(t *testing.T) {
	_ = config.InitDBConnectorSupportTest()

	_, exist := ExistFollowRecord(1, 2)
	if !exist {
		t.Error()
	}
}
func TestFindMultiConcern(t *testing.T) {
	_ = config.InitDBConnectorSupportTest()
	list, err := FindMultiConcern(1)
	if err != nil {
		t.Error()
	}
	log.Println(list)
}
