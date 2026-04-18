package infrastructure

import (
	"github.com/disintegration/imaging"
)

type ImagingProvider struct{}

func NewImagingProvider() *ImagingProvider {
	return &ImagingProvider{}
}

func (p *ImagingProvider) GenerateThumbnail(inputPath, outputPath string, width, height int) error {
	src, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	dst := imaging.Thumbnail(src, width, height, imaging.Lanczos)

	return imaging.Save(dst, outputPath)
}
