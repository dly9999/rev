package RabitMq

import (
	"cloud/Cloudlog"
	"log"

	//"time"

	//"fmt"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		Cloudlog.Logprint("mqerror:", err)
	}
}
func RevMq(ser string, quname string, inter chan interface{}) {
	//conn, err := amqp.Dial("amqp://admin:admin@172.16.172.98:5672/Admin")
	log.Println("begin")
	conn, err := amqp.Dial(ser)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	defer ch.Close()
	/*err = ch.ExchangeDeclare(
		"TestMQ_Exchange02", // name
		"direct",            // type
		false,               // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		nil,                 // arguments
	)*/
	q, err := ch.QueueDeclare(
		quname, //"TestMQ_Queue02", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")
	/*err = ch.QueueBind(
		q.Name,              // queue name
		"ewrwerw",           // routing key
		"TestMQ_Exchange02", // exchange
		false,
		nil)

	failOnError(err, "Failed to bind a queue")*/

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {

		log.Printf("Received a message ")

		inter <- d.Body
		//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	}
	//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

}
