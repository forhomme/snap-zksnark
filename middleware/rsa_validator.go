package middleware

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"smart-contract-service/configuration"
	"smart-contract-service/models"
	"strings"
	"time"
)

func RSASignatureValidator(cfg configuration.ConfigApp) echo.MiddlewareFunc {
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
			digest := sha256.Sum256([]byte(stringToSign))
			decodedSignature, _ := base64.StdEncoding.DecodeString(signature)
			pubKey, err := generatePublicKey(cfg)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.Response{
					Code:    http.StatusOK,
					Message: err.Error(),
				})
			}
			err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, digest[:], decodedSignature)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.Response{
					Code:    http.StatusOK,
					Message: err.Error(),
				})
			}

			c.Set("request-header", request)
			return next(c)
		}
	}
}

func generatePublicKey(cfg configuration.ConfigApp) (*rsa.PublicKey, error) {
	pub, err := os.ReadFile(cfg.PublicKeyLocation)
	if err != nil {
		return nil, err
	}
	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		return nil, err
	}
	if pubPem.Type != "PUBLIC KEY" {
		return nil, errors.New("Not Public Key")
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(pubPem.Bytes); err != nil {
		return nil, err
	}

	var pubKey *rsa.PublicKey
	var ok bool
	if pubKey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, err
	}
	return pubKey, nil
}
