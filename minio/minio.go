package minio

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"go.uber.org/fx"

	"pu/config"
)

type (
	FileStorage interface {
		Get(name string) (_ []byte, err error)
		Save(ctx context.Context, name string, content []byte) (err error)
		GetURL(ctx context.Context, name string, ttl time.Duration) (_ string, err error)
	}
	mConf struct {
		Endpoint        string `yaml:"endpoint"`
		Bucket          string `yaml:"bucket"`
		Location        string `yaml:"location"`
		AccessKeyID     string `yaml:"accessKeyId"`
		SecretAccessKey string `yaml:"secretAccessKey"`
	}
	S3 struct {
		conf *mConf
		conn *minio.Client
		lock *sync.RWMutex
	}
)

func New(lc fx.Lifecycle, config config.Config) (_ FileStorage, err error) {
	var mCfg *mConf
	if err = config.UnmarshalKey("minio", &mCfg); err != nil {
		return nil, err
	}

	s := &S3{
		conf: mCfg,
		lock: &sync.RWMutex{},
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if s.conn, err = minio.New(mCfg.Endpoint, &minio.Options{
				Creds: credentials.NewStaticV4(
					mCfg.AccessKeyID,
					mCfg.SecretAccessKey,
					"",
				),
				Region: mCfg.Location,
			}); err != nil {
				return err
			}

			var exists bool
			if exists, err = s.conn.BucketExists(context.Background(), mCfg.Bucket); err != nil {
				return err

			}

			if !exists {
				if err = s.conn.MakeBucket(context.Background(), mCfg.Bucket, minio.MakeBucketOptions{
					Region: mCfg.Location,
				}); err != nil {
					return err
				}
			}
			return nil
		},
	})

	return s, nil
}

func (s *S3) s3() *minio.Client {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.conn
}

// Get return []byte
// example github: https://github.com/minio/minio-go/blob/master/examples/s3/getobject.go
func (s *S3) Get(name string) (_ []byte, err error) {
	var reader *minio.Object
	if reader, err = s.s3().GetObject(
		context.Background(),
		s.conf.Bucket,
		name,
		minio.GetObjectOptions{},
	); err != nil {
		return nil, errors.Wrap(err, "get object error")
	}
	defer reader.Close()

	var stat minio.ObjectInfo
	if stat, err = reader.Stat(); err != nil {
		return nil, errors.Wrap(err, "reader error")
	}

	var buf = &bytes.Buffer{}
	if _, err = io.CopyN(buf, reader, stat.Size); err != nil {
		return nil, errors.Wrap(err, "copy error")
	}
	return buf.Bytes(), nil
}

func (s *S3) Save(ctx context.Context, name string, content []byte) (err error) {
	buf := bytes.NewBuffer(content)

	_, err = s.s3().PutObject(
		ctx,
		s.conf.Bucket,
		name,
		buf,
		int64(buf.Len()),
		minio.PutObjectOptions{
			ContentType: http.DetectContentType(buf.Bytes()),
		},
	)

	return err
}

func (s *S3) GetURL(ctx context.Context, name string, ttl time.Duration) (_ string, err error) {
	var reqParams = make(url.Values)
	reqParams.Add("Content-Disposition", `attachment; filename="`+name+`"`)

	var tempURL *url.URL
	if tempURL, err = s.s3().PresignedGetObject(
		ctx,
		s.conf.Bucket,
		name,
		ttl,
		reqParams,
	); err != nil {
		return "", err
	}

	return tempURL.String(), nil
}
