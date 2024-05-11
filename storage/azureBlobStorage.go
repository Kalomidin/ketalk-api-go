package storage

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

const (
	ContainerProfiles = "profiles"
	ContainerItems    = "items"
)

type AzureBlobStorage interface {
	GeneratePresignedUrlToUpload(imageUrl, containerName string) (string, error)
	// GeneratePresignedUrlToRead(imageUrl, containerName string) (string, error)
	GetFrontDoorUrl(imageUrl, containerName string) string
	GetUserImage(image string) string
}

type AzureBlobStorageConfig struct {
	AccountName  string `yaml:"accountName" env:"AZURE_BLOB_ACCOUNT_NAME" env-default:""`
	FrontDoorUrl string `yaml:"frontDoorUrl" env:"AZURE_BLOB_FRONT_DOOR_URL" env-default:""`
	BlobUrl      string `yaml:"blobUrl" env:"AZURE_BLOB_URL" env-default:""`
	AccountKey   string `yaml:"accountKey" env:"AZURE_BLOB_ACCOUNT_KEY" env-default:""`
}

type azureBlobStorage struct {
	AzureBlobStorageConfig
}

func NewAzureBlobStorage(cfg AzureBlobStorageConfig) AzureBlobStorage {
	return &azureBlobStorage{
		cfg,
	}
}

func (az *azureBlobStorage) GeneratePresignedUrlToUpload(imageUrl, containerName string) (string, error) {
	return az.generatePresignedUrl(imageUrl, containerName, azblob.BlobSASPermissions{Write: true, Permissions: true})
}

// func (az *azureBlobStorage) GeneratePresignedUrlToRead(imageUrl, containerName string) (string, error) {
// 	return az.generatePresignedUrl(imageUrl, containerName, azblob.BlobSASPermissions{Read: true, Permissions: true})
// }

func (az *azureBlobStorage) GetFrontDoorUrl(imageUrl, containerName string) string {
	// return fmt.Sprintf("https://%s/%s/%s", az.BlobUrl, containerName, imageUrl)
	// TODO: use front door url
	url, _ := az.generatePresignedUrl(imageUrl, containerName, azblob.BlobSASPermissions{Read: true, Permissions: true})
	return url
	// return fmt.Sprintf("http://%s/%s/%s", az.FrontDoorUrl, containerName, imageUrl)
}

func (az *azureBlobStorage) generatePresignedUrl(imageName, containerName string, accessPolicy azblob.BlobSASPermissions) (string, error) {
	credential, err := azblob.NewSharedKeyCredential(az.AccountName, az.AccountKey)
	if err != nil {
		return "", err
	}

	serviceURL := azblob.NewServiceURL(
		url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s.blob.core.windows.net", az.AccountName),
		},
		azblob.NewPipeline(credential, azblob.PipelineOptions{}),
	)

	containerURL := serviceURL.NewContainerURL(containerName)

	blobURL := containerURL.NewBlobURL(imageName)

	start := time.Now()
	expiry := start.Add(24 * time.Hour)

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     start.UTC(),
		ExpiryTime:    expiry.UTC(),
		ContainerName: containerName,
		BlobName:      imageName,
		Permissions:   accessPolicy.String(),
	}.NewSASQueryParameters(credential)

	if err != nil {
		return "", err
	}
	sasURL := blobURL.URL()
	sasURL.RawQuery = sasQueryParams.Encode()

	return sasURL.String(), nil
}

func (az *azureBlobStorage) GetUserImage(image string) string {
	if strings.Contains(image, "http") {
		return image
	}
	return az.GetFrontDoorUrl(image, ContainerProfiles)
}
