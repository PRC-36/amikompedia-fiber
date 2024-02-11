package main

import (
	"github.com/PRC-36/amikompedia-fiber/app"
	"github.com/PRC-36/amikompedia-fiber/shared/aws"
	"github.com/PRC-36/amikompedia-fiber/shared/mail"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"log"
)

func main() {
	viperConfig, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tokenMaker := token.NewJWTMaker()
	if err != nil {
		log.Fatalf("Failed to create JWT Maker: %v", err)
	}

	s3Client, err := aws.NewSessionAWSS3(viperConfig)
	if err != nil {
		log.Fatalf("Failed to create AWS S3 session: %v", err)
	}

	awsS3 := aws.NewAwsS3(s3Client, viperConfig.AWSS3Bucket)
	mailSender := mail.NewGmailSender(viperConfig.EmailName, viperConfig.EmailSender, viperConfig.EmailPassword)
	db := app.NewDatabaseConnection(viperConfig.DBDsn)
	validate := app.NewValidator()
	fiber := app.NewFiber(viperConfig)

	app.Bootstrap(&app.BootstrapConfig{
		DB:          db,
		App:         fiber,
		Validate:    validate,
		TokenMaker:  tokenMaker,
		EmailSender: mailSender,
		AwsS3:       awsS3,
		ViperConfig: viperConfig,
	})

	err = fiber.Listen(":" + viperConfig.PortApp)
	if err != nil {
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}

}
