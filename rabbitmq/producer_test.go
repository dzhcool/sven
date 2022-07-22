package rabbitmq

import (
	"testing"
)

//  go test -v common/modules/rabbitmq  -test.run Test_producer

func Test_producer(t *testing.T) {
	producer, err := NewProducer("amqp://admin:test@192.168.110.35:5672/")
	if err != nil {
		t.Fatalf("connect rabbitmq err:%s", err.Error())
		return
	}

	err = producer.Producer("ping-exchange", "direct", "ping-queue", "ping-queue", `{"name":"davistest"}`, false)
	if err != nil {
		t.Fatalf("rabbitmq producer err:%s", err.Error())
		return
	}
	t.Logf("success")
}
