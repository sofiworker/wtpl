package zip

import (
	gozip "archive/zip"
	"bytes"
)

func ZipFile(outputName string, name ...string) error {
	return zip(outputName, name...)
}

//func ZipDir() error {
//
//}
//
//func Unzip() error {
//
//}

func zip(outputName string, fileOrPath ...string) error {
	writer := gozip.NewWriter(&bytes.Buffer{})
	return writer.Copy(&gozip.File{})
}

func unzip() error {
	return nil
}
