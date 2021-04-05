package pkg

import "testing"

/*func TestUploadImage(t *testing.T) {
	client := CreateClient()
	var dirctory string = "./asset/photo/"
	var filename string = "mytest1.jpg"
	actual := UploadImage(client, "myron--bucket", dirctory, filename)
	if !actual {
		t.Errorf("UploadImage() failed dirctory=%s, filename=%s", dirctory, filename)
	}
}*/

func TestUploadImageList(t *testing.T) {
	client := CreateClient()
	var dirctory string = "./asset/photo/"
	fileList := ReadDir(dirctory)
	actual := UploadImageList(client, "myron--bucket", dirctory, fileList)
	if !actual {
		t.Errorf("UploadImageList() failed dirctory=%s", dirctory)
	}
}
