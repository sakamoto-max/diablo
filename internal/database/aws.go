package database

// import (
// 	"context"
// 	"fmt"

// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// )

// func a() error {

// 	// config.NewEnvConfig()
// 	config, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("ap-sount-2"))
// 	if err != nil {
// 		return fmt.Errorf("failed to create aws config : %w", err)
// 	}

// 	client := dynamodb.NewFromConfig(config)
// 	client.CreateTable()
// }
