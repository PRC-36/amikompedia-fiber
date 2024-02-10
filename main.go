package main

import (
	"github.com/PRC-36/amikompedia-fiber/app"
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

	tokenMaker, err := token.NewJWTMaker(viperConfig.TokenSymetricKey, viperConfig.TokenAccessDuration)
	if err != nil {
		log.Fatalf("Failed to create JWT Maker: %v", err)
	}

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
	})

	err = fiber.Listen(":" + viperConfig.PortApp)
	if err != nil {
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}

}
