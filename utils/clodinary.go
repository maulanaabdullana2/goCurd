package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/joho/godotenv"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Ambil kredensial dari variabel lingkungan
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return fmt.Errorf("one or more environment variables are not set")
	}

	cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("error creating Cloudinary client: %v", err)
	}

	return nil
}

// UploadImageToCloudinary uploads an image to Cloudinary and returns the URL.
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
