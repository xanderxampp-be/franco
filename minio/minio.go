package minio

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	S3Host     string
	S3Username string
	S3Password string
	Region     string
	Secure     bool
}

/*
	 This funcion for create session minio
		Do create a session first when uploading or downloading a file on minio
	    S3Host: Url Minio
	    S3Username: Username Minio
	    S3Password: Password Minio
	    Region: The set region was used on minio, confirm to the administrator to get this, by now we use sa-east-1
	    Secure: Set secure or not, I suggest using default false
*/
func CreateSession(config MinioConfig) (*minio.Client, error) {
	return minio.New(config.S3Host, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3Username, config.S3Password, ""),
		Secure: config.Secure,
		//Region: config.Region,
	})
}

// decoder base64 string to base64 byte
func DecodeBase64(base64String string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64String)
}

// encoder base64 string to base64 byte
func EncodeBase64(base64Byte []byte) string {
	return base64.StdEncoding.EncodeToString(base64Byte)
}

/*
	 This Function for upload a file to minio
	    minioClient: It was session that you already create using method CreateSession
		ctx: context
		bucket: Bucket minio, confirm to the administrator to get this, for example on revamp registration brimo we use bucket brimo-registration
		fileName: File name, don't forget to include the extention, for example: fr01_nasabah.pdf
		contentType: content type based64, so minio can detect the file type, for example: application/pdf, image/png
		decodeBase64: decode base64, you can use DecodeBase64(base64String) for get the decode byte
*/
func PutObject(minioClient *minio.Client, ctx context.Context, bucket, fileName, contentType string, decodeBase64 []byte) (minio.UploadInfo, error) {
	return minioClient.PutObject(ctx, bucket, fileName, bytes.NewReader(decodeBase64), -1,
		minio.PutObjectOptions{
			ContentType: contentType})
}

/*
This Function for downloading a file from minio

	minioClient: It was a session that you already create using the method CreateSession
	ctx: context
	bucket: Bucket minio, confirm to the administrator to get this, for example on revamp registration brimo we use bucket brimo-registration
	path: path inside bucket
	contentType: content type based64, so minio can detect the file type, for example: application/pdf, image/png
	decodeBase64: decode base64, you can use DecodeBase64(base64String) for get the decode byte
*/
func GetObject(minioClient *minio.Client, ctx context.Context, bucket, path string) (*minio.Object, error) {
	return minioClient.GetObject(ctx, bucket, path, minio.GetObjectOptions{})
}

// read object minio from GetObject, will return base64 byte
func ReadObjectMinio(obj *minio.Object) ([]byte, error) {
	return ioutil.ReadAll(obj)
}

/*
This Function for upload a file to minio using a presigned url

	minioClient: it was session that you already create using method CreateSession
	ctx: context
	bucketName: bucket minio, confirm to the administrator to get this, for example on revamp registration brimo we use bucket brimo-registration
	objectName: object or file name, don't forget to include the extention, for example: fr01_nasabah.pdf
	expires: time duration before the url expires, for example 1 * time.Hour
*/
func PresignedPutObject(minioClient *minio.Client, ctx context.Context, bucketName, objectName string, expires time.Duration) (*url.URL, error) {
	return minioClient.PresignedPutObject(ctx, bucketName, objectName, expires)
}

/*
This Function for downloading a file from minio

	minioClient: it was a session that you already create using the method CreateSession
	ctx: context
	bucketName: bucket minio, confirm to the administrator to get this, for example on revamp registration brimo we use bucket brimo-registration
	objectName: object or file name, don't forget to include the extention, for example: fr01_nasabah.pdf
	expires: time duration before the url expires, for example 1 * time.Hour
	reqParams: parameters to add to the url
*/
func PresignedGetObject(minioClient *minio.Client, ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams url.Values) (*url.URL, error) {
	return minioClient.PresignedGetObject(ctx, bucketName, objectName, expires, reqParams)
}
