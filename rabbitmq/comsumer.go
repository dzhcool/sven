package rabbitmq

import (
	"context"
	"time"

	log "github.com/dzhcool/sven/zapkit"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	amqpURI     string
	conn        *amqp.Connection
	channel     *amqp.Channel
	notifyClose chan *amqp.Error
	closer      chan bool
	payload     *Payload
}

type Payload struct {
	Exchange     string
	ExchangeType string
	QueueName    string
	Key          string
	Tag          string
	Handle       func(string, uint64, []byte) error
}

const (
	reconnectDelay = 5 * time.Second // 重连延迟
	resendDelay    = 5 * time.Second // 重发延迟
	resendTimes    = 3               // 重发次数
)

func NewConsumer(amqpURI string) (*Consumer, error) {
	consumer := &Consumer{
		amqpURI: amqpURI,
	}
	if err := consumer.newConn(); err != nil {
		log.Errorf("创建RabbitMQ连接失败: %s url: %s", err.Error(), amqpURI)
		return nil, err
	}

	return consumer, nil
}

func (p *Consumer) newConn() error {
	conn, err := amqp.Dial(p.amqpURI)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	p.conn = conn
	p.channel = ch

	p.closer = make(chan bool)
	p.notifyClose = make(chan *amqp.Error)
	p.channel.NotifyClose(p.notifyClose)
	return nil
}

// 消费
func (p *Consumer) Consumer(ctx context.Context, exchange, exchangeType, queueName, key, tag string, handle func(string, uint64, []byte) error) {
	p.payload = &Payload{
		exchange,
		exchangeType,
		queueName,
		key,
		tag,
		handle,
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.closer:
			return
		case err := <-p.notifyClose:
			log.Errorf("rabbitmq consumer failed:%s, Retrying...", err)
			for {
				if err := p.newConn(); err != nil {
					time.Sleep(reconnectDelay)
					continue
				}
				break
			}
		default:
			deliveries, err := p.consumer()
			if err != nil {
				log.Errorf("rabbitmq create consumer failed:%s", err)
				time.Sleep(reconnectDelay)
				continue
			}
			for d := range deliveries {
				if err := handle(d.MessageId, d.DeliveryTag, d.Body); err != nil {
					log.Errorf("RabbitMQ consumer err: %s msg: %s", err.Error(), string(d.Body))
					continue
				}
				// 确认消息，必须传递false
				if err := d.Ack(false); err != nil {
					log.Errorf("RabbitMQ ack err: %s msg: %s", err.Error(), string(d.Body))
				}
			}
		}
	}
}

// 消费
func (p *Consumer) consumer() (<-chan amqp.Delivery, error) {
	if err := p.channel.ExchangeDeclare(
		p.payload.Exchange,
		p.payload.ExchangeType,
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return nil, err
	}

	var queue amqp.Queue
	var err error

	if queue, err = p.channel.QueueDeclarePassive(p.payload.QueueName, true, false, false, false, nil); err != nil { // 检查是否创建过
		if queue, err = p.channel.QueueDeclare(
			p.payload.QueueName,
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // noWait
			nil,   // arguments
		); err != nil {
			return nil, err
		}
	}

	if err := p.channel.QueueBind(
		queue.Name,
		p.payload.Key, // bindingKey
		p.payload.Exchange,
		false, // noWait
		nil,   // arguments
	); err != nil {
		return nil, err
	}

	deliveries, err := p.channel.Consume(
		queue.Name,
		p.payload.Tag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	return deliveries, nil
}

// 关闭
func (p *Consumer) Close() error {
	p.closer <- true
	if err := p.channel.Close(); err != nil {
		return err
	}
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}
