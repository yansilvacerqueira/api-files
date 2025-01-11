package bucket

import (
	"io"
	"os"
)

const (
	// Descriptive constant name for clarity
	AWSS3BucketProvider BucketType = iota
)

type BucketType int

// Interface representing a storage bucket provider, allowing for flexibility in provider choice
type StorageProvider interface {
	Upload(io.Reader, string) error
	Download(src string, dest string) (*os.File, error)
	Remove(src string) error
}

type Bucket struct {
	// Use more descriptive variable name
	provider StorageProvider
}

// Function to initialize a new Bucket instance
// Receives a bucket type and configuration
func NewAWSBucket(cfg AWSconfig) (*Bucket, error) {
	awsSession, err := newAWSSession(cfg)
	if err != nil {
		return nil, err
	}

	return &Bucket{
		provider: awsSession,
	}, nil
}

// Upload a file to the bucket using the underlying provider
func (b *Bucket) Upload(file io.Reader, key string) error {
	return b.provider.Upload(file, key)
}

// Download a file from the bucket using the underlying provider
func (b *Bucket) Download(src string, dest string) (*os.File, error) {
	return b.provider.Download(src, dest)
}

// Remove (delete) a file from the bucket using the underlying provider
func (b *Bucket) Delete(src string) error {
	return b.provider.Remove(src)
}
