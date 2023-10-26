package storage

import (
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type AzureBlobStorage interface {
	GeneratePresignedUrlToUpload(imageUrl string) (string, error)
	GeneratePresignedUrlToRead(imageUrl string) (string, error)
	GetFrontDoorUrl(imageName string) string
}

type AzureBlobStorageConfig struct {
	AccountName   string `yaml:"accountName" env:"AZURE_BLOB_ACCOUNT_NAME" env-default:""`
	FrontDoorUrl  string `yaml:"frontDoorUrl" env:"AZURE_BLOB_FRONT_DOOR_URL" env-default:""`
	AccountKey    string `yaml:"accountKey" env:"AZURE_BLOB_ACCOUNT_KEY" env-default:""`
	ContainerName string `yaml:"containerName" env:"AZURE_BLOB_CONTAINER_NAME" env-default:""`
}

type azureBlobStorage struct {
	AzureBlobStorageConfig
}

func NewAzureBlobStorage(cfg AzureBlobStorageConfig) AzureBlobStorage {
	return &azureBlobStorage{
		cfg,
	}
}

func (az *azureBlobStorage) GeneratePresignedUrlToUpload(imageName string) (string, error) {
	return az.generatePresignedUrl(imageName, azblob.BlobSASPermissions{Write: true, Permissions: true})
}

func (az *azureBlobStorage) GeneratePresignedUrlToRead(imageName string) (string, error) {
	return az.generatePresignedUrl(imageName, azblob.BlobSASPermissions{Read: true, Permissions: true})
}

func (az *azureBlobStorage) GetFrontDoorUrl(imageName string) string {
	return fmt.Sprintf("http://%s/%s/%s", az.FrontDoorUrl, az.ContainerName, imageName)
}

func (az *azureBlobStorage) generatePresignedUrl(imageName string, accessPolicy azblob.BlobSASPermissions) (string, error) {
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

	containerURL := serviceURL.NewContainerURL(az.ContainerName)

	blobURL := containerURL.NewBlobURL(imageName)

	start := time.Now()
	expiry := start.Add(24 * time.Hour)

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     start.UTC(),
		ExpiryTime:    expiry.UTC(),
		ContainerName: az.ContainerName,
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
