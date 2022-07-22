package rabbitmq

import (
	log "github.com/dzhcool/sven/zapkit"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	amqpURI string
	conn    *amqp.Connection
	channel *amqp.Channel
}

// 初始化
func NewProducer(amqpURI string) (*Producer, error) {
	producer := &Producer{
		amqpURI: amqpURI,
	}
	if err := producer.newConn(); err != nil {
		log.Errorf("创建RabbitMQ连接失败: %s url: %s", err.Error(), amqpURI)
		return nil, err
	}
	return producer, nil
}

// 创建连接
func (p *Producer) newConn() error {
	conn, err := amqp.Dial(p.amqpURI)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	p.conn = conn
	p.channel = ch
	return nil
}

// 生产
func (p *Producer) Producer(exchange, exchangeType, queueName, routingKey, body string, reliable bool) error {
	var err error

	if _, err = p.channel.QueueDeclarePassive(queueName, true, false, false, true, nil); err != nil {
		if _, err = p.channel.QueueDeclare(queueName, true, false, false, true, nil); err != nil {
			log.Errorf("1: %s", err.Error())
			return err
		}
	}

	if err = p.channel.QueueBind(queueName, routingKey, exchange, true, nil); err != nil {
		log.Errorf("2: %s", err.Error())
		return err
	}

	if err = p.channel.ExchangeDeclarePassive(exchange, exchangeType, true, false, false, false, nil); err != nil {
		if err = p.channel.ExchangeDeclare(
			exchange,
			exchangeType,
			true,  // durable
			false, // delete when complete
			false, // internal
			false, // noWait
			nil,   // arguments
		); err != nil {
			log.Errorf("3: %s", err.Error())
			return err
		}
	}

	if err = p.channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		Headers: amqp.Table{},
		// ContentType: "application/json",
		ContentType: "text/plain",
		Body:        []byte(body),
		// Timestamp:   time.Now(),
	}); err != nil {
		return err
	}

	return nil
}

// 关闭
func (p *Producer) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	if err := p.conn.Close(); err != nil {
		return err
	}
	return nil
}
