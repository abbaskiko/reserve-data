package archive

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"go.uber.org/zap"

	// "github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Archive struct {
	uploader *s3manager.Uploader
	svc      *s3.S3
	awsConf  AWSConfig
	l        *zap.SugaredLogger
}

func (a *S3Archive) UploadFile(bucketName string, awsfolderPath string, filePath string) error {
	file, err := os.Open(filePath)
	defer func() {
		if cErr := file.Close(); cErr != nil {
			a.l.Errorf("File close error: %s", cErr.Error())
		}
	}()
	if err != nil {
		return err
	}
	_, err = a.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filepath.Join(awsfolderPath, getFileNameFromFilePath(filePath))),
		Body:   file,
	})
	return err
}

func (a *S3Archive) RemoveFile(bucketName string, awsfolderPath string, filePath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filepath.Join(awsfolderPath, getFileNameFromFilePath(filePath))),
	}
	_, err := a.svc.DeleteObject(input)
	return err
}

func getFileNameFromFilePath(filePath string) string {
	elems := strings.Split(filePath, "/")
	if len(elems) < 1 {
		return filePath
	}
	fileName := elems[len(elems)-1]
	return fileName
}

func (a *S3Archive) CheckFileIntergrity(bucketName string, awsfolderPath string, filePath string) (bool, error) {
	//get File info
	file, err := os.Open(filePath)
	defer func() {
		if cErr := file.Close(); cErr != nil {
			a.l.Errorf("File close error: %s", cErr.Error())
		}
	}()
	if err != nil {
		return false, err
	}
	fi, err := file.Stat()
	if err != nil {
		return false, err
	}
	//get AWS's file info

	x := s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(filepath.Join(awsfolderPath, getFileNameFromFilePath(filePath))),
	}
	resp, err := a.svc.ListObjects(&x)
	if err != nil {
		return false, err
	}

	for _, item := range resp.Contents {
		remoteFileName := getFileNameFromFilePath(*item.Key)
		localFileName := getFileNameFromFilePath(filePath)
		if (remoteFileName == localFileName) && (*item.Size == fi.Size()) {
			return true, nil
		}
	}
	return false, nil
}

func (a *S3Archive) GetReserveDataBucketName() string {
	return a.awsConf.ExpiredReserveDataBucketName
}

func (a *S3Archive) GetLogBucketName() string {
	return a.awsConf.LogBucketName
}

func NewS3Archive(conf AWSConfig) *S3Archive {
	crdtl := credentials.NewStaticCredentials(conf.AccessKeyID, conf.SecretKey, conf.Token)
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: crdtl,
	}))
	uploader := s3manager.NewUploader(sess)
	svc := s3.New(sess)
	archive := S3Archive{uploader: uploader,
		svc:     svc,
		awsConf: conf,
		l:       zap.S(),
	}

	return &archive
}
