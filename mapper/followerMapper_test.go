package mapper

import (
	"log"
	"testing"
)

func TestExistFollowRecord(t *testing.T) {
	_ = InitDBConnectorSupportTest()

	_, exist := ExistFollowRecord(1, 2)
	if !exist {
		t.Error()
	}
}
func TestFindMultiConcern(t *testing.T) {
	_ = InitDBConnectorSupportTest()
	list, err := FindMultiConcern(1)
	if err != nil {
		t.Error()
	}
	log.Println(list)
}
