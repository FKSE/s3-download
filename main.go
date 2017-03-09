package main

import (
	"flag"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	accessKey = flag.String("accessKey", "", "S3 Access Key")
	secretKey = flag.String("secretKey", "", "S3 Secret Key")
	region    = flag.String("region", "eu-west-1", "S3-Region")
	output    = flag.String("output", "", "Output-Filename")
	urlStr    = flag.String("url", "", "Object-Url")
)

func main() {
	flag.Parse()
	// parse url
	url, err := url.Parse(*urlStr)
	if err != nil {
		log.Fatalln(err)
	}
	if url.Scheme != "s3" {
		log.Fatalln("Url has to start with s3://")
	}
	// init s3
	sess, err := session.NewSession(&aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentials(*accessKey, *secretKey, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}
	svc := s3.New(sess)
	// get object
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &url.Host,
		Key:    &url.Path,
	})
	if err != nil {
		log.Fatalln(err)
	}
	// open file for wrting
	file, err := os.Create(*output)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	// copy body
	n, err := io.Copy(file, out.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Downloaded %d bytes to %s", n, *output)
}
