package database

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/viper"
)

var DirectoryPath string

func NewDatabase() (*s3manager.Uploader,error){
	
	region := viper.GetString("s3.region")
	access_key := viper.GetString("s3.access_key")
	secret_key := viper.GetString("s3.secret_key")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(access_key, secret_key, ""),
	})

	uploader := s3manager.NewUploader(sess)
	return uploader,err
}

func UploadS3(bucket string ,uploader *s3manager.Uploader, folder string,filename string) error {
	fmt.Println(DirectoryPath +folder +"/"+ filename)
	file, err  := os.Open(DirectoryPath +folder +"/"+ filename)
	if err != nil {
		fmt.Println("file not present ")
		return  err
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key: aws.String(folder+"/"+filename),
		Body: file,
	})
	if err != nil {
		// Print the error and exit.
		fmt.Printf("Unable to upload %q to %q, %v", filename, bucket, err)
		return err
	}
	return nil
}