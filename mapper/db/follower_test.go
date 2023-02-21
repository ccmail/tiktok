package db

import (
	"log"
	"testing"
	"tiktok/mapper"
)

func TestExistFollowRecord(t *testing.T) {
	_ = mapper.InitDBConnectorSupportTest()

	_, exist := CheckFollowRecord(1, 2)
	if !exist {
		t.Error()
	}
}
func TestFindMultiConcern(t *testing.T) {
	_ = mapper.InitDBConnectorSupportTest()
	list, err := GetMultiConcern(1)
	if err != nil {
		t.Error()
	}
	log.Println(list)
}
