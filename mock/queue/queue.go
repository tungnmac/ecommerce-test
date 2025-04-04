package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-mock/db"
	"go-mock/request"

	"github.com/adjust/rmq"
)

// Define queue connection
var queue rmq.Queue

func InitQueue() {
	// Connect to Redis for queue management
	connection := rmq.OpenConnection("worker", "tcp", "localhost:6379", 0)
	defer connection.Close()

	// Open the queue
	queue = connection.OpenQueue("data_queue")
	defer queue.Close()

	// Start consuming messages (10 prefetch limit, 1 poll duration)
	queue.StartConsuming(10, 5)
}

// Consumer struct implementing rmq.Consumer interface
type DataConsumer struct{}

// Consume processes queue messages
func (consumer *DataConsumer) Consume(delivery rmq.Delivery) {
	var batch []request.Data
	err := json.Unmarshal([]byte(delivery.Payload()), &batch)
	if err != nil {
		log.Println("Queue Unmarshal Error:", err)
		delivery.Reject() // Reject on failure
		return
	}

	// Insert data into PostgreSQL
	db.InsertBatch(batch)

	// Acknowledge successful processing
	delivery.Ack()
}

// Worker function to process queue messages
func StartWorker() {
	consumer := &DataConsumer{} // Create an instance of DataConsumer
	consumerName := fmt.Sprintf("data_worker-%d", time.Now().Unix())

	// Retry mechanism for Redis connection
	for {
		if queue == nil {
			log.Println("Queue not initialized, retrying in 5s...")
			time.Sleep(5 * time.Second)
			continue
		}

		err := queue.AddConsumer(consumerName, consumer)
		if err != "" {
			log.Println("Error adding consumer, retrying in 5s:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Worker started and consuming messages...")
		break // Exit loop when successful
	}
}

// Push batch to queue
func PushToQueue(batch []request.Data) {
	payload, _ := json.Marshal(batch)
	if err := queue.PublishBytes(payload); !err {
		fmt.Println("Queue Publish Error:", err)
	}
}
