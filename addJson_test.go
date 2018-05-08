package main

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticPersistence"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
	tearDown()
}

func setup() {
	conf = staticPersistence.ReadConfig("testResources/", "configNew.json")
}

func tearDown() {
	filepath := path.Join(getTestFileDirPath(), conf[0].Deploy.TargetDir)
	fs.RemoveDirContents(filepath)
}

func getTestFileDirPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func TestNewAddJson(t *testing.T) {
	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/add"
	destDir := "testResources/deploy"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	aj := NewAddJson(envName, srcDir, destDir, excerpt, url)
	if bucketName != aj.awsBucket {
		t.Error("Expected", aj.awsBucket, "to be", bucketName)
	}
}

func TestWriteToFs(t *testing.T) {
	dto := staticPersistence.NewFilledDto(42,
		"titleValue",
		"titlePlainValue",
		"thumbUrlValue",
		"imageUrlValue",
		"descriptionValue",
		"disqusIdValue",
		"createDateValue",
		"contentValue",
		"urlValue",
		"domainValue",
		"pathValue",
		"fspathValue",
		"htmlfilenameValue",
		"thumbBase64Value",
		"categoryValue")

	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/add"
	destDir := "testResources/deploy"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	aj := NewAddJson(envName, srcDir, destDir, excerpt, url)
	aj.dto = dto
	aj.WriteToFs()

	ba := fs.ReadByteArrayFromFile("testResources/deploy/page42.json")

	actual := len(ba)
	expected := 375

	if actual != expected {
		t.Error("Expected byte array to be of length", expected, "but it was", actual)
	}

}
