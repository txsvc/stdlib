package storage

import (
	"context"
	"io"

	"github.com/txsvc/stdlib/pkg/provider"
)

const (
	TypeStorage provider.ProviderType = 20
)

type (
	StorageProvider interface {
		Bucket(string) BucketHandle
	}

	BucketHandle interface {
		Object(string) ObjectHandle
	}

	ObjectHandle interface {
		Close() error
		NewReader(context.Context) (io.Reader, error)
		NewWriter(context.Context) (io.Writer, error)
	}
)

var (
	p *provider.Provider
)

func Bucket(name string) BucketHandle {
	imp, found := p.Find(TypeStorage)
	if !found {
		return nil
	}
	return imp.(StorageProvider).Bucket(name)
}
