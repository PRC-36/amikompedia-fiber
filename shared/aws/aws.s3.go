package aws

import (
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsSess "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewSessionAWSS3(c util.Config) (*s3.S3, error) {

	sess, err := awsSess.NewSession(&aws.Config{
		Region:      aws.String(c.AWSRegion),
		Credentials: credentials.NewStaticCredentials(c.AWSAccessKey, c.AWSSecretKey, ""),
	})

	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess)

	return s3Client, nil
}

//func NewS3Aws(session *awsSess.Session, bucket string) *s3.S3 {
//	return s3.New(session)
//}
