package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//Myapp start the app
func Myapp() {
	fmt.Println("start")

	input := bufio.NewScanner(os.Stdin)

	client := CreateClient()
	var bucketName = "myron--bucket"
	var directory string = ""
	var filename string = ""
	var role string = "N"
	var ready string = ""
	var toexit string = ""
	var photoList []string = make([]string, 0)
	for 1 == 1 {

		for 1 == 1 {
			fmt.Println("******************************************************************************")
			fmt.Println(" + Please Enter the photo directory(not file path) you want to use.")
			fmt.Println(" + Exit (Enter exit)")
			fmt.Println("******************************************************************************")
			fmt.Print(">>")

			for input.Scan() {
				if len(input.Text()) != 0 {
					break
				}
			}
			directory = input.Text()
			directory = strings.Replace(directory, "\\ ", " ", -1)

			if directory == "exit" {
				os.Exit(3)
			}
			if !IsDir(directory) {
				fmt.Println("'" + directory + "' is not a useful directory！ Make sure there are no spaces at the beginning and the ending")
				directory = ""
				continue
			}

			photoList = ReadDir(directory)
			if len(photoList) == 0 {
				fmt.Println("'" + directory + "' has no image file")
				directory = ""
				continue
			}
			fmt.Println("The Photo List:")
			for _, str := range photoList {
				fmt.Println("- " + str)
			}
			directory = directory + "/"
			break
		}

		for 1 == 1 {
			fmt.Println("******************************************************************************")
			fmt.Println(" + Do You want to upload one photo or photo list?")
			fmt.Println(" - A. One Photo")
			fmt.Println(" - B. Photo List")
			fmt.Println(" - E. Exit")
			fmt.Println("******************************************************************************")
			fmt.Print(">>")
			for input.Scan() {
				if len(input.Text()) != 0 {
					break
				}
			}
			role = input.Text()
			if role == "E" {
				os.Exit(3)
			} else if role == "A" {
				break
			} else if role == "B" {
				break
			} else {
				fmt.Println("Please Enter 'A', 'B' or 'E'")
			}
		}

		if role == "A" {
			for 1 == 1 {
				fmt.Println("******************************************************************************")
				fmt.Println(" + Please Enter the photo name (ex: mytest.jpg)")
				fmt.Println("******************************************************************************")
				fmt.Print(">>")
				for input.Scan() {
					if len(input.Text()) != 0 {
						break
					}
				}
				filename = input.Text()
				if !IsFile(directory + filename) {
					fmt.Println("It's not a file !!")
					continue
				}
				fmt.Println("******************************************************************************")
				fmt.Println(" + Do You Want to upload it?(Y/N)")
				fmt.Println("******************************************************************************")
				fmt.Print(">>")
				for input.Scan() {
					if len(input.Text()) != 0 {
						break
					}
				}
				ready = input.Text()
				if ready == "Y" {
					UploadImage(client, bucketName, directory, filename)
					UploadExifTxt(client, bucketName, directory, filename)
					break
				} else {
					continue
				}
			}

		} else if role == "B" {
			for 1 == 1 {
				fmt.Println("******************************************************************************")
				fmt.Println(" + Do You Want to upload All Photos?(Y/N)")
				fmt.Println("******************************************************************************")
				fmt.Print(">>")
				for input.Scan() {
					if len(input.Text()) != 0 {
						break
					}
				}
				ready = input.Text()
				if ready == "Y" {
					UploadImageList(client, bucketName, directory, photoList)
					break
				} else {
					continue
				}
			}
		}

		fmt.Println("******************************************************************************")
		fmt.Println(" + Do You Want to Exit (Y/N)")
		fmt.Println("******************************************************************************")
		fmt.Print(">>")
		for input.Scan() {
			if len(input.Text()) != 0 {
				break
			}
		}
		toexit = input.Text()

		if toexit == "Y" {
			os.Exit(3)
		} else {
			directory = ""
			filename = ""
			role = "N"
			ready = ""
			toexit = ""
			continue

		}
	}

}
