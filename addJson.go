package main

import (
	"fmt"
	"os"

	"github.com/ingmardrewing/staticBlogAdd"
	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticPersistence"
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
}

func (a *addJson) GenerateDto() {
	bda := staticBlogAdd.NewBlogDataAbstractor(a.awsBucket, a.srcDir, a.destDir, a.excerpt, a.url)
	a.dto = bda.GeneratePostDto()
}

func (a *addJson) WriteToFs() {
	filename := fmt.Sprintf("page%d.json", a.dto.Id())
	staticPersistence.WritePageDtoToJson(a.dto, a.destDir, filename)
}

func (a *addJson) CurlData() (string, string, string, string) {
	return a.dto.Title(), a.dto.Description(), a.dto.Url(), a.dto.ImageUrl()
}
