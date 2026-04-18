package domain

type ThumbnailProvider interface {
	GenerateThumbnail(inputPath, outputPath string, width, height int) error
}
