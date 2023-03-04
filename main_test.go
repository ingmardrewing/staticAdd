package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ingmardrewing/fs"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	os.Setenv("BLOG_UPLOAD_SSH_USER", "www.drewing.de")
	os.Setenv("BLOG_UPLOAD_SSH_PASS", "F0tmmctddacowmaebod2sopeg!")
	os.Setenv("BLOG_UPLOAD_SSH_SERVER", "ssh.strato.de")
	os.Setenv("BLOG_UPLOAD_SSH_PORT", "22")
	os.Setenv("BLOG_CONFIG_DIR", "./testResources/")
	for _, p := range givenDirPaths() {
		fs.CreateDir(p)
	}
	src := path.Join(getTestFileDirPath(), "testResources/image/test-image.png")
	dest := path.Join(getTestFileDirPath(), "testResources/src/add/testImage.png")
	fs.CopyFile(src, dest)
	DoUpload(false)
}

func tearDown() {
	for _, p := range givenDirPaths() {
		fmt.Print("Removing content in: " + p + "\n")
		fs.RemoveDirContents(p)
	}

	pth := path.Join(getTestFileDirPath(), "testResources/src/posts/")
	filename := "page42.json"
	if exist, _ := fs.PathExists(path.Join(pth, filename)); exist == true {
		fs.RemoveFile(pth, filename)
	}
	DoUpload(true)
}

func givenDirPaths() []string {
	return []string{
		path.Join(getTestFileDirPath(), "testResources/src/add/"),
		path.Join(getTestFileDirPath(), "testResources/src/posts/")}
}

func getTestFileDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func TestReadConf(t *testing.T) {
	readConf()
}
