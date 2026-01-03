package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"

	appConfig "github.com/ulumfr/ulumfr-be/pkg/config"
)

// R2Client handles Cloudflare R2 operations
type R2Client struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
	publicURL string
}

// PresignedURLRequest contains the request for generating a presigned URL
type PresignedURLRequest struct {
	FileName    string `json:"file_name" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
	Folder      string `json:"folder"` // e.g., "projects", "resumes", "careers"
}

// PresignedURLResponse contains the response with presigned URL
type PresignedURLResponse struct {
	UploadURL string `json:"upload_url"`
	FileURL   string `json:"file_url"`
	Key       string `json:"key"`
	ExpiresIn int    `json:"expires_in"` // seconds
}

// NewR2Client creates a new Cloudflare R2 client
func NewR2Client(cfg *appConfig.Config) (*R2Client, error) {
	if cfg.R2AccountID == "" || cfg.R2AccessKeyID == "" || cfg.R2SecretAccessKey == "" {
		return nil, fmt.Errorf("R2 configuration is incomplete")
	}

	// Cloudflare R2 endpoint
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID)

	// Create custom resolver
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: r2Endpoint,
		}, nil
	})

	// Load AWS config with custom credentials
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.R2AccessKeyID,
			cfg.R2SecretAccessKey,
			"",
		)),
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg)
	presigner := s3.NewPresignClient(client)

	log.Info().
		Str("bucket", cfg.R2BucketName).
		Msg("R2 client initialized")

	return &R2Client{
		client:    client,
		presigner: presigner,
		bucket:    cfg.R2BucketName,
		publicURL: cfg.R2PublicURL,
	}, nil
}

// GeneratePresignedPutURL generates a presigned URL for uploading a file
func (r *R2Client) GeneratePresignedPutURL(ctx context.Context, req PresignedURLRequest) (*PresignedURLResponse, error) {
	// Generate unique key
	timestamp := time.Now().Unix()
	key := fmt.Sprintf("%s/%d-%s", req.Folder, timestamp, req.FileName)

	// Set expiration (15 minutes)
	expiration := 15 * time.Minute

	// Generate presigned URL
	presignedReq, err := r.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(req.ContentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Build the public URL
	fileURL := fmt.Sprintf("%s/%s", r.publicURL, key)

	return &PresignedURLResponse{
		UploadURL: presignedReq.URL,
		FileURL:   fileURL,
		Key:       key,
		ExpiresIn: int(expiration.Seconds()),
	}, nil
}

// DeleteObject deletes an object from R2
func (r *R2Client) DeleteObject(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// IsConfigured checks if R2 client is properly configured
func (r *R2Client) IsConfigured() bool {
	return r != nil && r.client != nil
}
