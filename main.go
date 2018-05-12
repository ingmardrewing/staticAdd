package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ingmardrewing/fs"
	curl "github.com/ingmardrewing/gomicSocMedCurl"
	"github.com/ingmardrewing/staticPersistence"
)

var (
	conf        []staticPersistence.JsonConfig
	configFile  = "configNew.json"
	fconfigPath = ""
)

func readConf() {
	flag.StringVar(&fconfigPath, "configPath", os.Getenv("BLOG_CONFIG_DIR"), "path to config file")

	exists, _ := fs.PathExists(path.Join(fconfigPath, configFile))
	if exists {
		conf = staticPersistence.ReadConfig(fconfigPath, configFile)
	} else {
		conf = staticPersistence.ReadConfig("./testResources/", configFile)
	}
}

func main() {
	readConf()

	aj := NewAddJson("AWS_BUCKET", conf[0].AddPostDir, conf[0].WritePostDir, conf[0].DefaultMeta.BlogExcerpt, "https://drewing.de/blog/")
	aj.GenerateDto()
	aj.WriteToFs()
	fs.RemoveDirContents(conf[0].AddPostDir)

	fmt.Println("About to clear add post dir", conf[0].AddPostDir)
	fs.RemoveDirContents(conf[0].AddPostDir)

	title, desc, link, imgUrl := aj.CurlData()
	tagsCsv := strings.Join(curl.TAGS, ",")
	json := curl.Json(title, desc, link, imgUrl, tagsCsv)
	fmt.Println(json)
}
