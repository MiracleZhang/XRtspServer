package utils

import (
	
	"os"
	"io"
	"fmt"
	"path/filepath"
)

func copyDir(src,des string) {

	exit,_ := PathExists(des)

	if !exit {
		Mkdir(des)
	}


	filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		//fmt.Println("path:",path)
		//fmt.Println("path2:",filepath.Clean(path))
		path2 := Substr(filepath.Clean(path),len(filepath.Clean(src)),len(filepath.Clean(path)))

		fmt.Println("path3:",path2)

		if !info.IsDir() {
			copyFile(path,des + path2)
/*
			filecontent, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err)
			}
			err = WriteToFile(des + path2,string(filecontent))
			if err != nil {
				fmt.Println(err)
			}*/
		} else {
			Mkdir(des + path2)
		}

		return nil

	})
}

func copyFile(src,des string) {

	f, err := os.OpenFile(des, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return
	}

	defer f.Close()

	fi, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer fi.Close()

	_,err = io.Copy(f,fi)
	if err != nil {
		panic(err)
	}
}

func Copy(src,des string) {

	fmt.Println("copy src:",src)
	f, err := os.Stat(src)
	if err != nil {
		fmt.Println("Copy err:",err)
		return
	}

	if f.IsDir() {
		fmt.Println("copyDir des:",des)
		copyDir(src,des)
	} else {
		fmt.Println("copyfile des:",des)
		copyFile(src,des)
	}

}