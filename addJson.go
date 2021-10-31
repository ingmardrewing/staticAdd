package main

import (
	"os"

	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticPersistence"
	"github.com/ingmardrewing/staticUtil"
)

func NewAddJson(deployedStaticAssetsLocation, bucketEnv, srcDir, destDir, excerpt string,
	defaultByTag []staticPersistence.DefaultByTag,
	url string) *addJson {

	//fmt.Println(fmt.Sprintf("New blog content from content %s\n", srcDir))
	//fmt.Println(fmt.Sprintf("Creating JSON in %s\n", deployedStaticAssetsLocation))

	aj := new(addJson)
	aj.deployedStaticAssetsLocation = deployedStaticAssetsLocation
	aj.awsBucket = os.Getenv(bucketEnv)
	aj.srcDir = srcDir
	aj.destDir = destDir
	aj.excerpt = excerpt
	aj.url = url
	aj.dbt = defaultByTag
	return aj
}

type addJson struct {
	deployedStaticAssetsLocation string
	awsBucket                    string
	srcDir                       string
	destDir                      string
	excerpt                      string
	url                          string
	filename                     string
	titlePlain                   string
	dto                          staticIntf.PageDto
	dbt                          []staticPersistence.DefaultByTag
	tags                         []string
}

func (a *addJson) GenerateDto() {
	bda := NewBlogDataAbstractor(
		a.deployedStaticAssetsLocation,
		a.awsBucket,
		a.srcDir,
		a.destDir,
		a.excerpt,
		a.url,
		a.dbt)

	bda.ExtractData()
	a.dto = bda.GeneratePostDto()
	a.filename = bda.GetFilename()
	a.titlePlain = bda.GetTitlePlain()
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
	return a.dto.Title(), a.dto.Description(), url, a.dto.Images()[0].MaxResolution()
}
