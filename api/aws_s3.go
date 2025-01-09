package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/emmaahmads/summafy/util"
	"github.com/gin-gonic/gin"
)

type AwsConfig struct {
	s3_bucket string
	region    string
	creds     []string
}

type s3ObjectUploaded struct {
	filename string
	summary  string
}

var s3ObjectsNotifiedMap = make(map[string]s3ObjectUploaded)

func NewAwsConfig(s3_bucket string, region string, creds ...string) *AwsConfig {
	return &AwsConfig{
		s3_bucket: s3_bucket,
		region:    region,
		creds:     creds,
	}
}

//	@BasePath	/api/v1

// PingExample godoc
//
//	@Summary	ping example hello
//	@Schemes
//	@Description	do ping
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/example/helloworld [get]
func (server *Server) UploadFileToS3(fileDir string, file *os.File) (s3ObjectUploaded, error) {
	var newFile s3ObjectUploaded
	summary := "N/A"

	util.MyGinLogger("In UploadFileToS3")
	util.MyGinLogger("FileDir:", fileDir)
	util.MyGinLogger("File:", file.Name())
	util.MyGinLogger("Loading AWS config")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		util.MyGinLogger("Error loading AWS config:", err.Error())
		return newFile, err
	}

	util.MyGinLogger("Creating S3 client")
	client := s3.NewFromConfig(cfg)
	filename := strings.Split(file.Name(), "/")[1]

	util.MyGinLogger("Uploading file:", filename)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &server.s3bucket,
		Key:    &filename,
		Body:   file,
		// TODO: Add tags for username
	})

	if err != nil {
		util.MyGinLogger("Error uploading file:", err.Error())
		return newFile, err
	}

	util.MyGinLogger("File uploaded to S3:", filename)

	// go func() {
	// 	util.MyGinLogger("Waiting for notification")
	// 	for {
	// 		if obj, ok := s3ObjectsNotifiedMap[filename]; ok {
	// 			util.MyGinLogger("Summary received:", obj.summary)
	// 			summary = obj.summary
	// 			delete(s3ObjectsNotifiedMap, filename)
	// 			return
	// 		}
	// 		util.MyGinLogger("No notification yet")
	// 		// If no notification, update DB with N/A summary
	// 	}
	// }()

	newFile = s3ObjectUploaded{
		filename: filename,
		summary:  summary,
	}

	return newFile, nil
}

// TODO func (server *Server) DeleteFileFromS3(filename string) error {}

func (server *Server) DownloadFileFromS3(filename string) (string, error) {
	/* sess, err := session.NewSession(server.aws)

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
	} */

	dst := "NA"

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
		summary:  notification.Summary,
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
