package minio

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/xanderxampp-be/franco/contextwrap"
	"github.com/xanderxampp-be/franco/trace"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.elastic.co/apm"
)

// MiniOop wrap minioClient native to be customed with brimo's approach
type MinioOop struct {
	minioClient *minio.Client
	endpoint    string
}

// NewClient create / instantiate minio client
func NewClient(id string, secret string, token string, isSecure bool, endpointMinio string) (*MinioOop, error) {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(id, secret, token),
		Secure: isSecure,
	}
	minioClient, err := minio.New(endpointMinio, opts)
	if err != nil {
		return nil, err
	}

	minclientBrimo := &MinioOop{
		minioClient: minioClient,
		endpoint:    endpointMinio,
	}

	return minclientBrimo, nil
}

// PutObject do the upload process of an object to minio bucket by byte basis, return information of file if the upload process success
func (m *MinioOop) PutObject(ctx context.Context, bucketName string, objectName string, objectBase64 []byte, objectSize int64, contentType string) (minio.UploadInfo, error) {
	apmSpan, _ := apm.StartSpan(ctx, "PutObject", "Minio")
	defer apmSpan.End()

	uploadInfo, err := m.minioClient.PutObject(ctx, bucketName, objectName, bytes.NewReader(objectBase64), objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println("put object minio error : ", err)
	}
	return uploadInfo, err
}

// GetObject do the download process of an object from minio bucket by byte basis, return minioObject if the process sucess
func (m *MinioOop) GetObject(ctx context.Context, bucketName string, objectName string) (*minio.Object, error) {
	apmSpan, _ := apm.StartSpan(ctx, "GetObject", "Minio")
	defer apmSpan.End()

	minioObj, err := m.minioClient.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("get object minio error : ", err)
	}
	return minioObj, err
}

// FGetObject do the download process of an object from minio bucket by file basis, return error if any occured
func (m *MinioOop) FGetObject(ctx context.Context, bucketName, objectName, filepath string) (context.Context, error) {
	start := time.Now()
	apmSpan, _ := apm.StartSpan(ctx, "FGetObject", "Minio")
	defer apmSpan.End()

	trOri := contextwrap.GetTraceFromContext(ctx)
	err := m.minioClient.FGetObject(ctx, bucketName, objectName, filepath, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("get object minio error : ", err)
	}

	tr := &trace.TraceMinio{
		Host:       m.endpoint,
		ObjectName: objectName,
		BucketName: bucketName,
		Elapsed:    time.Since(start).String(),
	}

	trProcessed := append(trOri, tr)
	ctx = contextwrap.SetTraceFromContext(ctx, trProcessed)

	return ctx, err
}

func (m *MinioOop) FPutObject(ctx context.Context, bucketName, objectName, filepath string, contentType string) (context.Context, minio.UploadInfo, error) {
	start := time.Now()
	apmSpan, _ := apm.StartSpan(ctx, "FPutObject", "Minio")
	defer apmSpan.End()

	trOri := contextwrap.GetTraceFromContext(ctx)
	uploadInfo, err := m.minioClient.FPutObject(ctx, bucketName, objectName, filepath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println("put object minio error : ", err)
	}

	tr := &trace.TraceMinio{
		Host:       m.endpoint,
		ObjectName: objectName,
		BucketName: bucketName,
		Elapsed:    time.Since(start).String(),
	}

	trProcessed := append(trOri, tr)
	ctx = contextwrap.SetTraceFromContext(ctx, trProcessed)

	return ctx, uploadInfo, err
}

// StatObject return minio info related to the existence of the file in bucket minio
// if file not exist then the statObject return error
func (m *MinioOop) StatObject(ctx context.Context, bucketName string, objectName string) (minio.ObjectInfo, error) {
	apmSpan, _ := apm.StartSpan(ctx, "StatObject", "Minio")
	defer apmSpan.End()

	minioObj, err := m.minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("get object minio error :", err)
	}
	return minioObj, err
}

// BucketExists return bool, true if bucket in minio with correspondent bucketName param, exist
func (m *MinioOop) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	apmSpan, _ := apm.StartSpan(ctx, "BucketExists", "Minio")
	defer apmSpan.End()

	isExist, err := m.minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		fmt.Println("check existence of bucket get error : ", err)
	}

	return isExist, err
}

// ListObjects return <-chan minio.ObjectInfo, a list contains object info of an object in the folder
func (m *MinioOop) ListObjects(ctx context.Context, bucketName, folderName string) (context.Context, <-chan minio.ObjectInfo) {
	apmSpan, _ := apm.StartSpan(ctx, "ListObjects", "Minio")
	defer apmSpan.End()

	listObjects := m.minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{
		Prefix:    folderName,
		Recursive: true,
	})

	if len(listObjects) == 0 {
		fmt.Println(fmt.Printf("no object exist inside bucket %v folder %v", bucketName, folderName))
	}

	return ctx, listObjects
}

// RemoveObject return ctx, err, used for remove object
func (m *MinioOop) RemoveObject(ctx context.Context, bucketName string, object minio.ObjectInfo) (context.Context, error) {
	apmSpan, _ := apm.StartSpan(ctx, "RemoveObject", "Minio")
	defer apmSpan.End()

	err := m.minioClient.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{})
	if err != nil {
		fmt.Println(fmt.Printf("failed remove object %v, cause : %v", object.Key, err.Error()))
	}

	return ctx, err
}

// CopyObject return ctx, error. Used for copy object from a bucket to another bucket
func (m *MinioOop) CopyObject(ctx context.Context, objectName, destination, source string) (context.Context, minio.UploadInfo, error) {
	apmSpan, _ := apm.StartSpan(ctx, "CopyObject", "Minio")
	defer apmSpan.End()

	info, err := m.minioClient.CopyObject(ctx,
		minio.CopyDestOptions{
			Bucket: destination,
			Object: objectName,
		}, minio.CopySrcOptions{
			Bucket: source,
			Object: objectName,
		})
	if err != nil {
		fmt.Println(fmt.Printf("failed copy object %v, cause : %v", objectName, err.Error()))
	}

	return ctx, info, err
}

// RemoveObject return ctx, err, used for remove object
func (m *MinioOop) RemoveObjectWithBypassGovernance(ctx context.Context, bucketName string, object minio.ObjectInfo) (context.Context, error) {
	apmSpan, _ := apm.StartSpan(ctx, "RemoveObjectWithBypassGovernance", "Minio")
	defer apmSpan.End()

	err := m.minioClient.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{
		GovernanceBypass: true,
	})

	if err != nil {
		fmt.Println(fmt.Printf("failed remove object %v, cause : %v", object.Key, err.Error()))
	}

	return ctx, err
}

// PresignedPutObject return a presigned url to do the upload process of an object to minio, return error if any occured
func (m *MinioOop) PresignedPutObject(ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error) {
	apmSpan, _ := apm.StartSpan(ctx, "PresignedPutObject", "Minio")
	defer apmSpan.End()

	uploadUrl, err := m.minioClient.PresignedPutObject(ctx, bucketName, objectName, expires)
	if err != nil {
		fmt.Println("generate presigned url put object minio error : ", err)
	}
	return uploadUrl, err
}

// PresignedGetObject return a presigned url to do the download process of an object from minio, return error if any occured
func (m *MinioOop) PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	apmSpan, _ := apm.StartSpan(ctx, "PresignedGetObject", "Minio")
	defer apmSpan.End()

	downloadUrl, err := m.minioClient.PresignedGetObject(ctx, bucketName, objectName, expires, reqParams)
	if err != nil {
		fmt.Println("generate presigned url get object minio error : ", err)
	}
	return downloadUrl, err
}

// ForceRemoveObject remove all version of the object, return ctx, err, used for remove object
func (m *MinioOop) ForceRemoveObject(ctx context.Context, bucketName string, object minio.ObjectInfo) (context.Context, error) {
	apmSpan, _ := apm.StartSpan(ctx, "RemoveObject", "Minio")
	defer apmSpan.End()

	err := m.minioClient.RemoveObject(ctx, bucketName, object.Key, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		fmt.Println(fmt.Printf("failed remove object %v, cause : %v", object.Key, err.Error()))
	}

	return ctx, err
}
