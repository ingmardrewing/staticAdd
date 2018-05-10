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
	src := path.Join(getTestFileDirPath(), "testResources/image/test-image.png")
	dest := path.Join(getTestFileDirPath(), "testResources/src/add/image.png")
	fs.CopyFile(src, dest)
	staticBlogAdd.DoUpload(false)
}

func tearDown() {
	paths := []string{
		path.Join(getTestFileDirPath(), conf[0].AddPostDir),
		path.Join(getTestFileDirPath(), conf[0].Src[0].Dir)}
	for _, p := range paths {
		fs.RemoveDirContents(p)
	}
	staticBlogAdd.DoUpload(true)
}

func getTestFileDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func TestReadConf(t *testing.T) {
	readConf()
}
