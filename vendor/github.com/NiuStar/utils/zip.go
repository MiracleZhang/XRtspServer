package utils

import (
	"archive/zip"
	"os"
	"io"
	"io/ioutil"
	"fmt"
	"bytes"
	"path/filepath"
	"strings"
)


func CompressZip(dir, dest string) {
	//获取源文件列表
	f, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}
	fzip, _ := os.Create(dest)
	w := zip.NewWriter(fzip)
	defer func() {
		w.Close()
		fzip.Close()
	}()
	for _, file := range f {

		if !file.IsDir() {
			fw, _ := w.Create(file.Name())
			filecontent, err := ioutil.ReadFile(dir + "/" + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			n, err := fw.Write(filecontent)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(n)
		} else {
			w.Create(file.Name() + "/")
		}
	}
}

func CompressDir(dir, dest string) {

	fzip, _ := os.Create(dest)
	w := zip.NewWriter(fzip)
	defer func() {
		w.Close()
		fzip.Close()
	}()
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		fmt.Println("path:",path)
		path2 := Substr(filepath.Clean(path),len(filepath.Clean(dir)),len(filepath.Clean(path)))

		if strings.HasPrefix(path2,"/") {
			path2 = Substr(path2,1,len(path2))
		} else if strings.HasPrefix(path2,"\\") {
			path2 = Substr(path2,1,len(path2))
		}

		path2 = strings.Replace(path2,"\\","/",-1)
		fmt.Println("path2:",path2)
		if !info.IsDir() {
			//path2 = strings.Replace(path2,"\\","/")
			fw, _ := w.Create(path2)
			fi, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
			}
			defer fi.Close()

			io.Copy(fw,fi)

		} else {

			w.Create(path2 + "/")
		}

		return nil

	})
}
// 参数frm可以是文件或目录，不会给dst添加.zip扩展名
func Compress(frm, dst string) error {
	buf := bytes.NewBuffer(make([]byte, 0, 10*1024*1024)) // 创建一个读写缓冲
	myzip := zip.NewWriter(buf)              // 用压缩器包装该缓冲
	// 用Walk方法来将所有目录下的文件写入zip
	err := filepath.Walk(frm, func(path string, info os.FileInfo, err error) error {
		var file []byte
		if err != nil {
			return filepath.SkipDir
		}
		header, err := zip.FileInfoHeader(info) // 转换为zip格式的文件信息
		if err != nil {
			return filepath.SkipDir
		}
		header.Name, _ = filepath.Rel(filepath.Dir(frm), path)
		if !info.IsDir() {
			// 确定采用的压缩算法（这个是内建注册的deflate）
			//header.Method = 8
			file, err = ioutil.ReadFile(path) // 获取文件内容
			if err != nil {
				return filepath.SkipDir
			}
		} else {
			file = nil
		}
		// 上面的部分如果出错都返回filepath.SkipDir
		// 下面的部分如果出错都直接返回该错误
		// 目的是尽可能的压缩目录下的文件，同时保证zip文件格式正确
		w, err := myzip.CreateHeader(header) // 创建一条记录并写入文件信息
		if err != nil {
			return err
		}

		if !info.IsDir() {
			_, err = w.Write(file) // 非目录文件会写入数据，目录不会写入数据
			if err != nil {    // 因为目录的内容可能会修改
				return err     // 最关键的是我不知道咋获得目录文件的内容
			}
			return nil
		}
		return nil

	})
	if err != nil {
		return err
	}
	myzip.Close()        // 关闭压缩器，让压缩器缓冲中的数据写入buf
	file, err := os.Create(dst) // 建立zip文件
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = buf.WriteTo(file) // 将buf中的数据写入文件
	if err != nil {
		return err
	}
	return nil
}