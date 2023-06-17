package middleware

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"smart-contract-service/configuration"
	"smart-contract-service/models"
	"strings"
)

func AccessTokenValidator(cfg configuration.ConfigApp) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := new(models.RequestHeader)
			authorization := c.Request().Header.Get("Authorization")
			if len(authorization) > 0 {
				request.Authorization = authorization
			}
			if session, err := validateJWTtoken(request.Authorization, cfg); err != nil {
				return c.JSON(http.StatusUnauthorized, models.Response{
					Code:    http.StatusUnauthorized,
					Message: err.Error(),
				})
			} else {
				c.Set("session", session)
				c.Set("request", request)
			}
			return next(c)
		}
	}
}

func validateJWTtoken(auth string, cfg configuration.ConfigApp) (data models.JwtCustomClaims, err error) {
	if len(auth) == 0 {
		return
	}

	authorization := strings.Split(auth, "Bearer ")
	if len(authorization) < 0 {
		return
	}
	tokenString := authorization[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data.ID = claims["id"].(string)
		data.Username = claims["username"].(string)
		data.LoginAt = claims["loginAt"].(string)
		data.ExpireAt = claims["expireAt"].(string)
	} else {
		return
	}
	return
}
