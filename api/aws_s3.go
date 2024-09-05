package api

import (
	"bytes"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/emmaahmads/summafy/util"
)

type awsConfig struct {
	s3_bucket string
	region    string
	creds     []string
}

func NewAwsConfig(s3_bucket string, region string, creds ...string) *awsConfig {
	return &awsConfig{
		s3_bucket: s3_bucket,
		region:    region,
		creds:     creds,
	}
}

func (server *Server) UploadFileToS3(fileDir string, file *os.File) (string, error) {
	sess, err := session.NewSession(server.aws)

	if err != nil {
		return "", err
	}

	filename := strings.Split(file.Name(), "/")[2]
	util.MyGinLogger(filename)
	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(server.s3_bucket),
		Key:    aws.String(filename),
		Body:   file,
		// todo add tags for username
	})

	if err != nil {
		return "", err
	}

	return filename, nil
}

// TODO func (server *Server) DeleteFileFromS3(filename string) error {}

func (server *Server) DownloadFileFromS3(filename string) (string, error) {
	sess, err := session.NewSession(server.aws)

	if err != nil {
		return "", err
	}
	util.MyGinLogger(filename)
	svc := s3.New(sess)
	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(server.s3_bucket),
		Key:    aws.String(filename),
	})
	dst := "./download/" + filename
	if err != nil {
		return "", err
	}
	defer obj.Body.Close()

	// read the contents of the object
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj.Body)
	if err != nil {
		return "", err
	}
	contents := buf.String()

	// save the contents to a file
	err = os.WriteFile(dst, []byte(contents), 0644)
	if err != nil {
		return "", err
	}

	return dst, nil
}
