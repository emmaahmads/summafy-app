package api

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	DeleteObjectFunc func(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

func (m *mockS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return m.DeleteObjectFunc(ctx, input, opts...)
}

func TestDeleteFileFromS3(t *testing.T) {
	// Arrange
	server := &Server{s3bucket: "test-bucket"}
	mockClient := &mockS3Client{
		DeleteObjectFunc: func(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
			if input.Bucket == nil || input.Key == nil {
				return nil, errors.New("missing bucket or key")
			}
			if *input.Bucket != "test-bucket" || *input.Key != "test.txt" {
				return nil, errors.New("incorrect bucket or key")
			}
			return &s3.DeleteObjectOutput{}, nil
		},
	}

	// Act
	err := server.deleteFileFromS3WithClient("test.txt", mockClient)

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = server.deleteFileFromS3WithClient("tset.txt", mockClient)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
