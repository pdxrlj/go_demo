package main

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/airbusgeo/osio"
	osios3 "github.com/airbusgeo/osio/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/klauspost/compress/gzip"
)

const (
	S3Service    = "s3"
	MinioService = "minio"
)

func main() {
	zip()
}

func zip() {
	serviceName, regionName, accessKey, secretKey, endpoint :=
		MinioService,
		"cn-northwest-1",
		"3C0JZGAMQNDW4V79AFNM",
		"dC3JVwrefYyeBvRDnIS8XketdopCQW80E8wpJs9K",
		"https://obs.cn-north-4.myhuaweicloud.com"
	s, err := RegisterReader(serviceName, regionName, accessKey, secretKey, endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}
	s3r, err := osios3.Handle(context.Background(), osios3.S3Client(s), osios3.S3RequestPayer())
	if err != nil {
		fmt.Println(err)
		return
	}
	osr, err := osio.NewAdapter(s3r)
	if err != nil {
		fmt.Println(err)
		return
	}

	//uri:= "obs://pie-engine-image-data/USGS/collection02/level-1/standard/oli-tirs/2023/138/039/LC09_L1TP_138039_20230101_20230315_02_T1.tar.gz"
	//uri := "s3://USGS/collection02/level-1/standard/oli-tirs/2023/138/039/LC09_L1TP_138039_20230101_20230315_02_T1.tar.gz"
	uri := "s3://pie-engine-image-data/USGS/collection02/level-1/standard/oli-tirs/2023/138/039/LC09_L1TP_138039_20230101_20230315_02_T1.tar.gz"
	obj, err := osr.Reader(uri)
	if err != nil {
		fmt.Println(err)
		return
	}
	gzReader, err := gzip.NewReader(obj)
	if err != nil {
		panic(err)
	}

	reader := tar.NewReader(gzReader)
	for {
		header, err := reader.Next()
		if err != nil {
			return
		}
		fmt.Println(header.Name)
		continue
	}

}

func RegisterReader(serviceName, regionName, accessKey, secretKey, endpoint string) (*s3.Client, error) {
	var s3cl *s3.Client
	switch strings.ToLower(serviceName) {
	case strings.ToLower(MinioService):
		staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               endpoint, // or where ever you ran minio
				SigningRegion:     regionName,
				HostnameImmutable: true,
			}, nil
		})
		cfg := aws.Config{
			Region:           regionName,
			Credentials:      credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
			EndpointResolver: staticResolver,
		}
		s3cl = s3.NewFromConfig(cfg)
	case strings.ToLower(S3Service):
		provider := credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
				SessionToken:    "",
			},
		}
		s3cl = s3.New(s3.Options{
			Region:      regionName,
			Credentials: &provider,
		})
	default:
		return nil, errors.New("s3ServiceName only support s3 and minio")
	}

	return s3cl, nil
}
