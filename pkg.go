// Package aws_session provides a few small functions for creating an AWS SDK Session object. I find these hard to
// remember if it's been a while, so I created this small package.
package awssession

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

// Create makes an attempt to use the profile and/or the environment variables. If AWS_SDK_LOAD_CONFIG it will use the
// profile. If all of AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION then these will be used. If neither of
// those conditions is met, then AWS_SDK_LOAD_CONFIG will be set to true and the profile will be attempted.
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

	if sess, err := CreateUsingProfile(); err != nil {
		return sess, nil
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

// CreateUsingProfile sets AWS_SDK_LOAD_CONFIG to true and tries to use the profile.
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

// CreateUsingEvironment uses the AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY to create the Session.
func CreateUsingEnvironment(region string) (*session.Session, error) {
	conf := aws.Config{Region: aws.String(region)}
	return session.NewSession(&conf)
}

// CreateUsingStrings creates a Session with the strings passed in.
func CreateUsingStrings(awsAccessKey, awsSecretKey, awsRegion string) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
	})
}
