package bucket

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSconfig struct {
	Config         aws.Config
	BucketDownload string
	BucketUpload   string
}

// Represents an AWS session for S3 bucket operations
type AWSSession struct {
	session        *session.Session
	bucketDownload string
	bucketUpload   string
}

// Download method - Downloads a file from S3 bucket to the specified destination
func (awsSession *AWSSession) Download(src string, dest string) (*os.File, error) {
	// Create a file for the destination
	file, err := os.Create(dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer file.Close()

	// Initialize the S3 downloader
	downloader := s3manager.NewDownloader(awsSession.session)

	// Perform the download
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(awsSession.bucketDownload),
		Key:    aws.String(src),
	})
	if err != nil {
		return nil, fmt.Errorf("error downloading file from S3: %v", err)
	}

	return file, nil
}

// Remove (delete) method - Deletes a file from the S3 bucket
func (awsSession *AWSSession) Remove(src string) error {
	// Initialize the S3 service client
	svc := s3.New(awsSession.session)

	// Perform the delete operation
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(awsSession.bucketDownload),
		Key:    aws.String(src),
	})
	if err != nil {
		return fmt.Errorf("error deleting file from S3: %v", err)
	}

	// Wait until the object no longer exists
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(awsSession.bucketDownload),
		Key:    aws.String(src),
	})
	if err != nil {
		return fmt.Errorf("error waiting for object deletion: %v", err)
	}

	return nil
}

// Upload method - Uploads a file to the S3 bucket
func (awsSession *AWSSession) Upload(file io.Reader, key string) error {
	// Initialize the S3 uploader
	uploader := s3manager.NewUploader(awsSession.session)

	// Perform the upload operation
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awsSession.bucketUpload),
		Key:    aws.String(key),
		Body:   file, // Ensure the file is uploaded
	})
	if err != nil {
		return fmt.Errorf("error uploading file to S3: %v", err)
	}

	return nil
}

// Function to create a new AWS session
// Handles AWS session initialization and configuration
func newAWSSession(cfg AWSconfig) (*AWSSession, error) {
	// Create a new AWS session
	sess, err := session.NewSession(&cfg.Config)
	if err != nil {
		return nil, fmt.Errorf("error creating AWS session: %v", err)
	}

	return &AWSSession{
		session:        sess,
		bucketDownload: cfg.BucketDownload,
		bucketUpload:   cfg.BucketUpload,
	}, nil
}
