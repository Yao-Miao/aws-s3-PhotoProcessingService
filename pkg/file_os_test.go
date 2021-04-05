package pkg

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestIsDir(t *testing.T) {
	var dirctory string = "/Users/miaoyao/Desktop/NEU/CSYE 6225 Cloud and Distributed Systems/Assignments/Assignment2/CSYE 6225 Assignment 2 - Photo Processing Service - Yao Miao/photo"
	actual := IsDir(dirctory)
	if !actual {
		t.Errorf("TestIsDir() failed dirctory=%s", dirctory)
	}
}
func TestReadDir(t *testing.T) {
	fileList := ReadDir("./asset/photo")
	fmt.Println(fileList)
}

func TestGetMetadata(t *testing.T) {
	path := aws.String("/Users/miaoyao/GoProject/photo_processing_and_storage/asset/photo/Mytest1.png")
	file, err := os.Open(*path)
	if err != nil {
		fmt.Println("Unable to open file " + *path)
	}
	metadata := GetMetadata(file)
	fmt.Println(metadata)
}
