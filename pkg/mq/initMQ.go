package mq

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"tiktok/config"
	"tiktok/mapper"
	"tiktok/mapper/db"
)

var (
	wg sync.WaitGroup
	Q  *StreamMQ
)

func InitRedisMQ() {
	Q = NewStreamMQ(mapper.RedisConn, config.MQMaxLen, true)
	go func() {
		err := Q.Consume(context.Background(), config.MQLikeAdd, "g1", "c1", "0", 1, func(msg *Msg) error {
			body := msg.Body
			var temp LikeStruct
			err := json.Unmarshal(body, &temp)
			if err != nil {

				return err
			}
			err = db.CreateLikeRecord(temp.UserID, temp.VideoID, temp.IsLike)
			if err != nil {
				log.Panicln("点赞信息插入数据库失败")
			}
			wg.Done()
			return nil
		})
		if err != nil {
		}
	}()
	go func() {
		err := Q.Consume(context.Background(), config.MQLikeUpdate, "g1", "c1", "0", 1, func(msg *Msg) error {
			body := msg.Body
			var temp LikeStruct
			err := json.Unmarshal(body, &temp)
			if err != nil {
				return err
			}
			db.UpdateLikeRecord(temp.UserID, temp.VideoID, temp.IsLike)
			wg.Done()
			return nil
		})
		if err != nil {
			return
		}
	}()

}
