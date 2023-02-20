package redis

import (
	"fmt"
	"testing"
	"tiktok/config"
	"time"
)

func TestSetInfoForRedis(t *testing.T) {
	config.InitDBConnectorSupportTest()
	config.InitRedisConnector()
	fmt.Println(time.Now().Unix())
}
