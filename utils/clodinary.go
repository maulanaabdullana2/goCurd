package utils

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func InitCloudinary(cloudName, apiKey, apiSecret string) error {
	var err error
	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	return err
}

func UploadImageToCloudinary(file multipart.File) (string, error) {
	uploadResult, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: "your-folder", // Optional: specify a folder in Cloudinary
	})

	fmt.Printf("Upload Result: %+v\n", uploadResult)
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
