package util

import (
	"io"
	"log"
	"mime/multipart"
	"os"
)

func SaveFileLocal(data *multipart.FileHeader, path string) error {
	file, err := data.Open()
	if err != nil {
		log.Println("将视频保存到本地失败", err)
		return err
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Panicln("关闭文件时出现了错误", err)
		}
	}(file)
	out, err := os.Create(path)
	if err != nil {
		log.Println("本地文件创建失败, 创建目录为", path, err)
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(out)
	_, err = io.Copy(out, file)
	if err != nil {
		log.Println("本地文件保存失败", err)
		return err
	}
	return nil
}

func RemoveFileLocal(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Println("路径出现了问题, 推测为没有找到可以删除的视频", err)
	}
}
