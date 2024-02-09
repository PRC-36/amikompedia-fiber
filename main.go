package main

import (
	"github.com/PRC-36/amikompedia-fiber/app"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/go-playground/validator/v10"
	"log"
)

func main() {
	viperConfig, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	tokenMaker, err := token.NewJWTMaker(viperConfig.TokenSymetricKey, viperConfig.AccessTokenDuration)
	if err != nil {
		log.Fatalf("Failed to create JWT Maker: %v", err)
	}

	db := app.NewDatabaseConnection(viperConfig.DBDsn)
	validate := validator.New()
	fiber := app.NewFiber(viperConfig)

	app.Bootstrap(&app.BootstrapConfig{
		DB:         db,
		App:        fiber,
		Validate:   validate,
		TokenMaker: tokenMaker,
	})

	err = fiber.Listen(":" + viperConfig.PortApp)
	if err != nil {
		if err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}

}
