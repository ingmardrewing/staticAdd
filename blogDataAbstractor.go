package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/staticIntf"
	"github.com/ingmardrewing/staticPersistence"
	"github.com/ingmardrewing/staticUtil"

	"gopkg.in/russross/blackfriday.v2"
)

var (
	RX     = regexp.MustCompile("([0-9]+|[A-ZÄÜÖ]*[a-zäüöß]*)")
	RX2REL = regexp.MustCompile("[A-ZÄÜÖ][A-ZÄÜÖ]+")
	RX2    = regexp.MustCompile("([A-ZÄÜÖ]+)([A-ZÄÜÖ][a-zäüöß]+)")
)

func NewBlogDataAbstractor(deployedStaticAssetsLocation, bucket, addDir, postsDir, defaultExcerpt, domain string, dbt []staticPersistence.DefaultByTag) *BlogDataAbstractor {
	bda := new(BlogDataAbstractor)
	bda.deployedStaticAssetsLocation = deployedStaticAssetsLocation
	bda.addDir = addDir
	bda.postsDir = postsDir
	bda.defaultExcerpt = defaultExcerpt
	bda.defaultByTag = dbt
	bda.domain = domain
	bda.data = new(abstractData)

	bda.data.imageFileName = bda.findImageFileInAddDir()

	imageFileNameWithoutTags, fileNameTags := bda.findFileNameTags(bda.data.imageFileName)
	bda.data.imageFileNameWithoutTags = imageFileNameWithoutTags
	bda.data.fileNameTags = fileNameTags
	bda.cleanseFileName()

	bda.data.imageFileName = bda.data.imageFileNameWithoutTags

	imgPath := path.Join(addDir, bda.data.imageFileName)
	bda.im = NewImageManager(deployedStaticAssetsLocation, bucket, imgPath)

	return bda
}

type abstractData struct {
	id                       int
	htmlFilename             string
	imageFileName            string
	imageFileNameWithoutTags string
	title                    string
	titlePlain               string
	microThumbUrl            string
	thumbUrl                 string
	imgUrl                   string
	mdContent                string
	excerpt                  string
	tags                     []string
	fileNameTags             []string
	url                      string
	path                     string
	disqId                   string
	content                  string
	date                     string
	category                 string
	images                   []staticIntf.Image
}

type BlogDataAbstractor struct {
	deployedStaticAssetsLocation string
	data                         *abstractData
	domain                       string
	addDir                       string
	postsDir                     string
	defaultExcerpt               string
	im                           ImgManager
	dto                          *staticIntf.PageDto
	defaultByTag                 []staticPersistence.DefaultByTag
}

func (b *BlogDataAbstractor) ExtractData() {
	b.data.htmlFilename = "index.html"

	title, titlePlain := b.inferBlogTitleFromFilename(b.data.imageFileName)
	b.data.title = title
	b.data.titlePlain = titlePlain

	imgUrls := b.prepareImages()

	b.data.images = append(b.data.images,
		staticPersistence.NewImageDto(
			title,
			imgUrls[0],
			imgUrls[1],
			imgUrls[2],
			imgUrls[3],
			imgUrls[4],
			imgUrls[5],
			imgUrls[6],
			imgUrls[7],
			imgUrls[8],
			imgUrls[9],
			imgUrls[10]))

	b.data.microThumbUrl = imgUrls[3]
	b.data.thumbUrl = imgUrls[7]
	b.data.imgUrl = imgUrls[10]

	mdContent, excerpt, tags := b.readMdData()
	b.data.mdContent = mdContent
	b.data.excerpt = excerpt
	b.data.tags = tags

	tpl := `<a href=\"%s\"><img src=\"%s\" srcset=\"%s 2x\" width=\"800\" alt=\"%s\"></a>%s`
	b.data.content = fmt.Sprintf(
		tpl,
		imgUrls[10],
		imgUrls[8],
		imgUrls[7],
		title,
		mdContent)

	b.data.url = b.generateUrl(titlePlain)
	b.data.path = b.generatePath(titlePlain)
	b.data.id = b.getId()
	b.data.date = staticUtil.GetDate()
	b.data.category = "blog post"
}

func (b *BlogDataAbstractor) cleanseFileName() {
	if len(b.data.imageFileNameWithoutTags) > 0 {
		from := b.addDir + b.data.imageFileName
		to := b.addDir + b.data.imageFileNameWithoutTags
		if from != to {
			if os.Rename(from, to) != nil {
				panic("Could not rename file")
			}
		}
	}
}

func (b *BlogDataAbstractor) GeneratePostDto() staticIntf.PageDto {
	return staticPersistence.NewPageDto(
		b.data.title,
		b.data.excerpt,
		b.data.content,
		b.data.category,
		b.data.date,
		b.data.path,
		b.data.htmlFilename,
		b.data.tags,
		b.data.images)
}

func (b *BlogDataAbstractor) GetTags() []string {
	return append(b.data.fileNameTags, b.data.tags...)
}

func (b *BlogDataAbstractor) GetFileNameTags() []string {
	return b.data.fileNameTags
}

func (b *BlogDataAbstractor) GetTitlePlain() string {
	return b.data.titlePlain
}

func (b *BlogDataAbstractor) GetImageUrl() string {
	return b.data.imageFileName
}

func (b *BlogDataAbstractor) GetFilename() string {
	return fmt.Sprintf(
		staticPersistence.JsonFileNameTemplate(),
		b.data.id)
}

func (b *BlogDataAbstractor) generatePath(titlePlain string) string {
	return "/" + staticUtil.GenerateDatePath() + titlePlain + "/"
}
func (b *BlogDataAbstractor) generateUrl(titlePlain string) string {
	return b.domain + staticUtil.GenerateDatePath() + titlePlain + "/"
}

func (b *BlogDataAbstractor) getId() int {
	postJsons := fs.ReadDirEntries(b.postsDir, false)
	if len(postJsons) == 0 {
		return 0
	}
	sort.Strings(postJsons)
	lastFile := postJsons[len(postJsons)-1]
	rx := regexp.MustCompile("(\\d+)")
	m := rx.FindStringSubmatch(lastFile)
	i, _ := strconv.Atoi(m[1])
	i++
	return i
}

func (b *BlogDataAbstractor) stripLinksAndImages(text string) string {
	rx := regexp.MustCompile(`\[.*\]\(.*\)`)
	return rx.ReplaceAllString(text, "")
}

func (b *BlogDataAbstractor) prepareImages() []string {
	b.im.AddCropImageSize(80)
	b.im.AddCropImageSize(100)
	b.im.AddCropImageSize(190)
	b.im.AddCropImageSize(200)
	b.im.AddCropImageSize(390)
	b.im.AddCropImageSize(400)
	b.im.AddCropImageSize(800)
	b.im.AddCropImageSize(1600)

	b.im.AddImageSize(800)
	b.im.AddImageSize(1600)

	b.im.PrepareImages()
	b.im.UploadImages()

	return b.im.GetImageUrls()
}

func (b *BlogDataAbstractor) generateExcerpt(text string) string {
	text = b.stripLinksAndImages(text)
	if len(text) > 155 {
		txt := fmt.Sprintf("%.155s ...", text)
		return b.stripQuotes(b.stripNewlines(txt))
	} else if len(text) == 0 {
		return b.defaultExcerpt
	}
	txt := strings.TrimSuffix(text, "\n")
	return b.stripQuotes(b.stripNewlines(txt))
}

func (b *BlogDataAbstractor) stripNewlines(text string) string {
	return strings.Replace(text, "\n", " ", -1)
}

func (b *BlogDataAbstractor) generateHtmlFromMarkdown(input string) string {
	bytes := []byte(input)
	htmlBytes := blackfriday.Run(bytes, blackfriday.WithNoExtensions())
	htmlString := string(htmlBytes)
	trimmed := strings.TrimSuffix(htmlString, "\n")
	escaped := b.stripQuotes(trimmed)
	return strings.Replace(escaped, "\n", " ", -1)
}

// extracts social media hashtags from the given input
// and returns them as a slice of strings without the leading #
func (b *BlogDataAbstractor) extractTags(input string) []string {
	rx := regexp.MustCompile(`#[A-Za-zäüößÄÜÖ]+\b`)
	matches := rx.FindAllString(input, -1)
	resultSet := []string{}
	for _, m := range matches {
		trimmedTag := strings.TrimPrefix(m, "#")
		resultSet = append(resultSet, trimmedTag)
	}
	return resultSet
}

func (b *BlogDataAbstractor) stripQuotes(txt string) string {
	txt = strings.Replace(txt, `'`, `’`, -1)
	return strings.Replace(txt, `"`, `\"`, -1)
}

func (b *BlogDataAbstractor) readMdData() (string, string, []string) {
	pathToMdFile := b.findMdFileInAddDir()
	if len(pathToMdFile) > 0 {
		mdData := fs.ReadFileAsString(pathToMdFile)
		excerpt := b.generateExcerpt(mdData)
		content := b.generateHtmlFromMarkdown(mdData)
		tags := b.extractTags(mdData)
		return content, excerpt, tags
	} else if len(b.data.fileNameTags) > 0 && b.defaultByTag != nil {
		for _, d := range b.defaultByTag {
			if d.Tag == b.data.fileNameTags[0] {
				return d.Content, d.Excerpt, b.data.fileNameTags
			}
		}
	}
	return "", b.defaultExcerpt, []string{}
}

func (b *BlogDataAbstractor) findImageFileInAddDir() string {
	imgs := fs.ReadDirEntriesEndingWith(b.addDir, "PNG", "png", "jpg", "jpeg", "JPG", "JPEG")
	for _, i := range imgs {
		if !strings.Contains(i, "-w") {
			return i
		}
	}
	return ""
}

func (b *BlogDataAbstractor) findFileNameTags(filename string) (string, []string) {
	ext := filepath.Ext(filename)
	basename := strings.TrimSuffix(filename, ext)
	nameParts := strings.Split(basename, "+")
	tags := []string{}
	if len(nameParts) > 1 {
		for _, p := range nameParts[1:] {
			tags = append(tags, p)
		}
	}
	pureFilename := fmt.Sprintf("%s%s", nameParts[0], ext)
	return pureFilename, tags
}

func (b *BlogDataAbstractor) inferBlogTitleFromFilename(filename string) (string, string) {
	fname := strings.TrimSuffix(filename, filepath.Ext(filename))
	return b.inferBlogTitle(fname), b.inferBlogTitlePlain(fname)
}

func (b *BlogDataAbstractor) inferBlogTitle(filename string) string {
	sepBySpecChars := splitAtSpecialChars(filename)
	parts := []string{}
	for _, s := range sepBySpecChars {
		parts = append(parts, splitCamelCaseAndNumbers(s)...)
	}

	spaceSeparated := strings.Join(parts, " ")
	return strings.Title(spaceSeparated)
}

func splitCamelCaseAndNumbers(whole string) []string {
	found := []string{}
	parts := RX.FindAllString(whole, -1)
	for _, p := range parts {
		found = append(found, findUpperCaseSequence(p)...)
	}
	return found
}

func findUpperCaseSequence(chars string) []string {
	if RX2REL.MatchString(chars) {
		subParts := RX2.FindAllStringSubmatch(chars, -1)
		return subParts[0][1:]
	}
	return []string{chars}
}

func splitAtSpecialChars(whole string) []string {
	rx := regexp.MustCompile("[^-_ ,.]*")
	return rx.FindAllString(whole, -1)
}

func (b *BlogDataAbstractor) findMdFileInAddDir() string {
	mds := fs.ReadDirEntriesEndingWith(b.addDir, "md", "txt")
	for _, md := range mds {
		return path.Join(b.addDir, md)
	}
	return ""
}

func (b *BlogDataAbstractor) inferBlogTitlePlain(filename string) string {
	rx := regexp.MustCompile("(^[a-z]+)|([A-Z][a-z]*)|([0-9]+)")
	parts := rx.FindAllString(filename, -1)
	dashSeparated := strings.Join(parts, "-")
	return strings.ToLower(dashSeparated)
}
