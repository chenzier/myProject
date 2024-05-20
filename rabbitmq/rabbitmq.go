package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"product/datamodels"
	"product/fronted/middlerware"
	"product/services"
	"sync"
)

// url格式 amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://imoocuser:123456@127.0.0.1:5673/imooc"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	//队列名称
	QueueName string
	Exchange  string
	Key       string
	Mqurl     string
	sync.Mutex
}

// MQURL已经定义
// 创建RabbitMQ结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{QueueName: queueName, Exchange: exchange, Key: key, Mqurl: MQURL}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败")
	return rabbitmq
}

func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf(message, err)
		panic(fmt.Sprintf(message, err))
	}
}

// // 创建简单模式下的RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// 使用简单模式下的RabbitMQ实例：实现 1.生产 2.消费
// 生产
func (r *RabbitMQ) PublishSimple(message string) error {
	//1.申请队列,如果队列不存在就创建，存在则跳过
	//保证队列存在，消息能发送到队列中
	r.Lock()
	defer r.Unlock()
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		return err
	}
	//2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		//如果为true,根据exchange类型和routkey规则,如果无法找到符合条件的队列，那么会把发送的消息返回给发送者
		false,
		//如果为true,当exchange发送消息到队列后，发现队列上没有绑定消费者，会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "test/plain",
			Body:        []byte(message),
		})
	fmt.Println("写入消息")
	return nil
}

//	func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
//		//1.申请队列,如果队列不存在就创建，存在则跳过
//		//保证队列存在，消息能发送到队列中
//		_, err := r.channel.QueueDeclare(
//			r.QueueName,
//			//是否持久化
//			false,
//			//是否为自动删除
//			false,
//			//是否具有排他性
//			false,
//			//是否阻塞
//			false,
//			//额外属性
//			nil,
//		)
//		if err != nil {
//			fmt.Println(err)
//		}
//
//		//消费者流控，防止爆库
//		r.channel.Qos(1, //当前消费者一次能接受最大消息数量
//			0,     //服务器传递的最大容量(以八位字节为单位)
//			false, //如果设置为true 对channel可用
//		)
//		//2.接收消息
//		msgs, err := r.channel.Consume(
//			r.QueueName,
//			//用来区分多个消费者
//			"",
//			//是否自动应答
//			false,
//			//是否具有排他性
//			false,
//			//设置为ture则不能将同一个connection中发送的消息传递给这个connection中的消费者
//			false,
//			//队列消费是否阻塞
//			false,
//			//其他参数
//			nil,
//		)
//		if err != nil {
//			fmt.Println(err)
//		}
//		forever := make(chan bool)
//		//3.启用协程处理消息
//		go func() {
//			for d := range msgs {
//				log.Println("Received a message:%s", string(d.Body))
//				message := &datamodels.Message{}
//				err := json.Unmarshal([]byte(d.Body), message)
//				if err != nil {
//					//fmt.Println(3, err)
//				}
//				//插入订单
//				_, err = orderService.InsertOrderByMessage(message)
//				if err != nil {
//					//fmt.Println(2, err)
//				}
//
//				//扣除商品数量
//				err = productService.SubNumberOne(message.ProductID)
//				if err != nil {
//					fmt.Println(1, err)
//				}
//
//				//设为false表示 确认当前消息发送给rabbitmq
//				//true表示ack所有未确认的消息，一般用于批量
//				d.Ack(false)
//			}
//		}()
//		log.Printf("[*] Waiting for messages,To exit press CTRL+C")
//		<-forever
//	}
//
// // 订阅模式下创建RabbitMQ实例
//
//	func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
//		//创建RabbitMQ实例
//		rabbitmq := NewRabbitMQ("", exchangeName, "")
//		var err error
//		//获取connection
//		rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
//		rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
//		//获取channel
//		rabbitmq.channel, err = rabbitmq.conn.Channel()
//		rabbitmq.failOnErr(err, "failed to open a channel")
//		return rabbitmq
//	}
//
// // 订阅模式生产
//
//	func (r *RabbitMQ) PublishPub(message string) {
//		//1. 尝试创建交换机
//		err := r.channel.ExchangeDeclare(
//			r.Exchange,
//			"fanout", //广播模式
//			true,
//			false,
//			false,
//			false,
//			nil,
//		)
//		r.failOnErr(err, "Failed to declare an excha"+"nge")
//		//2.发送消息
//		err = r.channel.Publish(
//			r.Exchange,
//			"",
//			false,
//			false,
//			amqp.Publishing{
//				ContentType: "text/plain",
//				Body:        []byte(message),
//			})
//	}
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService, poolNumber int) {
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Qos(1, 0, false)

	msgs, err := r.channel.Consume(
		r.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	pool := middlerware.NewPool(poolNumber) // 设置协程池大小为3

	go func() {
		for d := range msgs {

			d.Ack(false) //立即确认

			//修改逻辑，如果消费失败，再尝试x次，如果x次还失败，存入死信队列
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}
			task := middlerware.Task{
				ID: 1, // 使用 DeliveryTag 作为任务 ID
				Job: func() {
					log.Printf("Received a message: %s", string(d.Body))
					_, err := orderService.InsertOrderByMessage(message)
					if err != nil {
						fmt.Println(err)
					}
					err = productService.SubNumberOne(message.ProductID)
					if err != nil {
						fmt.Println(err)
					}

				},
			}
			pool.AddTask(task)
		}
	}()

	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	select {}
}

// 订阅模式消费
func (r *RabbitMQ) RecieveSub() {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", //广播模式
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"ange")
	//2.试探性创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true, //排他，设置为true
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name, //上面创建的队列
		"",
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	fmt.Println("CTRL+C退出")
	<-forever
}

// 路由模式
// 创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 路由模型发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", //交换机——direct模式
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key, //要设置routingKey
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 路由模式消费
func (r *RabbitMQ) RecieveRouting() {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct", //direct模式
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"ange")
	//2.试探性创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true, //排他，设置为true
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name, //上面创建的队列
		r.Key,  //绑定Key
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	fmt.Println("CTRL+C退出")
	<-forever
}

// Topic模式
// 创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

// 话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic", //交换机——topic模式
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"nge")
	//2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key, //要设置routingKey
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

// 话题模式消费
// 要注意key,规则
// 其中"*"用于匹配一个单词,"#"用于匹配多个单词(可以是零个)
// 匹配imooc.* 表示匹配一个单词
func (r *RabbitMQ) RecieveTopic() {
	//1. 尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic", //话题模式
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an excha"+"ange")
	//2.试探性创建队列
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true, //排他，设置为true
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name, //上面创建的队列
		r.Key,  //绑定Key
		r.Exchange,
		false,
		nil,
	)
	//消费消息
	messages, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	forever := make(chan bool)
	go func() {
		for d := range messages {
			log.Printf("Received a message:%s", d.Body)
		}
	}()
	fmt.Println("CTRL+C退出")
	<-forever
}
