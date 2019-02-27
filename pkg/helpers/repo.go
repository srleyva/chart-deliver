package helpers

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Repo is the  interface used to talk to repo
type Repo interface {
	GetFiles() (map[string]string, error)
}

// NewRepo returns new repo based on type
func NewRepo(provider, bucket, path string) (Repo, error) {
	switch provider {
	case "GCS":
		return newGCS(bucket, path)
	case "S3":
		return nil, fmt.Errorf("Not implemented %s", provider)
	}

	return nil, fmt.Errorf("Provider type %s not known", provider)
}

// GCS talks to GCS to get template
type GCS struct {
	client *storage.Client
	bucket string
	path   string
}

func newGCS(bucket, path string) (*GCS, error) {
	client, err := storage.NewClient(context.TODO())
	if err != nil {
		return nil, err
	}
	return &GCS{
		client: client,
		bucket: bucket,
		path:   path,
	}, nil
}

// GetFiles grabs templates from bucket
func (g *GCS) GetFiles() (map[string]string, error) {
	bh := g.client.Bucket(g.bucket)
	// Next check if the bucket exists
	if _, err := bh.Attrs(context.TODO()); err != nil {
		return nil, err
	}

	files := make(map[string]string)

	// Get items in dir
	items := bh.Objects(context.TODO(), &storage.Query{
		Prefix:    g.path,
		Delimiter: "/",
	})

	for {
		attrs, err := items.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		reader, err := bh.Object(attrs.Name).NewReader(context.TODO())
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		files[attrs.Name] = string(data)
	}

	return files, nil
}

// S3 talks to s3 to pull templates
type S3 struct{}
