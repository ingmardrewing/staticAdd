package main

import (
	"os"

	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticPersistence"
	"github.com/ingmardrewing/staticUtil"
)

func NewAddJson(bucketEnv, srcDir, destDir, excerpt, url string) *addJson {
	aj := new(addJson)
	aj.awsBucket = os.Getenv(bucketEnv)
	aj.srcDir = srcDir
	aj.destDir = destDir
	aj.excerpt = excerpt
	aj.url = url
	return aj
}

type addJson struct {
	awsBucket  string
	srcDir     string
	destDir    string
	excerpt    string
	url        string
	filename   string
	titlePlain string
	imageUrl   string
	dto        staticIntf.PageDto
	tags       []string
}

func (a *addJson) GenerateDto() {
	bda := NewBlogDataAbstractor(
		a.awsBucket,
		a.srcDir,
		a.destDir,
		a.excerpt,
		a.url)

	bda.ExtractData()
	a.dto = bda.GeneratePostDto()
	a.filename = bda.GetFilename()
	a.titlePlain = bda.GetTitlePlain()
	a.imageUrl = bda.GetImageUrl()
	a.tags = bda.GetTags()
}

func (a *addJson) GetTags() []string {
	return a.tags
}

func (a *addJson) WriteToFs() {
	staticPersistence.WritePageDtoToJson(
		a.dto,
		a.destDir,
		a.filename)
}

func (a *addJson) CurlData() (string, string, string, string) {
	url := a.url + staticUtil.GenerateDatePath() + a.titlePlain + "/"
	return a.dto.Title(), a.dto.Description(), url, a.imageUrl
}
