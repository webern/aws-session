package aws_session

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	"os"
)

const (
	constKey    = "AWS_ACCESS_KEY_ID"
	constSecret = "AWS_SECRET_ACCESS_KEY"
	constRegion = "AWS_REGION"
	constConfig = "AWS_SDK_LOAD_CONFIG"
)

func Create() (*session.Session, error) {
	if val, ok := os.LookupEnv(constConfig); ok {
		if val != "false" {
			return CreateUsingProfile()
		}
	}
	_, keyOK := os.LookupEnv(constKey)
	_, secretOK := os.LookupEnv(constSecret)
	region, regionOK := os.LookupEnv(constRegion)
	if keyOK && secretOK && regionOK {
		return CreateUsingEnvironment(region)
	}

	return &session.Session{}, errors.New(
		fmt.Sprintf("Either %s must be set to true, or all of %s, %s, and %s must be set",
			constConfig,
			constKey,
			constSecret,
			constRegion,
		),
	)
}

func CreateUsingProfile() (*session.Session, error) {
	err := os.Setenv(constConfig, "true")
	if err != nil {
		return &session.Session{}, errors.Wrap(
			err,
			"Error while setting the environment variable AWS_SDK_LOAD_CONFIG",
		)
	}
	return session.NewSession(&aws.Config{})
}

func CreateUsingEnvironment(region string) (*session.Session, error) {
	conf := aws.Config{Region: aws.String(region)}
	return session.NewSession(&conf)
}

func CreateUsingStrings(awsAccessKey, awsSecretKey, awsRegion string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
}
