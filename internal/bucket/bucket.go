package bucket

import "io"

type BucketUser interface {
	PushFileToBucket(bucket string, prefix string, file io.Reader) error
}
