package rabbitmq

import "log"

type ServiceContext struct {
	Config   Config
	MQ       *RabbitMQ
	Consumer *Consumer
}

func NewServiceContext(c Config) *ServiceContext {
	// 创建 RabbitMQ 实例
	mq, err := NewRabbitMQ(&Config{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建消费者服务
	consumer := NewConsumer(mq)

	// 注册消息处理器
	// articleHandler := NewArticleMessageHandler() // 创建处理器实例
	// consumer.RegisterHandler("article_queue", articleHandler)

	// 启动消费者服务
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}

	return &ServiceContext{
		Config:   c,
		MQ:       mq,
		Consumer: consumer,
	}
}
