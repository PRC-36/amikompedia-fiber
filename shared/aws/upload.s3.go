package aws

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
)

const (
	ImgHeader = "HEADER"
	ImgAvatar = "AVATAR"
	ImgPost   = "POST"
)

type AwsS3 struct {
	S3        *s3.S3
	AwsBucket string
}

func NewAwsS3(s3 *s3.S3, awsBucket string) AwsS3Action {
	return &AwsS3{S3: s3, AwsBucket: awsBucket}
}

type AwsS3Action interface {
	UploadFile(file *multipart.FileHeader, imgType string) (Image, error)
	//ValidateImageType(contentType string) bool
}

func (a *AwsS3) UploadFile(file *multipart.FileHeader, imgType string) (Image, error) {
	src, err := file.Open()
	if err != nil {
		return Image{}, err
	}
	defer src.Close()

	// Read the file content into a byte slice
	fileBytes := make([]byte, file.Size)
	_, err = src.Read(fileBytes)
	if err != nil {
		return Image{}, err
	}

	// Resize the image
	resizedBytes, objectKey, err := a.resizeImageBytes(fileBytes, file.Filename, imgType) // Adjust dimensions as needed
	if err != nil {
		return Image{}, err
	}

	// Specify the destination path in the S3 bucket (e.g., "uploads/")

	//objectKey := "uploads/" + file.Filename

	// Get file information for content type and size
	fileType := http.DetectContentType(resizedBytes)
	size := int64(len(resizedBytes))

	if !a.validateImageType(fileType) {
		return Image{}, fmt.Errorf("invalid image type")
	}

	// Create an io.Reader from the resizedBytes
	fileReader := bytes.NewReader(resizedBytes)

	// Upload the resized file to S3
	//_, err = s3.New(a.Session).PutObject(&s3.PutObjectInput{
	//	Bucket:        aws.String(a.AwsBucket),
	//	Key:           aws.String(objectKey),
	//	ACL:           aws.String("public-read"),
	//	Body:          fileReader,
	//	ContentLength: aws.Int64(size),
	//	ContentType:   aws.String(fileType),
	//})

	_, err = a.S3.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(a.AwsBucket),
		Key:           aws.String(objectKey),
		ACL:           aws.String("public-read"),
		Body:          fileReader,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
	})

	if err != nil {
		return Image{}, err
	}

	//log.Println("object:", object)
	//log.Println("object url:", object.GoString())

	newImage := Image{
		FilePath:  objectKey,
		ImageType: imgType,
		ImageUrl:  fmt.Sprintf("https://%s.s3.amazonaws.com/%s", a.AwsBucket, objectKey),
	}

	return newImage, nil
}

func (a *AwsS3) validateImageType(contentType string) bool {
	// Add more allowed content types if needed
	allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
	for _, t := range allowedTypes {
		if t == contentType {
			return true
		}
	}
	return false
}

// ResizeImageBytes resizes the image represented as bytes to the specified width and height
func (a *AwsS3) resizeImageBytes(imageBytes []byte, fileName, imgType string) ([]byte, string, error) {
	var resizedImg image.Image
	var objectKey string
	// Check Image Type and Resize Image

	// Decode the original image
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, "", err
	}

	// Resize the image (adjust dimensions as needed) as ImgType
	switch imgType {
	case ImgHeader:
		resizedImg = resize.Resize(1200, 400, img, resize.Lanczos2)
		objectKey = "profile-header/" + fileName
	case ImgAvatar:
		resizedImg = resize.Resize(500, 500, img, resize.Lanczos3)
		objectKey = "profile-avatar/" + fileName
	case ImgPost:
		resizedImg = img
		objectKey = "post/" + fileName
	default:
		resizedImg = img
		objectKey = fileName
	}
	//resizedImg := resize.Resize(width, height, img, resize.Lanczos3)

	// Encode the resized image
	var resizedBuffer bytes.Buffer
	ext := a.detectImageFormat(imageBytes)
	switch ext {
	case "png":
		err = png.Encode(&resizedBuffer, resizedImg)
	case "jpeg":
		err = jpeg.Encode(&resizedBuffer, resizedImg, nil)
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, "", err
	}

	return resizedBuffer.Bytes(), objectKey, nil
}

// detectImageFormat detects the image format based on the provided image bytes
func (a *AwsS3) detectImageFormat(imageBytes []byte) string {
	contentType := http.DetectContentType(imageBytes)
	switch contentType {
	case "image/jpeg":
		return "jpeg"
	case "image/png":
		return "png"
	// Add additional cases for other supported formats
	default:
		return ""
	}
}
