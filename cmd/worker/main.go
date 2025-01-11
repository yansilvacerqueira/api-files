package worker

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/yansilvacerqueira/api-files/internal/bucket"
	"github.com/yansilvacerqueira/api-files/internal/queue"
)

// TODO: improving the architecture of this code
func main() {
	// RabbitMQ queue configuration
	rabbitConfig := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBIT_URL"),
		QueueName: os.Getenv("RABBIT_TOPIC_NAME"),
		Timeout:   time.Second * 30,
	}

	queueClient, err := queue.NewQueue(queue.RabbitMQ, rabbitConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the queue: %v", err)
	}

	msgChannel := make(chan queue.QueueMessage)
	queueClient.ReceiveMessage(msgChannel)

	// AWS bucket configuration
	bucketConfig := bucket.AWSconfig{
		Config: aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), ""),
		},
		BucketDownload: "drive-raw",
		BucketUpload:   "drive-compact",
	}

	awsBucket, err := bucket.NewAWSBucket(bucketConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the bucket: %v", err)
	}

	// Processing messages from the queue
	for message := range msgChannel {
		sourcePath := fmt.Sprintf("%s/%s", message.Path, message.Filename)
		destinationPath := fmt.Sprintf("%d/%s", message.ID, message.Filename)

		file, err := awsBucket.Download(sourcePath, destinationPath)
		if err != nil {
			log.Printf("Error downloading file: %v", err)
			continue
		}

		fileContent, err := io.ReadAll(file)
		if err != nil {
			log.Printf("Error reading file content: %v", err)
			continue
		}

		// Compressing the file
		var compressedBuffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressedBuffer)

		if _, err = gzipWriter.Write(fileContent); err != nil {
			log.Printf("Error compressing file: %v", err)
			continue
		}

		if err = gzipWriter.Close(); err != nil {
			log.Printf("Error closing gzip writer: %v", err)
			continue
		}

		gzipReader, err := gzip.NewReader(&compressedBuffer)
		if err != nil {
			log.Printf("Error creating gzip reader: %v", err)
			continue
		}

		// Uploading the compressed file
		if err = awsBucket.Upload(gzipReader, sourcePath); err != nil {
			log.Printf("Error uploading compressed file: %v", err)
			continue
		}

		// Removing the local file after processing
		if err = os.Remove(destinationPath); err != nil {
			log.Printf("Error removing temporary file: %v", err)
			continue
		}
	}
}
