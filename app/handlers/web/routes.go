package web

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"smart-contract-service/configuration"
	middleware2 "smart-contract-service/middleware"
)

type Routes struct {
	config configuration.ConfigApp
}

func NewRoutes(config configuration.ConfigApp) *Routes {
	return &Routes{
		config: config,
	}
}

func (route *Routes) RegisterServices(r *echo.Echo, handler *HTTP) {
	openRoutes := r.Group(route.config.RootURL)
	apiRoutes := r.Group(route.config.RootURL, middleware2.RSASignatureValidator(route.config))
	hmacRoutes := r.Group(route.config.RootURL, middleware2.SignatureHMACValidator(route.config))
	accessTokenRoute := r.Group(route.config.RootURL, middleware2.AccessTokenValidator(route.config))
	route.setMiddleware(apiRoutes)
	route.setMiddleware(accessTokenRoute)

	// Routes Endpoint
	openRoutes.POST("/signup", handler.SignUp)
	openRoutes.POST("/token", handler.Token)
	openRoutes.POST("/token-hmac", handler.TokenHMAC)
	openRoutes.GET("/proof", handler.GetProof)
	apiRoutes.POST("/rsa/login", handler.Login)
	hmacRoutes.POST("/hmac/login", handler.Login)
	apiRoutes.POST("/refresh", handler.RefreshToken)
	hmacRoutes.POST("/hmac/refresh", handler.RefreshToken)
	accessTokenRoute.GET("/rsa/ping", handler.PingHandler, middleware2.RSASignatureValidator(route.config))
	accessTokenRoute.GET("/hmac/ping", handler.PingHandler, middleware2.SignatureHMACValidator(route.config))
	accessTokenRoute.POST("/rsa/proof", handler.VerifyProof, middleware2.RSASignatureValidator(route.config))
	accessTokenRoute.POST("/hmac/proof", handler.VerifyProof, middleware2.SignatureHMACValidator(route.config))
	accessTokenRoute.POST("/transaction/payment", handler.PaymentTransaction, middleware2.RSASignatureValidator(route.config))
	accessTokenRoute.POST("/transaction/payment-proof", handler.PaymentTransactionWithProof, middleware2.RSASignatureValidator(route.config))
}

func (route *Routes) setMiddleware(rGroup *echo.Group) {
	rGroup.Use(middleware.Recover())
	rGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXRealIP},
		AllowMethods: []string{http.MethodGet, http.MethodPut},
	}))
	rGroup.Use(middleware.BodyLimit("2M"))
}
