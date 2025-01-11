package bucket

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

// Download method - yet to be implemented
func (a *AWSSession) Download(src string, dest string) (*os.File, error) {
	return nil, fmt.Errorf("Download not implemented for AWSSession")
}

// Remove (delete) method - yet to be implemented
func (a *AWSSession) Remove(src string) error {
	return fmt.Errorf("Remove not implemented for AWSSession")
}

// Upload method - yet to be implemented
func (a *AWSSession) Upload(file io.Reader, key string) error {
	return fmt.Errorf("Upload not implemented for AWSSession")
}

// Function to create a new AWS session
// Handles AWS session initialization and configuration
func newAWSSession(cfg AWSconfig) (*AWSSession, error) {
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
