package main

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

type BackupServiceI interface {
	uploadToR2(ctx context.Context, name, path string) error
	dumpToFile(filepath string) (string, error)
	deleteFile(compressedFilePath string) error
}

type BackupService struct {
	config Config
}

func NewBackUpService(config Config) BackupServiceI {
	return &BackupService{
		config: config,
	}
}

func (b *BackupService) uploadToR2(ctx context.Context, name, path string) error {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				b.config.CLOUD_FLARE_ACCESS_KEY_ID,
				b.config.CLOUD_FLARE_SECRET_ACCESS_KEY,
				"",
			),
		),
		config.WithRegion(b.config.CLOUD_FLARE_R2_REGION),
		config.WithRetryMaxAttempts(5),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		// o.BaseEndpoint = aws.String(b.config.CLOUD_FLARE_R2_ENDPOINT)
	})

	if b.config.BUCKET_SUBFOLDER != "" {
		name = b.config.BUCKET_SUBFOLDER + "/" + name
	}

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Couldn't open file %v to upload. Here's why: %v\n", name, err)
	} else {
		defer file.Close()
		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(b.config.CLOUD_FLARE_R2_BUCKET),
			Key:    aws.String(name),
			Body:   file,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
				log.Printf("Error while uploading object to %s. The object is too large.\n"+
					"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
					"or the multipart upload API (5TB max).", b.config.CLOUD_FLARE_R2_BUCKET)
			} else {
				log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					path, b.config.CLOUD_FLARE_R2_BUCKET, name, err)
			}
		}
		// } else {
		// 	err = s3.NewObjectExistsWaiter(client).Wait(
		// 		ctx, &s3.HeadObjectInput{Bucket: aws.String(b.config.CLOUD_FLARE_R2_BUCKET), Key:
		// aws.String(name)}, time.Minute)
		// 	if err != nil {
		// 		log.Printf("Failed attempt to wait for object %s to exist.\n", name)
		// 	}
		// 	log.Println("Upload successful")
		// }
	}

	return err
}

func (b *BackupService) dumpToFile(filepath string) (string, error) {
	log.Println("Dumping to file...")
	compressedFilePath := fmt.Sprintf("%s.gz", filepath)
	cmd := exec.Command(
		"sh", "-c",
		fmt.Sprintf(
			"MYSQL_PWD=%s mysqldump --single-transaction --protocol=TCP --add-drop-table --quick --lock-tables=false --user=%s --host=%s --port=%s %s | gzip > %s",
			b.config.BACKUP_DATABASE_PASSWORD,
			b.config.BACKUP_DATABASE_USER,
			b.config.BACKUP_DATABASE_HOST,
			b.config.BACKUP_DATABASE_PORT,
			b.config.BACKUP_DATABASE_NAME,
			compressedFilePath,
		),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run mysqldump: %v\nOutput: %s", err, output)
	}

	// err = compressFile(filepath, compressedFilePath)
	// if err != nil {
	// 	return fmt.Errorf("failed to compress snapshot file: %v", err)
	// }

	return compressedFilePath, nil
}

func (b *BackupService) deleteFile(compressedFilePath string) error {
	log.Println("Deleting files...")
	if err := os.Remove(compressedFilePath); err != nil {
		return fmt.Errorf("failed to delete compressed snapshot file: %w", err)
	}

	return nil
}

func compressFile(source, target string) error {
	log.Println("Compressing file...")
	inFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	_, err = io.Copy(gzipWriter, inFile)
	return err
}
