package main

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticBlogAdd"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {
	os.Setenv("BLOG_CONFIG_DIR", "./testResources/")
	for _, p := range givenDirPaths() {
		fs.CreateDir(p)
	}
	src := path.Join(getTestFileDirPath(), "testResources/image/test-image.png")
	dest := path.Join(getTestFileDirPath(), "testResources/src/add/image.png")
	fs.CopyFile(src, dest)
	staticBlogAdd.DoUpload(false)

}

func tearDown() {
	for _, p := range givenDirPaths() {
		fs.RemoveDirContents(p)
	}
	staticBlogAdd.DoUpload(true)
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
