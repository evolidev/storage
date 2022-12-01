package storage

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/evolidev/blitza/disk"
	"github.com/evolidev/blitza/fs"
	"io"
	"testing"
)

func TestS3Config(t *testing.T) {
	t.Run("simple credentials should work", func(t *testing.T) {
		config := disk.S3Config{}
		config.Key = "test"
		config.Secret = "test"
		config.EndpointResolver = s3.EndpointResolverFromURL("https://test.com")

		storage := disk.NewS3(config)

		options := storage.Options()

		cred, _ := options.Credentials.Retrieve(context.TODO())

		if cred.AccessKeyID != "test" {
			t.Errorf("Wrong access key")
		}

		if cred.SecretAccessKey != "test" {
			t.Errorf("Wrong access key")
		}
	})

	t.Run("simple credentials objects should work", func(t *testing.T) {
		tmpCred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("key", "secret", ""))
		config := disk.S3Config{}
		config.Credentials = tmpCred
		config.Endpoint = "test"

		storage := disk.NewS3(config)

		options := storage.Options()

		cred, _ := options.Credentials.Retrieve(context.TODO())

		if cred.AccessKeyID != "key" {
			t.Errorf("Wrong access key")
		}

		if cred.SecretAccessKey != "secret" {
			t.Errorf("Wrong secret key")
		}
	})

	t.Run("callback options should be invoked", func(t *testing.T) {
		config := disk.S3Config{}
		config.OptionsFuncs = make([]func(options *s3.Options), 0)
		tmp := false
		config.OptionsFuncs = append(config.OptionsFuncs, func(options *s3.Options) {
			options.Region = "test"
			tmp = true
		})

		disk.NewS3(config)

		if !tmp {
			t.Errorf("wrong region")
		}
	})
}

func TestAttributes(t *testing.T) {
	t.Run("Attributes should fallback to Head if object attributes fails", func(t *testing.T) {
		config := disk.S3Config{}
		config.Client = &headFallback{MemoryClient: disk.NewMemoryClient()}
		c := disk.NewS3(config)
		c.Put("test.txt", []byte("test"), fs.PUBLIC)

		a := c.Attributes("test.txt")

		if a.Size != int64(len([]byte("test"))) {
			t.Errorf("Size %d does not match expected %d", a.Size, int64(len([]byte("test"))))
		}
	})
}

func TestGetShouldPassReadError(t *testing.T) {
	config := disk.S3Config{}
	config.Client = &getFail{MemoryClient: disk.NewMemoryClient()}
	c := disk.NewS3(config)

	_, err := c.Get("test")
	if err == nil {
		t.Errorf("error expected")
	}
}

func TestDeleteShouldPassError(t *testing.T) {
	t.Run("Fail for files", func(t *testing.T) {
		config := disk.S3Config{}
		config.Client = &deleteFail{MemoryClient: disk.NewMemoryClient(), size: 0}
		c := disk.NewS3(config)

		err := c.DeleteDirectory("test")
		if err == nil {
			t.Errorf("error expected")
		}
	})

	t.Run("Fail for directories", func(t *testing.T) {
		config := disk.S3Config{}
		config.Client = &deleteFail{MemoryClient: disk.NewMemoryClient(), size: 1}
		c := disk.NewS3(config)

		err := c.DeleteDirectory("test")
		if err == nil {
			t.Errorf("error expected")
		}
	})
}

func TestListShouldPassError(t *testing.T) {
	t.Parallel()
	config := disk.S3Config{}
	config.Client = &listFail{MemoryClient: disk.NewMemoryClient()}
	c := disk.NewS3(config)

	t.Run("Files should return error", func(t *testing.T) {
		err := c.Files("test")
		if err == nil {
			t.Errorf("error expected")
		}
	})

	t.Run("Directories should return error", func(t *testing.T) {
		err := c.Directories("test")
		if err == nil {
			t.Errorf("error expected")
		}
	})

}

type headFallback struct {
	*disk.MemoryClient
}

func (h headFallback) GetObjectAttributes(ctx context.Context, params *s3.GetObjectAttributesInput, optFns ...func(*s3.Options)) (*s3.GetObjectAttributesOutput, error) {
	return nil, errors.New("fail")
}

type getFail struct {
	*disk.MemoryClient
}

func (g *getFail) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	rc := io.NopCloser(failReader{})
	o := &s3.GetObjectOutput{Body: rc}

	return o, nil
}

type failReader struct {
}

func (f failReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("fail")
}

type deleteFail struct {
	size int64
	*disk.MemoryClient
}

func (d *deleteFail) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return nil, errors.New("fail")
}

func (d *deleteFail) ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	r := make([]types.Object, 0)
	t := "test"

	r = append(r, types.Object{Key: &t, Size: d.size})

	o := &s3.ListObjectsOutput{}
	o.Contents = r

	return o, nil
}

type listFail struct {
	*disk.MemoryClient
}

func (l *listFail) ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	return nil, errors.New("fail")
}
