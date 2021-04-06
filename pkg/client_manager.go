package pkg

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3ListObjectsAPI defines the interface for the ListObjectsV2 function.
// We use this interface to test the function using a mocked service.
type S3ListObjectsAPI interface {
	ListObjectsV2(ctx context.Context,
		params *s3.ListObjectsV2Input,
		optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

// S3PutObjectAPI defines the interface for the PutObject function.
// We use this interface to test the function using a mocked service.
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// GetObjects retrieves the objects in an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a ListObjectsV2Output object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to ListObjectsV2
func GetObjects(c context.Context, api S3ListObjectsAPI, input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error) {
	return api.ListObjectsV2(c, input)
}

// PutFile uploads a file to an Amazon Simple Storage Service (Amazon S3) bucket
// Inputs:
//     c is the context of the method call, which includes the AWS Region
//     api is the interface that defines the method call
//     input defines the input arguments to the service call.
// Output:
//     If success, a PutObjectOutput object containing the result of the service call and nil
//     Otherwise, nil and an error from the call to PutObject
func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

//CreateClient create a s3 client
func CreateClient() *s3.Client {
	mycredentials := credentials.NewStaticCredentialsProvider("key_id", "access_key", "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-1"),
		config.WithCredentialsProvider(mycredentials),
	)
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := s3.NewFromConfig(cfg)

	return client
}

//GetImages get s3 image info
func GetImages(client *s3.Client, bucketName string) {
	bucket := aws.String(bucketName)

	if *bucket == "" {
		fmt.Println("You must supply the name of a bucket (-b BUCKET)")
		return
	}

	input := &s3.ListObjectsV2Input{
		Bucket: bucket,
	}

	resp, err := GetObjects(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error retrieving list of objects:")
		fmt.Println(err)
		return
	}

	fmt.Println("Objects in " + *bucket + ":")

	for _, item := range resp.Contents {
		fmt.Println("Name:          ", *item.Key)
		fmt.Println("Last modified: ", *item.LastModified)
		fmt.Println("Size:          ", item.Size)
		fmt.Println("Storage class: ", item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", *bucket)
	fmt.Println("")
}

//UploadImage to upload one image to S3
func UploadImage(client *s3.Client, bucketName string, directory string, filename string) bool {

	bucket := aws.String(bucketName)
	path := aws.String(directory + filename)
	key := aws.String("album1/" + filename)

	if *bucket == "" || *path == "" {
		fmt.Println("You must supply a bucket name and file path")
		return false
	}

	file, err := os.Open(*path)
	if err != nil {
		fmt.Println("Unable to open file " + *path)
		return false
	}

	defer file.Close()

	//metadata := GetMetadata(file)
	//fileInfo, _ := file.Stat()
	//var size int64 = fileInfo.Size()
	//buffer := make([]byte, size)
	//file.Read(buffer)
	//fileBytes := bytes.NewReader(buffer)
	//fileType := http.DetectContentType(buffer)

	input := &s3.PutObjectInput{
		Bucket: bucket,
		Key:    key,
		Body:   file,
		//ContentLength: size,
		ContentType: aws.String("Image/jpeg"),
		//Metadata:    metadata,
	}
	_, err = PutFile(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
		return false
	}

	fmt.Println(">>>>>>Upload " + filename + " Successfully!!")
	return true
}

//UploadImageList to upload photo list
func UploadImageList(client *s3.Client, bucketName string, directory string, fileList []string) bool {

	for _, filename := range fileList {
		if IsFile(directory + filename) {
			UploadImage(client, bucketName, directory, filename)
			UploadExifTxt(client, bucketName, directory, filename)
		}
	}
	return true
}

//UploadExifTxt upload exif information of photo
func UploadExifTxt(client *s3.Client, bucketName string, directory string, filename string) bool {

	bucket := aws.String(bucketName)
	path := aws.String(directory + filename)

	if *bucket == "" || *path == "" {
		fmt.Println("You must supply a bucket name and file path")
		return false
	}
	file, err := os.Open(*path)
	if err != nil {
		fmt.Println("Unable to open file " + *path)
		return false
	}

	defer file.Close()

	filenameOnly := GetFileNameOnly(filename)
	jsonString := GetMetadataString(file)
	exifkey := aws.String("exif/" + filenameOnly + "_exif.txt")
	stringinput := &s3.PutObjectInput{
		Bucket:      bucket,
		Key:         exifkey,
		Body:        strings.NewReader(jsonString),
		ContentType: aws.String("plain/text"),
	}

	_, err = PutFile(context.TODO(), client, stringinput)
	if err != nil {
		fmt.Println("Got error uploading file description:")
		fmt.Println(err)
		return false
	}
	return true
}
