package mapper

import (
	"fmt"
	"testing"
	"time"
)

func TestSetInfoForRedis(t *testing.T) {
	InitDBConnectorSupportTest()
	InitRedisConnector()
	fmt.Println(time.Now().Unix())
}
