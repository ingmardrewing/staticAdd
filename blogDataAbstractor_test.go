package main

import (
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
	expected = `<a href=\"https://drewing.de/just/another/path/TestImage.png\"><img src=\"https://drewing.de/just/another/path/TestImage-w800.png\" width=\"800\"></a>`
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
	expected = fmt.Sprintf("%d-%d-%d %d:%d:%d", n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())

	if actual != expected {
		t.Error("Expected", expected, "but got", actual)
	}

	actualInt := dto.Id()
	expectedInt := 13

	if actualInt != expectedInt {
		t.Errorf("Expected %d, but got %d\n", expectedInt, actualInt)
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

func TestSplitCamelCaseAndNumbers(t *testing.T) {
	expected := []string{"another", "Test", "4", "this"}
	actual := splitCamelCaseAndNumbers("anotherTest4this")

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
	addDir := getTestFileDirPath() + "/testResources/src/add/"
	postsDir := getTestFileDirPath() + "/testResources/src/posts/"
	dExcerpt := "A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing."
	domain := "https://drewing.de/blog/"

	bda := NewBlogDataAbstractor("drewingde", addDir, postsDir, dExcerpt, domain)
	bda.im = &imgManagerMock{}
	bda.ExtractData()
	dto := bda.GeneratePostDto()

	filename := fmt.Sprintf("page%d.json", dto.Id())

	staticPersistence.WritePageDtoToJson(dto, postsDir, filename)

	data := fs.ReadFileAsString(path.Join(postsDir, filename))

	date := staticUtil.GetDate()
	expected := `{
	"version":2,
	"filename":"index.html",
	"path_from_doc_root":"",
	"category":"blog post",
	"tags":[],
	"create_date":"` + date + `",
	"title":"Test Image",
	"title_plain":"test-image",
	"excerpt":"A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing.",
	"content":"<a href=\"https://drewing.de/just/another/path/TestImage.png\"><img src=\"https://drewing.de/just/another/path/TestImage-w800.png\" width=\"800\"></a>",
	"thumb_base64":"",
	"images_urls":[{"title":"Test Image","w_190":"https://drewing.de/just/another/path/TestImage-w190.png","w_390":"https://drewing.de/just/another/path/TestImage-w390.png","w_800":"https://drewing.de/just/another/path/TestImage-w800.png","max_resolution":""}]
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
		"https://drewing.de/just/another/path/TestImage-w190.png",
		"https://drewing.de/just/another/path/TestImage-w390.png",
		"https://drewing.de/just/another/path/TestImage-w800.png",
		"https://drewing.de/just/another/path/TestImage.png"}
}
func (i *imgManagerMock) AddImageSize(size int) string {
	return "TestImage-w" + strconv.Itoa(size) + ".png"
}

func givenBlogDataAbstractor() *BlogDataAbstractor {
	addDir := getTestFileDirPath() + "/testResources/src/add/"
	postsDir := getTestFileDirPath() + "/testResources/src/posts/"
	dExcerpt := "A blog containing texts, drawings, graphic narratives/novels and (rarely) code snippets by Ingmar Drewing."

	return NewBlogDataAbstractor("drewingde",
		addDir,
		postsDir,
		dExcerpt,
		"https://drewing.de/blog/")
}
