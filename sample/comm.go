package sample

import (
	"fmt"
	"github.com/wilac-pv/ksyun-ks3-go-sdk/ks3"
	"os"
	"strings"
	"time"
)

var (
	pastDate   = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	futureDate = time.Date(2049, time.January, 10, 23, 0, 0, 0, time.UTC)
)

// HandleError is the error handling method in the sample code
func HandleError(err error) {
	fmt.Println("occurred error:", err)
	os.Exit(-1)
}

// GetTestBucket creates the test bucket
func GetTestBucket(bucketName string) (*ks3.Bucket, error) {
	// New client
	client, err := ks3.New(endpoint, accessID, accessKey)
	if err != nil {
		return nil, err
	}

	// Create bucket
	err = client.CreateBucket(bucketName)
	if err != nil {
		return nil, err
	}

	// Get bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// DeleteTestBucketAndLiveChannel 删除sample的channelname和bucket，该函数为了简化sample，让sample代码更明了
func DeleteTestBucketAndLiveChannel(bucketName string) error {
	// New Client
	client, err := ks3.New(endpoint, accessID, accessKey)
	if err != nil {
		return err
	}

	// Get Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	marker := ""
	for {
		result, err := bucket.ListLiveChannel(ks3.Marker(marker))
		if err != nil {
			HandleError(err)
		}

		for _, channel := range result.LiveChannel {
			err := bucket.DeleteLiveChannel(channel.Name)
			if err != nil {
				HandleError(err)
			}
		}

		if result.IsTruncated {
			marker = result.NextMarker
		} else {
			break
		}
	}

	// Delete Bucket
	err = client.DeleteBucket(bucketName)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTestBucketAndObject deletes the test bucket and its objects
func DeleteTestBucketAndObject(bucketName string) error {
	// New client
	client, err := ks3.New(endpoint, accessID, accessKey)
	if err != nil {
		return err
	}

	// Get bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// Delete part
	keyMarker := ks3.KeyMarker("")
	uploadIDMarker := ks3.UploadIDMarker("")
	for {
		lmur, err := bucket.ListMultipartUploads(keyMarker, uploadIDMarker)
		if err != nil {
			return err
		}
		for _, upload := range lmur.Uploads {
			var imur = ks3.InitiateMultipartUploadResult{Bucket: bucket.BucketName,
				Key: upload.Key, UploadID: upload.UploadID}
			err = bucket.AbortMultipartUpload(imur)
			if err != nil {
				return err
			}
		}
		keyMarker = ks3.KeyMarker(lmur.NextKeyMarker)
		uploadIDMarker = ks3.UploadIDMarker(lmur.NextUploadIDMarker)
		if !lmur.IsTruncated {
			break
		}
	}

	// Delete objects
	marker := ks3.Marker("")
	for {
		lor, err := bucket.ListObjects(marker)
		if err != nil {
			return err
		}
		for _, object := range lor.Objects {
			err = bucket.DeleteObject(object.Key)
			if err != nil {
				return err
			}
		}
		marker = ks3.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}

	// Delete bucket
	err = client.DeleteBucket(bucketName)
	if err != nil {
		return err
	}

	return nil
}

// Object defines pair of key and value
type Object struct {
	Key   string
	Value string
}

// CreateObjects creates some objects
func CreateObjects(bucket *ks3.Bucket, objects []Object) error {
	for _, object := range objects {
		err := bucket.PutObject(object.Key, strings.NewReader(object.Value))
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteObjects deletes some objects.
func DeleteObjects(bucket *ks3.Bucket, objects []Object) error {
	for _, object := range objects {
		err := bucket.DeleteObject(object.Key)
		if err != nil {
			return err
		}
	}
	return nil
}
