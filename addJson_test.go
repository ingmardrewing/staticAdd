package main

import (
	"os"
	"path"
	"testing"

	"github.com/ingmardrewing/fs"
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
	aj := givenAddJson()

	aj.GenerateDto()
	expected := `<a href=\"testResources/src/add/testImage.png\"><img src=\"testResources/src/add/testImage-w800.png\" width=\"800\"></a>`
	actual := aj.dto.Content()

	if actual != expected {
		t.Error("Expected\n", expected, "\nbut got\n", actual)
	}
}

func TestWriteToFs(t *testing.T) {
	aj := givenAddJson()
	aj.dto = givenPageDto()
	aj.WriteToFs()
	expected := `{
	"version":2,
	"filename":"htmlfilenameValue",
	"path_from_doc_root":"pathValue",
	"category":"categoryValue",
	"tags":["tag1","tag2"],
	"create_date":"createDateValue",
	"title":"titleValue",
	"title_plain":"titlePlainValue",
	"excerpt":"descriptionValue",
	"content":"contentValue",
	"thumb_base64":"thumbBase64Value",
	"images_urls":[{"title":"titleValue","w_190":"microThumbValue","w_390":"thumbUrlValue","w_800":"imageUrlValue","max_resolution":""}]
}`

	actual := fs.ReadFileAsString(path.Join(getTestFileDirPath(), "testResources/src/posts/doc00012.json"))

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
	return staticPersistence.NewFilledDto(12,
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
		"categoryValue",
		"microThumbValue",
		[]string{"tag1", "tag2"},
		[]staticIntf.Image{})
}
