package disk

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/evolidev/blitza/fs"
	"strings"
)

type S3 struct {
	client    Client
	bucket    string
	delimiter string
	*Common
	options *s3.Options
	config  S3Config
}

type Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	GetObjectAttributes(ctx context.Context, params *s3.GetObjectAttributesInput, optFns ...func(*s3.Options)) (*s3.GetObjectAttributesOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
}

func NewS3(config S3Config) *S3 {
	s := &S3{
		bucket:    config.Bucket,
		delimiter: "/",
		Common:    &Common{},
		config:    config,
	}

	s.client = s.getClient(config)

	s.disk = s

	return s
}

func (s *S3) Attributes(file string) fs.Attributes {
	var size, lastModified int64
	size = 0
	lastModified = 0

	var attributes []types.ObjectAttributes
	attributes = append(attributes, types.ObjectAttributesObjectSize)

	key := aws.String(s.getPath(file))

	result, err := s.client.GetObjectAttributes(context.TODO(), &s3.GetObjectAttributesInput{
		Bucket:           aws.String(s.bucket),
		Key:              key,
		ObjectAttributes: attributes,
	})

	if err == nil {
		size = result.ObjectSize
		lastModified = result.LastModified.Unix()
	}

	if size == 0 {
		head, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    key,
		})

		if err == nil {
			size = head.ContentLength
		}
	}

	return fs.Attributes{
		Size:         size,
		LastModified: lastModified,
	}
}

func (s *S3) Put(file string, content []byte, visibility fs.Visibility) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.getPath(file)),
		Body:   bytes.NewReader(content),
	})

	return err
}

func (s *S3) Get(file string) ([]byte, error) {
	content, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.getPath(file)),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(content.Body)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *S3) Exists(file string) bool {
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.getPath(file)),
	})

	return err == nil
}

func (s *S3) Path(file string) string {
	return s.getPath(file)
}

func (s *S3) Delete(files ...string) error {
	for _, file := range files {
		tmp, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(s.getPath(file)),
		})

		if err == nil && !tmp.DeleteMarker {
			err = errors.New("could not delete file")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *S3) MakeDirectory(dir string, visibility fs.Visibility) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(s.getPath(dir)),
	})

	return err
}

func (s *S3) DeleteDirectory(dir string) error {
	for _, file := range s.AllFiles(s.getPath(dir)) {
		err := s.Delete(file.Path())

		if err != nil {
			return err
		}
	}

	for _, d := range s.AllDirectories(s.getPath(dir)) {
		err := s.Delete(d.Cwd())

		if err != nil {
			return err
		}
	}

	s.Delete(s.getPath(dir))

	return nil
}

func (s *S3) Files(dir string) []*fs.File {
	r := make([]*fs.File, 0)

	objects, err := s.client.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(s.getPath(dir) + s.delimiter),
	})

	if err != nil {
		return r
	}

	for _, object := range objects.Contents {
		tmp := strings.TrimPrefix(*object.Key, s.getPath(dir)+s.delimiter)
		if strings.Count(tmp, s.delimiter) > 0 || object.Size == 0 {
			continue
		}

		f := fs.NewFile(s, s.getPath(dir), tmp)
		r = append(r, f)
	}

	return r
}

func (s *S3) Directories(dir string) []fs.Disk {
	r := make([]fs.Disk, 0)
	files := make(map[string]fs.Disk, 0)

	objects, err := s.client.ListObjects(context.TODO(), &s3.ListObjectsInput{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String(s.getPath(dir) + s.delimiter),
	})

	if err != nil {
		return r
	}

	for _, object := range objects.Contents {
		tmp := strings.TrimPrefix(*object.Key, s.getPath(dir)+s.delimiter)
		if strings.Count(tmp, s.delimiter) == 0 && object.Size > 0 {
			continue
		}

		sp := strings.Split(tmp, s.delimiter)

		dirName := s.getPath(dir) + s.delimiter + sp[0]

		//f := fs.NewDirectory(s, s.getPath(dir), sp[0])
		f := s.Prefix(s.getPath(dir) + s.delimiter + sp[0])

		files[dirName] = f
	}

	for _, t := range files {
		r = append(r, t)
	}

	return r
}

func (s *S3) Prefix(prefix string) fs.Disk {
	c := s.config
	c.Prefix = prefix

	return NewS3(c)
}

func (s *S3) Cwd() string {
	return s.config.Prefix
}

func (s *S3) Options() *s3.Options {
	return s.options
}

func (s *S3) getClient(config S3Config) Client {
	if config.Client != nil {
		return config.Client
	}

	options := s3.Options{}
	buildCredentials(&config, &options)
	buildEndpoint(&config, &options)

	fns := make([]func(*s3.Options), 0)

	if config.OptionsFuncs != nil {
		fns = config.OptionsFuncs
	}

	s.options = &options

	return s3.New(options, fns...)
}

func buildCredentials(config *S3Config, options *s3.Options) {
	var c aws.CredentialsProvider
	if config.Credentials != nil {
		c = config.Credentials
	} else if config.Key != "" && config.Secret != "" {
		c = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(config.Key, config.Secret, ""))
	}

	if c != nil {
		options.Credentials = c
	}
}

func buildEndpoint(config *S3Config, options *s3.Options) {
	if config.EndpointResolver != nil {
		options.EndpointResolver = config.EndpointResolver
	} else if config.Endpoint != "" {
		options.EndpointResolver = s3.EndpointResolverFromURL(config.Endpoint)
	}
}

func (s *S3) getPath(path string) string {
	if s.config.Prefix == "" {
		return path
	}

	return s.config.Prefix + "/" + path
}
