package disk

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
)

type LocalConfig struct {
	PermModeFilePublic       os.FileMode
	PermModeFilePrivate      os.FileMode
	PermModeDirectoryPublic  os.FileMode
	PermModeDirectoryPrivate os.FileMode
	Prefix                   string
}

type S3Config struct {
	Client           Client
	OptionsFuncs     []func(*s3.Options)
	Credentials      aws.CredentialsProvider
	Key              string
	Secret           string
	Bucket           string
	Endpoint         string
	EndpointResolver s3.EndpointResolver
	Prefix           string
}

type MemoryConfig struct {
	Prefix string
}
