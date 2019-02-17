package main

import (
	"fmt"
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
	awsBucket string
	srcDir    string
	destDir   string
	excerpt   string
	url       string
	dto       staticIntf.PageDto
	tags      []string
}

func (a *addJson) GenerateDto() {
	bda := NewBlogDataAbstractor(
		a.awsBucket,
		a.srcDir,
		a.destDir,
		a.excerpt,
		a.url)

	bda.ExtractData()
	a.tags = bda.GetTags()
	a.dto = bda.GeneratePostDto()
}

func (a *addJson) GetTags() []string {
	return a.tags
}

func (a *addJson) WriteToFs() {
	filename := fmt.Sprintf(
		staticPersistence.JsonFileNameTemplate(),
		a.dto.Id())

	staticPersistence.WritePageDtoToJson(
		a.dto,
		a.destDir,
		filename)
}

func (a *addJson) CurlData() (string, string, string, string) {
	url := a.url + staticUtil.GenerateDatePath() + a.dto.TitlePlain() + "/"
	return a.dto.Title(), a.dto.Description(), url, a.dto.ImageUrl()
}
