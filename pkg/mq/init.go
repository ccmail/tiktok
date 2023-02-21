package mq

import (
	"github.com/streadway/amqp"
	"log"
	"tiktok/config"
)

var RMQ *RabbitMQ

type RabbitMQ struct {
	conn      *amqp.Connection
	mqAddress string
}

// InitRabbitMQ 初始化RabbitMQ的连接和通道。
func InitRabbitMQ() {

	RMQ = &RabbitMQ{
		mqAddress: config.RMQAddress,
	}
	dial, err := amqp.Dial(RMQ.mqAddress)
	if err != nil {
		log.Panicln("链接出错")
	}
	log.Println("连接成功!")
	RMQ.conn = dial

}
