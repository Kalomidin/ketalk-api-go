package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type r2CloudFlare struct {
	client *s3.Client
	cfg    *R2CloudFlareConfig
}

type R2CloudFlareConfig struct {
	Bucket          string `yaml:"bucket" env:"R2_CLOUDFLARE_BUCKET" env-default:""`
	AccountId       string `yaml:"accountId" env:"R2_CLOUDFLARE_ACCOUNT_ID" env-default:""`
	AccessKeyId     string `yaml:"accessKeyId" env:"R2_CLOUDFLARE_ACCESS_KEY_ID" env-default:""`
	AccessKeySecret string `yaml:"accessKeySecret" env:"R2_CLOUDFLARE_ACCESS_KEY_SECRET" env-default:""`
	Env             string `yaml:"env" env:"ENV" env-default:"local"`
	PublicR2Url     string `yaml:"publicR2Url" env:"PUBLIC_R2_URL" env-default:""`
}

func NewR2CloudFlare(ctx context.Context, cfg *R2CloudFlareConfig) Storage {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountId),
		}, nil
	})
	r2Cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyId, cfg.AccessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed to load config, %v", err))
	}
	client := s3.NewFromConfig(r2Cfg)
	return &r2CloudFlare{
		client: client,
		cfg:    cfg,
	}
}

func (r *r2CloudFlare) generateKey(imageUrl, containerName string) string {
	return fmt.Sprintf("%s/%s/%s", r.cfg.Env, containerName, imageUrl)
}

func (r *r2CloudFlare) GeneratePresignedUrlToUpload(ctx context.Context, imageUrl, containerName string) (string, error) {
	presignedClient := s3.NewPresignClient(r.client)
	key := r.generateKey(imageUrl, containerName)
	presignResult, err := presignedClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: &r.cfg.Bucket,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}
	return presignResult.URL, nil
}

func (r *r2CloudFlare) GetURLToRead(imageUrl, containerName string) string {
	key := r.generateKey(imageUrl, containerName)
	return fmt.Sprintf("%s/%s", r.cfg.PublicR2Url, key)
}

func (r *r2CloudFlare) GetUserImage(image string) string {
	if strings.Contains(image, "http") {
		return image
	}
	return r.GetURLToRead(image, ContainerProfiles)
}
