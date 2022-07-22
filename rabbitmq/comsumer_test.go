package rabbitmq

import (
	"context"
	"testing"
	"time"

	log "github.com/dzhcool/sven/zapkit"
)

func init() {
	log.ThirdInit("/tmp/test.log", "debug")
}

func Test_consumer(t *testing.T) {
	consumer, err := NewConsumer("amqp://admin:test@192.168.110.35:5672/")
	if err != nil {
		t.Fatalf("connect rabbitmq err:%s", err.Error())
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Consumer(ctx, "ping-exchange", "direct", "ping-queue", "ping-queue", "", consumer_handle)

	time.Sleep(10 * time.Second)
	cancel()
}

// 消费函数
func consumer_handle(messageId string, deliveryTag uint64, body []byte) error {
	log.Infof("收到消息：%s", string(body))
	return nil
}
