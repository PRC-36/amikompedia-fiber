package middleware

import (
	"errors"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "Authorization"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeaderKey := ctx.Get(authorizationHeaderKey)

		fields := strings.Fields(authorizationHeaderKey)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			res, code := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   http.StatusUnauthorized,
					Status: err.Error(),
				},
			)
			return ctx.Status(code).JSON(res)
		}

		accessToken := fields[1]

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			res, code := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   http.StatusUnauthorized,
					Status: err.Error(),
				},
			)
			return ctx.Status(code).JSON(res)
		}

		ctx.Locals(AuthorizationPayloadKey, payload)
		return ctx.Next()
	}
}
