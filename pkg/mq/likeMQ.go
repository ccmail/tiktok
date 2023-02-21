package mq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"tiktok/config"
	"tiktok/mapper/db"
)

type LikeMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

type LikeStruct struct {
	VideoID uint
	UserID  uint
	IsLike  bool
}

// NewLikeRabbitMQ 获取likeMQ的对应队列。
func NewLikeRabbitMQ(queueName string) *LikeMQ {
	likeMQ := &LikeMQ{
		RabbitMQ:  *RMQ,
		queueName: queueName,
	}
	cha, err := likeMQ.conn.Channel()
	likeMQ.channel = cha
	if err != nil {
		log.Panicln("获取通道失败", err)
	}
	return likeMQ
}

// Publish like操作的发布配置。
func (l *LikeMQ) Publish(msg LikeStruct) {
	marshal, err := json.Marshal(msg)
	if err != nil {
		log.Panicln("序列化时失败")
	}
	_, err = l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	err1 := l.channel.Publish(l.exchange, l.queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        marshal,
	})
	if err1 != nil {
		panic(err)
	}

}

// Consumer like关系的消费逻辑。
func (l *LikeMQ) Consumer() {

	_, err := l.channel.QueueDeclare(l.queueName, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	messages, err1 := l.channel.Consume(
		l.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err1 != nil {
		panic(err1)
	}

	forever := make(chan bool)
	switch l.queueName {
	case "like_add":
		//点赞消费队列
		go l.consumerLikeAdd(messages)
	case "like_del":
		//取消赞消费队列
		go l.consumerLikeUpdate(messages)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

//consumerLikeAdd 赞关系添加的消费方式。
func (l *LikeMQ) consumerLikeAdd(messages <-chan amqp.Delivery) {
	for d := range messages {
		var temp LikeStruct
		err := json.Unmarshal(d.Body, &temp)
		if err != nil {
			log.Panicln("解析信息时失败", err, "信息为:", d.Body, "===.")
		}
		err = db.CreateLikeRecord(temp.UserID, temp.VideoID, temp.IsLike)
		if err != nil {
			log.Panicln("点赞关系再插入到数据库时失败")
		}
	}
}

//consumerLikeUpdate 赞关系删除的消费方式。
func (l *LikeMQ) consumerLikeUpdate(messages <-chan amqp.Delivery) {
	for d := range messages {
		var temp LikeStruct
		err := json.Unmarshal(d.Body, &temp)
		if err != nil {
			log.Panicln("解析信息时失败", err, "信息为:", d.Body, "===.")
		}
		db.UpdateLikeRecord(temp.UserID, temp.VideoID, temp.IsLike)
	}
}

var RmqLikeAdd *LikeMQ
var RmqLikeUpdate *LikeMQ

// InitLikeRabbitMQ 初始化rabbitMQ连接。
func InitLikeRabbitMQ() {
	RmqLikeAdd = NewLikeRabbitMQ(config.MQLikeAdd)
	go RmqLikeAdd.Consumer()

	RmqLikeUpdate = NewLikeRabbitMQ(config.MQLikeUpdate)
	go RmqLikeUpdate.Consumer()
}
