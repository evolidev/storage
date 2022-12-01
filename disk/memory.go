package disk

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/evolidev/blitza/fs"
	"io"
	"strings"
	"time"
)

type Memory struct {
	*S3
	config MemoryConfig
}

func NewMemory(config MemoryConfig) *Memory {
	return &Memory{S3: NewS3(S3Config{Client: NewMemoryClient(), Prefix: config.Prefix})}
}

func (m *Memory) Prefix(prefix string) fs.Disk {
	c := m.config
	c.Prefix = prefix

	return &Memory{config: c, S3: m.S3.Prefix(prefix).(*S3)}
}

type MemoryClient struct {
	data map[string]file
}

func NewMemoryClient() *MemoryClient {
	return &MemoryClient{
		data: make(map[string]file),
	}
}

func (m *MemoryClient) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	buf := new(bytes.Buffer)

	if params.Body != nil {
		buf.ReadFrom(params.Body)
	}

	t := time.Now()

	o := types.Object{
		Key:          params.Key,
		LastModified: &t,
		Size:         int64(len(buf.Bytes())),
	}

	f := file{
		object:  o,
		content: buf.Bytes(),
	}

	m.data[*params.Key] = f

	return &s3.PutObjectOutput{}, nil
}

func (m *MemoryClient) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if err := m.check(*params.Key); err != nil {
		return nil, err
	}

	o := s3.GetObjectOutput{}
	o.Body = io.NopCloser(bytes.NewReader(m.data[*params.Key].content))

	return &o, nil
}

func (m *MemoryClient) GetObjectAttributes(ctx context.Context, params *s3.GetObjectAttributesInput, optFns ...func(*s3.Options)) (*s3.GetObjectAttributesOutput, error) {
	if err := m.check(*params.Key); err != nil {
		return nil, err
	}

	o := &s3.GetObjectAttributesOutput{}

	o.ObjectSize = m.data[*params.Key].object.Size
	o.LastModified = m.data[*params.Key].object.LastModified

	return o, nil
}

func (m *MemoryClient) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if err := m.check(*params.Key); err != nil {
		return nil, err
	}

	return &s3.HeadObjectOutput{
		ContentLength: m.data[*params.Key].object.Size,
		LastModified:  m.data[*params.Key].object.LastModified,
	}, nil
}

func (m *MemoryClient) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	ok := false

	for k, _ := range m.data {
		if strings.HasPrefix(k, *params.Key) {
			ok = true
			break
		}
	}

	delete(m.data, *params.Key)

	return &s3.DeleteObjectOutput{DeleteMarker: ok}, nil
}

func (m *MemoryClient) ListObjects(ctx context.Context, params *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	r := make([]types.Object, 0)

	for k, v := range m.data {
		if strings.HasPrefix(k, *params.Prefix) {
			r = append(r, v.object)
		}
	}

	o := &s3.ListObjectsOutput{}
	o.Contents = r

	return o, nil
}

func (m *MemoryClient) check(key string) error {
	if _, ok := m.data[key]; !ok {
		return errors.New("file does not exists")
	}

	return nil
}

type file struct {
	object  types.Object
	content []byte
}
