package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"smart-contract-service/configuration"
	"smart-contract-service/models"
	"strings"
	"time"
)

func SignatureHMACValidator(cfg configuration.ConfigApp) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request := new(models.RequestHeader)
			signature := c.Request().Header.Get("X-SIGNATURE")
			partnerId := c.Request().Header.Get("X-PARTNER-ID")
			externalId := c.Request().Header.Get("X-EXTERNAL-ID")
			channelId := c.Request().Header.Get("CHANNEL-ID")
			deviceId := c.Request().Header.Get("X-DEVICE-ID")
			timestamp := c.Request().Header.Get("X-TIMESTAMP")
			request = &models.RequestHeader{
				Signature:  signature,
				PartnerId:  partnerId,
				ChannelId:  channelId,
				DeviceId:   deviceId,
				Timestamp:  timestamp,
				ExternalId: externalId,
			}

			timeParse, _ := time.Parse("2006-01-02T15:04:05.999TZ7", timestamp)
			if time.Now().After(timeParse.Add(time.Hour * 5)) {
				return c.JSON(http.StatusUnauthorized, models.Response{
					Code:    http.StatusOK,
					Message: fmt.Sprintf("Extend datetime"),
				})
			}
			method := c.Request().Method
			endpointUrl := c.Path()
			req, _ := io.ReadAll(c.Request().Body)
			c.Request().Body.Close()
			c.Request().Body = io.NopCloser(bytes.NewBuffer(req))
			var dataRequest = make([]byte, 0)
			if len(req) > 0 {
				var buf bytes.Buffer
				err := json.Compact(&buf, req)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, models.Response{
						Code:    http.StatusOK,
						Message: err.Error(),
					})
				}
				dataRequest = buf.Bytes()
			} else {
				dataRequest = []byte("")
			}
			stringToSign := fmt.Sprintf("%s:%s:%s:%s", method, endpointUrl, strings.ToLower(signHMAC256(dataRequest, []byte(cfg.Secret))), timestamp)
			if !verifyHMAC([]byte(stringToSign), []byte(cfg.Secret), signature) {
				return c.JSON(http.StatusUnauthorized, models.Response{
					Code:    http.StatusOK,
					Message: fmt.Sprintf("Signature not valid"),
				})
			}
			c.Set("request-header", request)
			return next(c)
		}
	}
}

func signHMAC256(msg, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return hex.EncodeToString(mac.Sum(nil))
}

func verifyHMAC(msg, key []byte, hash string) bool {
	sig, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false
	}

	mac := hmac.New(sha512.New, key)
	mac.Write(msg)

	return hmac.Equal(sig, mac.Sum(nil))
}
