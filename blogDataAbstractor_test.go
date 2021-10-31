package main

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticPersistence"
	"github.com/ingmardrewing/staticUtil"
)

func TestGenerateDatePath(t *testing.T) {
	actual := staticUtil.GenerateDatePath()
	now := time.Now()
	expected := fmt.Sprintf("%d/%d/%d/", now.Year(), now.Month(), now.Day())

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestBlogDataAbstractor(t *testing.T) {
	bda := givenBlogDataAbstractor()
	bda.im = &imgManagerMock{}
	bda.ExtractData()
	dto := bda.GeneratePostDto()

	actual := dto.Title()
	expected := "Test Image"

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}

	actual = dto.Content()
	expected = `<a href=\"https://drewing.de/just/another/path/TestImage.png\"><img src=\"https://drewing.de/just/another/path/TestImage-w800.png\" srcset=\"https://drewing.de/just/another/path/TestImage-w1600.png 2x\" width=\"800\" alt=\"Test Image\"></a>`
	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}

	actual = dto.Description()
	expected = "A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing."

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}

	actual = dto.CreateDate()
	n := time.Now()
	expected = fmt.Sprintf("%d-%d-%d", n.Year(), n.Month(), n.Day())

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestSplitAtSpecialChars(t *testing.T) {
	expected := []string{"a", "b", "c", "d"}
	actual := splitAtSpecialChars("a-b,c_d")

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestExtractTagsFromMarkdownText(t *testing.T) {
	bda := givenBlogDataAbstractor()
	md := `# this is headline not a tag
#butthis and #this is`

	actual := bda.extractTags(md)
	expected := []string{"butthis", "this"}

	for i := 0; i <= 1; i++ {
		if actual[i] != expected[i] {
			t.Error("Expected", expected[i], "but got", actual[i], "at index", i)
		}
	}
}

func TestFindFileNameTags(t *testing.T) {
	bda := givenBlogDataAbstractor()

	givenFilename := "filename+one+two.png"
	expectedPureFilename := "filename.png"
	actualPureFilename, actualTags := bda.findFileNameTags(givenFilename)
	if actualPureFilename != expectedPureFilename {
		t.Error("Expected", expectedPureFilename, "but got", actualPureFilename)
	}

	if len(actualTags) != 2 {
		t.Error("Expected actualTags to be of length 2, but it isn't.")
	}

	if actualTags[0] != "one" {
		t.Error("Expected first found tag to be 'one', but it is:", actualTags[0])
	}

	if actualTags[1] != "two" {
		t.Error("Expected second found tag to be 'two', but it is:", actualTags[1])
	}
}

func TestSplitCamelCaseAndNumbers(t *testing.T) {
	expected := []string{"another", "AOC", "Test", "4", "this"}
	actual := splitCamelCaseAndNumbers("anotherAOCTest4this")

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestFindUpperCaseSequence(t *testing.T) {
	expected := []string{"AOC", "Test"}
	actual := findUpperCaseSequence("AOCTest")

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestInferBlogTitleFromFilename(t *testing.T) {
	bda := givenBlogDataAbstractor()

	filename2expected := map[string]string{
		"iPadTest.png":       "I Pad Test",
		"this-is-a-test.png": "This Is A Test",
		"BeeTwoPointO.png":   "Bee Two Point O",
		"test_image.png":     "Test Image",
		"even4me.png":        "Even 4 Me"}
	for filename, expected := range filename2expected {
		actual, _ := bda.inferBlogTitleFromFilename(filename)
		if actual != expected {
			t.Error("Expected", expected, "but got", actual)
		}
	}

}

func TestWriteData(t *testing.T) {
	staticAssetsLoc := "/static-assets/blog"
	addDir := getTestFileDirPath() + "/testResources/src/add/"
	postsDir := getTestFileDirPath() + "/testResources/src/posts/"
	dExcerpt := "A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing."
	domain := "https://drewing.de/blog/"

	bda := NewBlogDataAbstractor(staticAssetsLoc, "drewingde", addDir, postsDir, dExcerpt, domain, nil)
	bda.im = &imgManagerMock{}
	bda.ExtractData()
	dto := bda.GeneratePostDto()

	filename := fmt.Sprintf("page%d.json", bda.data.id)

	staticPersistence.WritePageDtoToJson(dto, postsDir, filename)

	data := fs.ReadFileAsString(path.Join(postsDir, filename))
	path := "/" + staticUtil.GenerateDatePath() + "test-image/"

	date := staticUtil.GetDate()
	expected := `{
	"version":2,
	"filename":"index.html",
	"path_from_doc_root":"` + path + `",
	"category":"blog post",
	"tags":[],
	"create_date":"` + date + `",
	"title":"Test Image",
	"excerpt":"A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing.",
	"content":"<a href=\"https://drewing.de/just/another/path/TestImage.png\"><img src=\"https://drewing.de/just/another/path/TestImage-w800.png\" srcset=\"https://drewing.de/just/another/path/TestImage-w1600.png 2x\" width=\"800\" alt=\"Test Image\"></a>",
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
		"w_1600_portrait":"https://drewing.de/just/another/path/TestImage-w800.png",
		"max_resolution":"https://drewing.de/just/another/path/TestImage.png"
	}]
}`

	if data != expected {
		t.Error("Expected", expected, "but got", data)
	}
}

type imgManagerMock struct{}

func (i *imgManagerMock) PrepareImages() {}
func (i *imgManagerMock) UploadImages()  {}
func (i *imgManagerMock) GetImageUrls() []string {
	return []string{
		"https://drewing.de/just/another/path/TestImage-w80-square.png",
		"https://drewing.de/just/another/path/TestImage-w100-square.png",
		"https://drewing.de/just/another/path/TestImage-w190-square.png",
		"https://drewing.de/just/another/path/TestImage-w200-square.png",
		"https://drewing.de/just/another/path/TestImage-w390-square.png",
		"https://drewing.de/just/another/path/TestImage-w400-square.png",
		"https://drewing.de/just/another/path/TestImage-w800-square.png",
		"https://drewing.de/just/another/path/TestImage-w800.png",
		"https://drewing.de/just/another/path/TestImage-w800.png",
		"https://drewing.de/just/another/path/TestImage-w1600.png",
		"https://drewing.de/just/another/path/TestImage.png"}
}
func (i *imgManagerMock) AddImageSize(size int) string {
	return "TestImage-w" + strconv.Itoa(size) + ".png"
}
func (i *imgManagerMock) AddCropImageSize(size int) string {
	return "TestImage-w" + strconv.Itoa(size) + ".png"
}

func givenBlogDataAbstractor() *BlogDataAbstractor {
	addDir := getTestFileDirPath() + "/testResources/src/add/"
	postsDir := getTestFileDirPath() + "/testResources/src/posts/"
	dExcerpt := "A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing."
	staticAssetsLoc := "/static-assets/blog"

	jsonString := `{"tag":"sketch","excerpt":"excerpt","content":"content"}`
	var dbt staticPersistence.DefaultByTag
	json.Unmarshal([]byte(jsonString), &dbt)

	dbts := []staticPersistence.DefaultByTag{dbt}

	return NewBlogDataAbstractor(staticAssetsLoc,
		"drewingde",
		addDir,
		postsDir,
		dExcerpt,
		"https://drewing.de/blog/",
		dbts)
}
