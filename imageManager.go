package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/ingmardrewing/fs"
	"github.com/ingmardrewing/img"
	"github.com/ingmardrewing/staticUtil"
)

var doUpload = true

func DoUpload(val bool) {
	doUpload = val
}

type ImgManager interface {
	PrepareImages()
	UploadImages()
	GetImageUrls() []string
	AddImageSize(size int) string
	AddCropImageSize(size int) string
}

// Upload images to the given awsbucket using
// environmental data as required by the aws
// packages
func NewImageManager(deployedStaticAssetsLocation, awsbucket, sourceimagepath string) *ImageManager {
	im := new(ImageManager)
	im.deployedStaticAssetsLocation = deployedStaticAssetsLocation
	im.awsbucket = awsbucket
	im.sourceimagepath = sourceimagepath
	return im
}

type ImageManager struct {
	deployedStaticAssetsLocation string
	sourceimagepath              string
	uploadimgagepaths            []string
	awsimageurls                 []string
	imagesizes                   []int
	cropimagesizes               []int
	awsbucket                    string
}

func (i *ImageManager) PrepareImages() {
	imgdir := fs.GetPathWithoutFilename(i.sourceimagepath)

	imgcropscaler := img.NewImgScaler(i.sourceimagepath, imgdir)
	imgcropscaler.PrepareResizeTo(i.cropimagesizes...)
	paths := imgcropscaler.GetResizeAndCropPaths()

	imgscaler := img.NewImgScaler(i.sourceimagepath, imgdir)
	imgscaler.PrepareResizeTo(i.imagesizes...)
	paths = append(paths, imgscaler.GetResizeOnlyPaths()...)
	i.uploadimgagepaths = append(paths, i.sourceimagepath)

	imgcropscaler.ResizeAndCrop()
	imgscaler.Resize()
}

func (i *ImageManager) UploadImages() {
	if !doUpload {
		return
	}
	for _, filepath := range i.uploadimgagepaths {
		filename := fs.GetFilenameFromPath(filepath)
		fullPathFromDocRoot := i.getPathFromDocRoot(filename)
		i.awsimageurls = append(i.awsimageurls, fullPathFromDocRoot)
	}
}

func (i *ImageManager) getPathFromDocRoot(filename string) string {
	return i.deployedStaticAssetsLocation + staticUtil.GenerateDatePath() + filename
}

func (i *ImageManager) GetImageUrls() []string {
	if !doUpload {
		log.Println("constructed image paths (not acquired via aws):")
		log.Println(i.uploadimgagepaths)
		return i.uploadimgagepaths
	}
	log.Println("image paths (acquired via aws):")
	log.Println(i.awsimageurls)

	imgUrls := []string{}
	for _, imgUrl := range i.awsimageurls {
		imgUrls = append(
			imgUrls,
			"/static-assets/blog/"+strings.Replace(imgUrl, "%2F", "/", -1))
	}
	return imgUrls
}

func (i *ImageManager) AddCropImageSize(size int) string {
	i.cropimagesizes = append(i.cropimagesizes, size)
	return i.getFileNameFor(size)
}

func (i *ImageManager) AddImageSize(size int) string {
	i.imagesizes = append(i.imagesizes, size)
	return i.getFileNameFor(size)
}

func (i *ImageManager) getFileNameFor(w int) string {
	tag := "-w" + strconv.Itoa(w)

	sf := fs.GetFilenameFromPath(i.sourceimagepath)
	parts := strings.Split(sf, ".")
	n := strings.Join(parts[:len(parts)-1], "")
	return n + tag + "." + parts[len(parts)-1]
}
