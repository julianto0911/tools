package tools

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3 struct {
	AccessKeyID     string
	SecretAccessKey string
	Token           string
	Region          string
	Endpoint        string
	Bucket          string
	S3TenantID      string
	Path            string
}

func NewS3(cfg S3Config) S3 {
	resp := S3{}
	resp.AccessKeyID = cfg.AccessKeyID
	resp.SecretAccessKey = cfg.SecretAccessKey
	resp.Token = cfg.Token
	resp.Region = cfg.Region
	resp.Endpoint = cfg.Endpoint
	resp.Bucket = cfg.Bucket
	resp.S3TenantID = cfg.S3TenantID
	resp.Path = cfg.Path

	return resp
}

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Token           string
	Region          string
	Endpoint        string
	Bucket          string
	S3TenantID      string
	DoStream        bool
	Path            string
}

func (c *S3) GetSVC() (*s3.S3, error) {
	creds := credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, c.Token)
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}

	cfg := aws.NewConfig().WithRegion(c.Region).WithCredentials(creds).WithEndpoint(c.Endpoint).WithS3ForcePathStyle(true)
	mySession := session.Must(session.NewSession())
	return s3.New(mySession, cfg), nil
}

// func (c *S3) objPath(obj string) string {
// 	return c.Endpoint + "/" + c.S3TenantID + ":" + c.Bucket + obj
// }

func (c *S3) InitUploader() (*s3manager.Uploader, error) {
	return c.initUploader()
}

func (c *S3) initDownloader() (*s3manager.Downloader, error) {
	//get credentials
	creds := credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, c.Token)
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}

	//create session
	cfg := aws.NewConfig().WithRegion(c.Region).WithCredentials(creds).WithEndpoint(c.Endpoint).WithS3ForcePathStyle(true)
	s3Session, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(s3Session)

	return downloader, nil
}

func (c *S3) initUploader() (*s3manager.Uploader, error) {
	//get credentials
	creds := credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, c.Token)
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}

	//create session
	cfg := aws.NewConfig().WithRegion(c.Region).WithCredentials(creds).WithEndpoint(c.Endpoint).WithS3ForcePathStyle(true)
	s3Session, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	//create uploader
	uploader := s3manager.NewUploader(s3Session)

	return uploader, nil
}

func (c *S3) CreateFolder(path string) (*s3.PutObjectOutput, error) {
	svc, err := c.GetSVC()
	if err != nil {
		return nil, err
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(path),
	}

	resp, err := svc.PutObject(params)
	return resp, err
}

func (c *S3) DeleteImage(path string) (*s3.DeleteObjectOutput, error) {
	svc, err := c.GetSVC()
	if err != nil {
		return nil, err
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(path),
	}

	// use the iterator to delete the files.
	resp, err := svc.DeleteObject(input)
	return resp, err
}

func (c *S3) GetFileList(ctx context.Context, path string) ([]string, error) {
	//get credentials
	creds := credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, c.Token)
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}

	//create session
	cfg := aws.NewConfig().WithRegion(c.Region).WithCredentials(creds).WithEndpoint(c.Endpoint).WithS3ForcePathStyle(true)
	s3Session, err := session.NewSession(cfg)
	if err != nil {
		return nil, err
	}

	s3client := s3.New(s3Session)

	s3Keys := make([]string, 0)

	if err := s3client.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(c.Bucket),
		Prefix: aws.String(path), // list files in the directory.
	}, func(o *s3.ListObjectsOutput, b bool) bool { // callback func to enable paging.
		for _, o := range o.Contents {
			s3Keys = append(s3Keys, *o.Key)
		}
		return true
	}); err != nil {
		return nil, err
	}

	return s3Keys, nil
}

func (c *S3) DownloadFile(file *os.File, path string) (int64, error) {
	downloader, err := c.initDownloader()
	if err != nil {
		return 0, err
	}

	byt, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(c.Bucket),
			Key:    aws.String(path),
		})
	if err != nil {
		return 0, err
	}

	return byt, err
}
