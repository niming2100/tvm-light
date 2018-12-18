package util

import (
	"io"
	"net/http"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func FileDownLoad(localPath string,fileName string, urlPath string) error {
	res, err := http.Get(urlPath);
	if err != nil {
		return err;
	}
	os.MkdirAll(localPath,os.ModePerm)
	f, err := os.Create(localPath+fileName)

	if err != nil {
		return err;
	}
	io.Copy(f, res.Body)
	return nil;
}
