/*
Copyright 2017 The GoStor Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package storage

import (
	"errors"
	"time"

	"github.com/gostor/gofs/pkg/api"
	minio "github.com/minio/minio-go"
)

type MinioStorage struct {
	Endpoint string
}

type MinioBucket struct {
	client *minio.Client
	Name   string
}

// ObjectInfo - represents object metadata.
type ObjectInfo struct {
	// Name of the bucket.
	Bucket string

	// Name of the object.
	Name string

	// Date and time when the object was last modified.
	ModTime time.Time

	// Total object size.
	Size int64

	// IsDir indicates if the object is prefix.
	IsDir bool

	// Hex encoded unique entity tag of the object.
	ETag string

	// A standard MIME type describing the format of the object.
	ContentType string

	// Specifies what content encodings have been applied to the object and thus
	// what decoding mechanisms must be applied to obtain the object referenced
	// by the Content-Type header field.
	ContentEncoding string
}

func NewMinioStorage(ep string) *MinioStorage {
	return &MinioStorage{
		Endpoint: ep,
	}
}

func (ms *MinioStorage) Ping() error {
	return nil
}

func (ms *MinioStorage) Bucket(name string, cfg *api.Config) (Bucket, error) {
	minioClient, err := minio.New(ms.Endpoint, cfg.AccessKey, cfg.SecretKey, true)
	if err != nil {
		return nil, err
	}

	return &MinioBucket{
		client: minioClient,
		Name:   name,
	}, nil
}

func (mb *MinioBucket) Auth() error {
	return nil
}

func (mb *MinioBucket) Create() error {
	return nil
}

func (mb *MinioBucket) Delete() error {
	return nil
}

func (mb *MinioBucket) Get() error {
	return nil
}

func (mb *MinioBucket) Object(name string) Object {
	return &ObjectInfo{
		Bucket: mb.Name,
		Name:   name,
	}
}

func (mo *ObjectInfo) Stat() error {
	return ErrNoSuchObject
}

func (mo *ObjectInfo) Delete() error {
	return ErrNoSuchObject
}

// ErrNoSuchObject - returned when object is not found.
var ErrNoSuchObject = errors.New("No such object")

// ErrNoSuchBucket - returned when bucket is not found.
var ErrNoSuchBucket = errors.New("No such bucket")

// IsNoSuchObject - is err ErrNoSuchObject ?
func IsNoSuchObject(err error) bool {
	if err == nil {
		return false
	}
	// Validate if the type is same as well.
	if err == ErrNoSuchObject {
		return true
	} else if err.Error() == ErrNoSuchObject.Error() {
		// Reaches here when type did not match but err string matches.
		// Someone wrapped this error? - still return true since
		// they are the same.
		return true
	}
	errorResponse := minio.ToErrorResponse(err)
	return errorResponse.Code == "NoSuchKey"
}
