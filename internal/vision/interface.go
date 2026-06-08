package vision

type VisionService interface {
	ExtractTextFromImage(imageBytes []byte) (string, error)
}
