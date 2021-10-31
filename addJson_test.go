package main

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticPersistence"
)

func TestNewAddJson(t *testing.T) {
	staticAssetsLoc := "/static-assets/blog"
	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/src/add"
	destDir := "testResources/deploy"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	aj := NewAddJson(staticAssetsLoc, envName, srcDir, destDir, excerpt, nil, url)
	if bucketName != aj.awsBucket {
		t.Error("Expected", aj.awsBucket, "to be", bucketName)
	}
}

func TestGenerateDto(t *testing.T) {
	aj := givenAddJson()

	aj.GenerateDto()

	expected := `<a href=\"testResources/src/add/testImage.png\"><img src=\"testResources/src/add/testImage-w800.png\" srcset=\"testResources/src/add/testImage-w1600.png 2x\" width=\"800\" alt=\"Test Image\"></a>`
	actual := aj.dto.Content()

	if actual != expected {
		t.Error("Expected\n", expected, "\nbut got\n", actual)
	}
}

func TestWriteToFs(t *testing.T) {
	aj := givenAddJson()
	aj.dto = givenPageDto()
	aj.filename = "doc00012.json"
	aj.WriteToFs()
	expected := `{
	"version":2,
	"filename":"htmlfilenameValue",
	"path_from_doc_root":"pathValue",
	"category":"categoryValue",
	"tags":["tag1","tag2"],
	"create_date":"createDateValue",
	"title":"titleValue",
	"excerpt":"descriptionValue",
	"content":"contentValue",
	"images_urls":[{
		"title":"Test Image",
		"w_85":"https://drewing.de/just/another/path/TestImage-w80-square.png",
		"w_100":"https://drewing.de/just/another/path/TestImage-w100-square.png",
		"w_190":"https://drewing.de/just/another/path/TestImage-w190-square.png",
		"w_200":"https://drewing.de/just/another/path/TestImage-w200-square.png",
		"w_390":"https://drewing.de/just/another/path/TestImage-w390-square.png",
		"w_400":"https://drewing.de/just/another/path/TestImage-w400-square.png",
		"w_800":"https://drewing.de/just/another/path/TestImage-w800-square.png",
		"w_800_portrait":"https://drewing.de/just/another/path/TestImage-w800.png",
		"w_1600_portrait":"https://drewing.de/just/another/path/TestImage-w1600.png",
		"max_resolution":"https://drewing.de/just/another/path/TestImage.png"
	}]
}`

	actual := fs.ReadFileAsString(path.Join(getTestFileDirPath(), "testResources/src/posts/doc00012.json"))

	if actual != expected {
		t.Error("expected\n", expected, "\nbut got\n", actual)
	}
}

func TestCurlData(t *testing.T) {
	aj := givenAddJson()
	aj.GenerateDto()

	expectedTitle := "Test Image"
	expectedDescription := "Test 1, 2"

	n := time.Now()
	expectedUrl := fmt.Sprintf("https://drewing.de/blog%d/%d/%d/test-image/", n.Year(), n.Month(), n.Day())
	expectedImageUrl := "testResources/src/add/testImage.png"

	title, description, url, imageUrl := aj.CurlData()

	if title != expectedTitle {
		t.Error("expected", expectedTitle, "but got", title)
	}
	if description != expectedDescription {
		t.Error("expected", expectedDescription, "but got", description)
	}
	if url != expectedUrl {
		t.Error("expected", expectedUrl, "but got", url)
	}
	if imageUrl != expectedImageUrl {
		t.Error("expected", expectedImageUrl, "but got", imageUrl)
	}
}

func givenAddJson() *addJson {
	staticAssetsLoc := "/static-assets/blog"
	envName := "TEST_AWS_BUCKET"
	bucketName := "testBucketName"
	srcDir := "testResources/src/add"
	destDir := "testResources/src/posts"
	excerpt := "Test 1, 2"
	url := "https://drewing.de/blog"
	os.Setenv(envName, bucketName)

	return NewAddJson(staticAssetsLoc, envName, srcDir, destDir, excerpt, nil, url)
}

func givenPageDto() staticIntf.PageDto {
	img := givenImage()
	return staticPersistence.NewPageDto(
		"titleValue",
		"descriptionValue",
		"contentValue",
		"categoryValue",
		"createDateValue",
		"pathValue",
		"htmlfilenameValue",
		[]string{"tag1", "tag2"},
		[]staticIntf.Image{img})
}

func givenImage() staticIntf.Image {
	return staticPersistence.NewImageDto(
		"Test Image",
		"https://drewing.de/just/another/path/TestImage-w80-square.png",
		"https://drewing.de/just/another/path/TestImage-w100-square.png",
		"https://drewing.de/just/another/path/TestImage-w190-square.png",
		"https://drewing.de/just/another/path/TestImage-w200-square.png",
		"https://drewing.de/just/another/path/TestImage-w390-square.png",
		"https://drewing.de/just/another/path/TestImage-w400-square.png",
		"https://drewing.de/just/another/path/TestImage-w800-square.png",
		"https://drewing.de/just/another/path/TestImage-w800.png",
		"https://drewing.de/just/another/path/TestImage-w1600.png",
		"https://drewing.de/just/another/path/TestImage.png")
}
