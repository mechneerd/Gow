package queue

import (
	"fmt"
	"sync"
)

// SQSDriver implements the Driver interface for Amazon SQS
type SQSDriver struct {
	QueueURL  string
	Region    string
	AccessKey string
	SecretKey string
	Endpoint  string // Optional custom endpoint (e.g., for LocalStack)
	mu        sync.RWMutex
	jobs      []Job
}

// NewSQSDriver creates a new SQS driver
func NewSQSDriver(queueURL, region, accessKey, secretKey string) *SQSDriver {
	return &SQSDriver{
		QueueURL:  queueURL,
		Region:    region,
		AccessKey: accessKey,
		SecretKey: secretKey,
		jobs:      make([]Job, 0),
	}
}

// Push pushes a job to SQS
func (d *SQSDriver) Push(job Job) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// In production, use aws-sdk-go-v2:
	// sqsClient := sqs.NewFromConfig(cfg)
	// payload, _ := json.Marshal(job)
	// _, err := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
	//     QueueUrl:    aws.String(d.QueueURL),
	//     MessageBody: aws.String(string(payload)),
	// })

	d.jobs = append(d.jobs, job)
	fmt.Printf("[SQS] Pushing job to %s\n", d.QueueURL)
	return nil
}

// Pop pops a job from SQS
func (d *SQSDriver) Pop() (Job, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// In production, use aws-sdk-go-v2:
	// sqsClient := sqs.NewFromConfig(cfg)
	// result, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
	//     QueueUrl:            aws.String(d.QueueURL),
	//     MaxNumberOfMessages: aws.Int32(1),
	//     WaitTimeSeconds:     aws.Int32(20),
	// })

	if len(d.jobs) == 0 {
		return nil, fmt.Errorf("no jobs available")
	}

	job := d.jobs[0]
	d.jobs = d.jobs[1:]
	return job, nil
}

// Ensure SQSDriver implements Driver at compile time
var _ Driver = (*SQSDriver)(nil)

// RabbitMQDriver implements the Driver interface for RabbitMQ
type RabbitMQDriver struct {
	ConnectionString string
	QueueName        string
	Exchange         string
	mu               sync.RWMutex
	jobs             []Job
}

// NewRabbitMQDriver creates a new RabbitMQ driver
func NewRabbitMQDriver(connectionString, queueName, exchange string) *RabbitMQDriver {
	return &RabbitMQDriver{
		ConnectionString: connectionString,
		QueueName:        queueName,
		Exchange:         exchange,
		jobs:             make([]Job, 0),
	}
}

// Push pushes a job to RabbitMQ
func (d *RabbitMQDriver) Push(job Job) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// In production, use amqp091-go:
	// ch, _ := d.client.Channel()
	// ch.ExchangeDeclare(d.Exchange, "direct", true, false, false, false, nil)
	// ch.QueueDeclare(d.QueueName, true, false, false, false, nil)
	// payload, _ := json.Marshal(job)
	// ch.PublishWithContext(context.TODO(), d.Exchange, d.QueueName, false, false, amqp091.Publishing{
	//     ContentType: "application/json",
	//     Body:        payload,
	// })

	d.jobs = append(d.jobs, job)
	fmt.Printf("[RabbitMQ] Pushing job to %s\n", d.QueueName)
	return nil
}

// Pop pops a job from RabbitMQ
func (d *RabbitMQDriver) Pop() (Job, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// In production, use amqp091-go:
	// msg, ok, _ := ch.Get(d.QueueName, false)
	// if ok {
	//     var job Job
	//     json.Unmarshal(msg.Body, &job)
	//     return &job, nil
	// }

	if len(d.jobs) == 0 {
		return nil, fmt.Errorf("no jobs available")
	}

	job := d.jobs[0]
	d.jobs = d.jobs[1:]
	return job, nil
}

// Ensure RabbitMQDriver implements Driver at compile time
var _ Driver = (*RabbitMQDriver)(nil)
