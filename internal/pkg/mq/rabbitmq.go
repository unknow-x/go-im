// Package mq
/**
  @author:kk
  @data:2021/9/15
  @note
**/
package mq

import (
	"github.com/streadway/amqp"
	"im_app/config"
	"log"
	"strconv"
)

var RabbitMq *amqp.Connection
var err error

// ConnectMQ 加载mq
func ConnectMQ() *amqp.Connection {
	RabbitMq, err = amqp.Dial("amqp://" + config.Conf.Rabbitmp.User + ":" +
		config.Conf.Rabbitmp.Password + "@" +
		config.Conf.Rabbitmp.Host + ":" +
		strconv.Itoa(config.Conf.Rabbitmp.Port) + "/")
	if err != nil {
		log.Fatal("rabbitmq连接失败")
	}
	//defer RabbitMq.Close()

	return RabbitMq
}
