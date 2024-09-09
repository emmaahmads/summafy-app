package api

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

type awsConfig struct {
	s3_bucket string
	region    string
	creds     []string
}

type s3ObjectUploaded struct {
	filename string
	summary  string
}

var s3ObjectsNotifiedMap = make(map[string]s3ObjectUploaded)

func NewAwsConfig(s3_bucket string, region string, creds ...string) *awsConfig {
	return &awsConfig{
		s3_bucket: s3_bucket,
		region:    region,
		creds:     creds,
	}
}

func (server *Server) UploadFileToS3(fileDir string, file *os.File) (s3ObjectUploaded, error) {
	var new_file s3ObjectUploaded
	summary := string("N/A")
	sess, err := session.NewSession(server.aws)
	// TODO add timeout for SNS notification
	//timeout := time.NewTimer(5 * time.Second)

	if err != nil {
		return new_file, err
	}

	filename := strings.Split(file.Name(), "/")[2]

	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(server.s3_bucket),
		Key:    aws.String(filename),
		Body:   file,
		// todo add tags for username
	})

	if err != nil {
		return new_file, err
	}

	go func() {
		for {
			if _, ok := s3ObjectsNotifiedMap[filename]; ok {
				if filename == s3ObjectsNotifiedMap[filename].filename {
					summary = s3ObjectsNotifiedMap[filename].summary
				}
				delete(s3ObjectsNotifiedMap, filename)
				return
			}
			// if there is no notification we will just update db with N/A summary
		}
	}()

	new_file = s3ObjectUploaded{
		filename: filename,
		summary:  summary,
	}

	return new_file, nil
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

func (server *Server) HandlerNotification(c *gin.Context) {
	type Notification struct {
		Name    string `json:"filename"`
		Summary string `json:"summary"`
	}
	var notification Notification

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	util.MyGinLogger("In HandlerNotification")
	util.MyGinLogger("Received name:", notification.Name)
	util.MyGinLogger("Received summary:", string(notification.Summary))
	s3ObjectsNotifiedMap[notification.Name] = s3ObjectUploaded{
		filename: notification.Name,
		summary:  notification.Summary, //notification.Summary,
	}
	util.MyGinLogger("Stored notification:", notification.Name, string(notification.Summary))
}

func (server *Server) HandlerDisplayNotifiedObjects(c *gin.Context) {
	util.MyGinLogger("Displaying notified objects:")
	for filename, obj := range s3ObjectsNotifiedMap {
		util.MyGinLogger(fmt.Sprintf("Filename: %s, Summary: %s", filename, string(obj.summary)))
	}
	c.Status(200)
}
