package web

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httptest"
	"net/url"
	"smart-contract-service/app/repo"
	"smart-contract-service/app/usecase"
	"smart-contract-service/configuration"
	"smart-contract-service/middleware"
	"strings"
	"testing"
	"time"
)

const (
	RSA_TOKEN_LOGIN   = "eQMzb57HHJyqkYanGcw5aQmwVcxd6IlMmUabZ6JeGo4d+QrYf/8qOIxW4xvzRgI3kZoZheQxId0Qu0KLc2LPqNfODNr6MfQC7aJluXHd7tNHYyvyYAcj6HjbmoU+TAc5wvz4NKgluaBYcFaHBR2vqPp91ZCQcncA2elLNwE8SGxfg7DPZfc3Jv7gPMkJfuctu6lMM7SrF9QlsXY/8m6wo2madLH3S5ZWHIB5iMd+r/u6tCJutbge4onBT2Q+hcEMystVespzaUJBiLW9/mfEuKEUZwzONBg7+mMsEUbVeJcng5a3Xn51awRsCGmVcoLCnPqzAl09iFlIbziZpe9+25RUbMWkCX9JppFNtCtffsf2pU+siFwaXmVuxW3E42JLLkehG5Upz9bf3XZSLCOf057iWsKSTfjQxTdnUZeuWjFiyOTPSWx6nfxm8WmS8fVEDP4vlvGxfZQtvsJ5SbfjKtHOJ/U3CL7pB4fKyfT8/HtnuskG7s3TeyrKV5WHlwIuKNJLtSt8Ppcx0N4Diew4VC7fGZcPi5B/ko7Vcw01q2NWiqEL5QxgZj8gYtKZlqaHucvS/Riu+gL1C4ArjnA11AkMHtmEy+sKLVsQvWM98AS9mez+UIAUCiupMdS/CvwWxOEqw79yW54xvkd+pYLlXm43CucnOfzWMcQOF85jPsA="
	RSA_TOKEN_VERIFY  = "ZTa+vPBCFzGFjcSC/3ww1rEsr75ZnSWdwedz9wcu3AKp7CnqKCmVwUaonk1M0fPzpyUT+gjfniCix2pbHmzj2AlwBouW4C5/I+r6MxBPPjIDeH7cDtdvNG67n8yvXUds/NFgqmN84uiiyeL3kbkrPehuvntZVXn/DCH3ITynbYTySJyaC/MaehI3onFQx/6LWQDTfpiBFxc/ad0rTPeM3EXjucZ2XMsJAokU3tHVa50q2UheGP5bzo6JzRVPSGTHUt/squUvQollzZTAUVtUZQ6btJ3hbW4WzB++bl3aRouGyIzqbDfDORrLFGpUOkXI3QFFfz7RMn1lClAhFfytQd9h7c0TKEfUtnJu6vD39cAKwkRAL7+qmdow0e5jdDMk1R07qubfOJhwPQoej0gClFakubm+TxFlgCrRyvz8SGlV+1JnNXPRtrnw6QfEXtwtfpLMUq+WAgxnOve36qU6fYsYvd4sMGfTwZD5R+0NQPiHDKrB9vnmXi7Rv2g5hSEdbNm24bJW4SUsmjVbydfIKqfeHGqgcoELZw99K6gcAT7N6+qgNGtoJ0MLR0qIgC2YlOyNFDxp1pMQx6FZDFoX9NGypowNuDT3dLK1ED6BF2ecBkHr0nAFTuE5n0MJZkKCfIQcmhxx14n0Us4UZmwOBbbU8STncm9LsnbntsbmrO8="
	RSA_TOKEN_PING    = "cyS5FULcQMlZaPI1quh0dNPkL7PbGNsTB4itkAUxzjxfzfas+HivevvmIbqYn/VqEJLkmdGGOWvdJ8ObBk1tg9lvWYTLGE2KOJhlCmm+oZBDHLfs7Xm1AgcuVtYB43x0CPVl6Ju3G9/hHwQtRIPGTvFkdLTPhzHg9yQzb3sVdQJxUqKrtentoYXW65XSz0bWmhxWqAKPF8ElAw+WeXyeSYSqC9fyfI0QvobTyJUuA02f98h2afGRtGFMXkdmHTvJDiNk/enIU+bvxA+ZEWcXTPaFUdPUT/dCWbTXe0y2UxgLJ2KU+ZboATyQtOfy9bJsRQIpLL/uTvW0E3p70AnpNPvInhWD7ORUC9XUIMOVwGfUc2Y4JYkhQWHNbK2Q6McFT6oEhd8y2Yj8+cW5k90OrRZr9teGX/GZ+3TktFBKZG3Sxb17GUt800v1J5o1djuuWSjrFSI7R2EM73bmPAc/PxumNITMXdLYZNDwPZy0VfjdgK38TBWQVMfFOQDY7Qnvf5e0OKjAjMB3XFj5E3/0jwUv3CU1jicamfw36WOVQ13d2tBILGbGXwUYRqrovQnCGDyIJKur/dyjQzTDBRki6E60HIjySb90QyiWUrJVMhSBhnToDLXRQ4aN9A2UEqiOcSzfQyIiZ13c0XhGXNPK2qZCd9PyMfHLc8V4KZpmvVg="
	HMAC_TOKEN_LOGIN  = "Gn1ifOFJscla8mnSsiX45BysvVRzYlpOM79BVJBcayAkloygdF4t9/vzzXPovX1j7jWoJRTQgtmld+zysyl93g=="
	HMAC_TOKEN_PING   = "6wZMmvy9dme0IziG3XH32eTWgsDQjgX5v9n3Op9OEcaQlVWlQZUANRCdhwBaGxLDCiJw0cgxjXixUYy3qX7Mfg=="
	HMAC_TOKEN_VERIFY = "Q7IR3OuwuWACTOW7sGMnfsAWsg7Xu+7QYj5Ry17MrTrW81CvXi3RaIcqFtVyiSuC/8wDggPdp5Z+C+hGU8phHQ=="
	ACCESS_TOKEN      = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjY2OGEwMDNkLTM4ZDItNDQ3OS1hODUyLTM1MTEzZGEwMzE5OSIsInVzZXJuYW1lIjoiZGV2ZWxvcGVyIiwibG9naW5BdCI6IjIwMjMtMDYtMTRUMDc6MjA6NTcrMDc6MDAiLCJleHBpcmVBdCI6IjIwMjMtMDYtMTRUMTI6MjA6NTcrMDc6MDAiLCJleHAiOjE2ODY3MjAwNTd9.BaqSKe9yRP19Mai61bFFH1xCUDehmmMbj9pkjeg_wPc"
	PARTNER_ID        = "668a003d-38d2-4479-a852-35113da03199"
	CHANNEL_ID        = "111111"
	DEVICE_ID         = "android"
	PROOF_ELLIPTIC    = "AAAAAQAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEl8fGU0Y2IzYjIwLTRlODEtNDY1OS04ZTlkLWNkMDVhM2E5YTMyNA=="
	PROOF_HASH        = "AAAAAQAAAAAAAAABJnnM/EuwL9Og8Y+E6uVyw67lw97tTvqeghW+eco6msF8fGU0Y2IzYjIwLTRlODEtNDY1OS04ZTlkLWNkMDVhM2E5YTMyNA=="
	PROOF_EDDSA       = "AAAABgAAAAAAAAAGAsXLbg4fcxt9MwszGmmipeQqAhs12w8U/Gn3u5xFVH0Tt/j3efLGm4qWraKoKFUrRgQ6b8nR6yBA7laCazX+tCOoxO36EYX43Eg7q85p8BKpmIrCCHoz35utpROxBsOrFu/hgJSy/R6XE3TKclFiWj9xPOqpmLhg89gM5M1yxe8Et1j1sKY9aqu2sMbwR2G7hQWownwyikKp0rZqwru7ZiJlOTNlMjk5NGE5ZGVhYmE0NTNiOGRlMDVkLTM1Mi04fHxlNGNiM2IyMC00ZTgxLTQ2NTktOGU5ZC1jZDA1YTNhOWEzMjQ="
)

func BenchmarkHTTP_RSAPingHandler(b *testing.B) {
	b.Run("Endpoint: GET /rsa/ping", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/rsa/ping", nil)
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", RSA_TOKEN_PING)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.RSASignatureValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.PingHandler(c)
			}
		})
	})
}

func BenchmarkHTTP_HMACPingHandler(b *testing.B) {
	b.Run("Endpoint: GET /hmac/ping", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/rsa/ping", nil)
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", HMAC_TOKEN_PING)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.SignatureHMACValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.PingHandler(c)
			}
		})
	})
}

func BenchmarkHTTP_RSAVerifyEllipticProof(b *testing.B) {
	b.Run("Endpoint: POST /rsa/proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				f := make(url.Values)
				f.Set("proof", PROOF_ELLIPTIC)
				f.Set("algo", "elliptic")
				req := httptest.NewRequest(http.MethodGet, "/rsa/proof", strings.NewReader(f.Encode()))
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", RSA_TOKEN_VERIFY)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.RSASignatureValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.VerifyProof(c)
			}
		})
	})
}

func BenchmarkHTTP_HMACVerifyEllipticProof(b *testing.B) {
	b.Run("Endpoint: POST /hmac/proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				f := make(url.Values)
				f.Set("proof", PROOF_ELLIPTIC)
				f.Set("algo", "elliptic")
				req := httptest.NewRequest(http.MethodGet, "/rsa/proof", strings.NewReader(f.Encode()))
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", HMAC_TOKEN_VERIFY)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.SignatureHMACValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.VerifyProof(c)
			}
		})
	})
}

func BenchmarkHTTP_RSAVerifyHashProof(b *testing.B) {
	b.Run("Endpoint: POST /rsa/proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				f := make(url.Values)
				f.Set("proof", PROOF_HASH)
				f.Set("algo", "hash")
				req := httptest.NewRequest(http.MethodGet, "/rsa/proof", strings.NewReader(f.Encode()))
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", RSA_TOKEN_VERIFY)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.RSASignatureValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.VerifyProof(c)
			}
		})
	})
}

func BenchmarkHTTP_HMACVerifyHashProof(b *testing.B) {
	b.Run("Endpoint: POST /hmac/proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				f := make(url.Values)
				f.Set("proof", PROOF_HASH)
				f.Set("algo", "hash")
				req := httptest.NewRequest(http.MethodGet, "/rsa/proof", strings.NewReader(f.Encode()))
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", HMAC_TOKEN_VERIFY)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.SignatureHMACValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.VerifyProof(c)
			}
		})
	})
}

func BenchmarkHTTP_RSAVerifyEddsaProof(b *testing.B) {
	b.Run("Endpoint: POST /rsa/proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				f := make(url.Values)
				f.Set("proof", PROOF_EDDSA)
				f.Set("algo", "eddsa")
				req := httptest.NewRequest(http.MethodGet, "/rsa/proof", strings.NewReader(f.Encode()))
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", RSA_TOKEN_VERIFY)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.RSASignatureValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.VerifyProof(c)
			}
		})
	})
}

func BenchmarkHTTP_HMACVerifyEddsaProof(b *testing.B) {
	b.Run("Endpoint: POST /hmac/proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				f := make(url.Values)
				f.Set("proof", PROOF_EDDSA)
				f.Set("algo", "eddsa")
				req := httptest.NewRequest(http.MethodGet, "/rsa/proof", strings.NewReader(f.Encode()))
				req.Header.Add("Content-Type", "application/json;charset=UTF-8")
				req.Header.Add("Authorization", "Bearer "+ACCESS_TOKEN)
				req.Header.Add("X-SIGNATURE", HMAC_TOKEN_VERIFY)
				req.Header.Add("X-PARTNER-ID", PARTNER_ID)
				req.Header.Add("CHANNEL-ID", CHANNEL_ID)
				req.Header.Add("X-DEVICE-ID", DEVICE_ID)
				req.Header.Add("X-EXTERNAL-ID", uuid.New().String())
				req.Header.Add("X-TIMESTAMP", time.Now().Format("2006-01-02T15:04:05.999TZ7"))
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				e.Use(middleware.SignatureHMACValidator(configMain.Config))
				e.Use(middleware.AccessTokenValidator(configMain.Config))
				handler.VerifyProof(c)
			}
		})
	})
}

func BenchmarkHTTP_GetEllipticProof(b *testing.B) {
	b.Run("Endpoint: GET /proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				q := make(url.Values)
				q.Set("id", "e4cb3b20-4e81-4659-8e9d-cd05a3a9a324")
				q.Set("algo", "elliptic")
				req := httptest.NewRequest(http.MethodGet, "/proof?"+q.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				handler.GetProof(c)
			}
		})
	})
}

func BenchmarkHTTP_GetHashProof(b *testing.B) {
	b.Run("Endpoint: GET /proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				q := make(url.Values)
				q.Set("id", "e4cb3b20-4e81-4659-8e9d-cd05a3a9a324")
				q.Set("algo", "hash")
				req := httptest.NewRequest(http.MethodGet, "/proof?"+q.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				handler.GetProof(c)
			}
		})
	})
}

func BenchmarkHTTP_GetEddsaProof(b *testing.B) {
	b.Run("Endpoint: GET /proof", func(b *testing.B) {
		configMain := configuration.ServiceApp{
			EnvVariable: "DEV",
			Path:        "/home/ramadhoni/GolandProjects/zkSnark",
		}
		configMain.Load()
		handler := newHandler(configMain)

		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				// Setup
				e := echo.New()
				q := make(url.Values)
				q.Set("id", "e4cb3b20-4e81-4659-8e9d-cd05a3a9a324")
				q.Set("algo", "eddsa")
				req := httptest.NewRequest(http.MethodGet, "/proof?"+q.Encode(), nil)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)

				handler.GetProof(c)
			}
		})
	})
}

func newHandler(configMain configuration.ServiceApp) *HTTP {

	dbConn := configuration.InitSingleDB(configMain.Config.PostgreConnection, configMain.Config.LogMode)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     configMain.Config.RedisConnection,
		Password: "",
		DB:       0,
	})
	repoDb := repo.NewDatabaseConnection(dbConn)
	repoRedis := repo.NewRedisConnection(redisClient)

	uc := usecase.NewUsecase(repoRedis, repoDb, configMain.Config)

	handler := NewHTTP(configMain.Config, uc)
	return handler
}
