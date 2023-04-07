package main

import (
	"context"
	"fmt"
	"log"

	"github.com/airbusgeo/godal"
	"github.com/airbusgeo/osio"
	osioS3 "github.com/airbusgeo/osio/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	region    = "cn-north-4"
	endpoint  = "https://obs.cn-north-4.myhuaweicloud.com"
	accessKey = "3C0JZGAMQNDW4V79AFNM"
	secretKey = "dC3JVwrefYyeBvRDnIS8XketdopCQW80E8wpJs9K"
)

func initGDAL2() {
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               endpoint, // or where ever you ran minio
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})
	cfg := aws.Config{
		Region:           region,
		Credentials:      credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		EndpointResolver: staticResolver,
	}
	s3client := s3.NewFromConfig(cfg)

	s3r, err := osioS3.Handle(context.Background(), osioS3.S3Client(s3client), osioS3.S3RequestPayer())
	if err != nil {
		log.Fatalf("osio s3 err:%s\n", err.Error())
	}
	adapter, err := osio.NewAdapter(s3r)
	if err != nil {
		log.Fatalf("osio NewAdapter err:%s\n", err.Error())
	}

	godal.RegisterAll()
	if err = godal.RegisterVSIHandler("s3://", adapter); err != nil {
		log.Fatalf("godal RegisterVSIHandler err:%s\n", err.Error())
	}
}

func initGDAL() {
	resolve := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint,
			PartitionID:       "aws",
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})

	awsConfig := aws.Config{
		Region:                      region,
		Credentials:                 credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		EndpointResolverWithOptions: resolve,
	}
	s3cl := s3.NewFromConfig(awsConfig)

	s3r, err := osioS3.Handle(context.Background(), osioS3.S3Client(s3cl), osioS3.S3RequestPayer())
	if err != nil {
		panic(err)
	}
	adapter, err := osio.NewAdapter(s3r, osio.NumCachedBlocks(10000))
	if err != nil {
		panic(err)
	}

	godal.RegisterAll()
	//godal.RegisterInternalDrivers()
	if err = godal.RegisterVSIHandler("s3://", adapter); err != nil {
		panic(err)
	}
}

type BoundGeom struct {
	MinX float64 `json:"minx,omitempty"`
	MinY float64 `json:"miny,omitempty"`
	MaxX float64 `json:"maxx,omitempty"`
	MaxY float64 `json:"maxy,omitempty"`
}

func main() {
	initGDAL2()

	file := "s3:///pie-engine-uav/20221028/783459288919322624/Product/DOM.tif"
	dataset, err := godal.Open(file)
	if err != nil {
		panic(err)
	}
	structure := dataset.Structure()
	fmt.Printf("dataset size: %dx%dx%d\n", structure.SizeX, structure.SizeY, structure.NBands)

	ref, err := godal.NewSpatialRefFromEPSG(4326)
	if err != nil {
		panic(err)
	}
	bounds, err := dataset.Bounds(ref)
	if err != nil {
		panic(err)
	}

	bound := &BoundGeom{
		MinX: bounds[0],
		MinY: bounds[1],
		MaxX: bounds[2],
		MaxY: bounds[3],
	}
	fmt.Println(bound)
}
