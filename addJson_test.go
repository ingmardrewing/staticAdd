package main

import (
	"os"
	"path"
	"testing"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticBlogAdd"
	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticPersistence"
)

func TestNewAddJson(t *testing.T) {
	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/src/add"
	destDir := "testResources/deploy"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	aj := NewAddJson(envName, srcDir, destDir, excerpt, url)
	if bucketName != aj.awsBucket {
		t.Error("Expected", aj.awsBucket, "to be", bucketName)
	}
}
func TestGenerateDto(t *testing.T) {
	staticBlogAdd.DontUpload()
	aj := givenAddJson()

	aj.GenerateDto()
	t.Error(aj.dto)
}

func TestWriteToFs(t *testing.T) {
	aj := givenAddJson()
	aj.dto = givenPageDto()

	aj.WriteToFs()
	expected := `{
	"version":1,
	"thumbImg":"thumbUrlValue",
	"postImg":"imageUrlValue",
	"filename":"htmlfilenameValue",
	"id":42,
	"date":"createDateValue",
	"url":"urlValue",
	"title":"titleValue",
	"title_plain":"titlePlainValue",
	"excerpt":"descriptionValue",
	"content":"contentValue",
	"dsq_thread_id":"disqusIdValue",
	"thumbBase64":"thumbBase64Value",
	"category":"categoryValue"
}`
	actual := fs.ReadFileAsString(path.Join(getTestFileDirPath(), "testResources/src/posts/page42.json"))

	if actual != expected {
		t.Error("expected\n", expected, "\nbut got\n", actual)
	}
}

func givenAddJson() *addJson {
	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/src/add"
	destDir := "testResources/src/posts"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	return NewAddJson(envName, srcDir, destDir, excerpt, url)
}

func givenPageDto() staticIntf.PageDto {
	return staticPersistence.NewFilledDto(42,
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
}
