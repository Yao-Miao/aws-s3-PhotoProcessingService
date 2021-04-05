package pkg

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/tidwall/gjson"
)

//Exists return bool: check if the directory or file exists
func Exists(filepath string) bool {
	_, err := os.Stat(filepath) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//IsDir return bool: checke if it's a directory
func IsDir(filepath string) bool {
	s, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//IsFile return bool: checke if it's a file
func IsFile(filepath string) bool {
	return !IsDir(filepath)
}

//ReadDir return a list of file name
func ReadDir(filepath string) []string {
	var fileList []string = make([]string, 0)
	files, _ := ioutil.ReadDir(filepath)
	for _, f := range files {
		var temp string = f.Name()
		fileExt := path.Ext(temp)
		if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
			continue
		}
		fileList = append(fileList, f.Name())
	}
	return fileList
}

//GetMetadata get GetMetadata map
func GetMetadata(image *os.File) map[string]string {
	var metadata map[string]string = make(map[string]string)
	imageInfo, err := exif.Decode(image)
	if err != nil {
		log.Fatal(err)
	}
	//var ImageWidth *tiff.Tag
	jsonByte, err := imageInfo.MarshalJSON()

	if err != nil {
		log.Fatal(err.Error())
	}

	jsonString := string(jsonByte)
	metadata["model"] = gjson.Get(jsonString, "Model").String()
	metadata["datetime"] = gjson.Get(jsonString, "DateTime").String()
	metadata["LensMake"] = gjson.Get(jsonString, "LensMake").String()
	return metadata
}

//GetMetadataString get GetMetadata string
func GetMetadataString(image *os.File) string {
	imageInfo, err := exif.Decode(image)
	if err != nil {
		log.Fatal(err)
	}
	//var ImageWidth *tiff.Tag
	jsonByte, err := imageInfo.MarshalJSON()

	if err != nil {
		log.Fatal(err.Error())
	}

	jsonString := string(jsonByte)
	return jsonString
}

//GetFileNameOnly get file name
func GetFileNameOnly(filefullname string) string {
	filenameWithSuffix := path.Base(filefullname)
	fileSuffix := path.Ext(filenameWithSuffix)
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return filenameOnly
}
