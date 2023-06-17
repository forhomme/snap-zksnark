package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"smart-contract-service/app/usecase"
	"smart-contract-service/configuration"
	"smart-contract-service/models"
)

const (
	EllipticAlgorithm = "elliptic"
	HashAlgorithm     = "hash"
	EddsaAlgorithm    = "eddsa"
)

type HTTP struct {
	config configuration.ConfigApp
	uc     usecase.InputPort
}

func NewHTTP(config configuration.ConfigApp, uc usecase.InputPort) *HTTP {
	return &HTTP{
		config: config,
		uc:     uc,
	}
}

func (h *HTTP) SignUp(c echo.Context) (err error) {
	request := new(models.LoginRequest)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	data, _ := json.Marshal(request)

	var buf bytes.Buffer
	err = json.Compact(&buf, data)
	id, err := h.uc.SignUp(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, models.Response{
		Code:    http.StatusCreated,
		Message: models.SUCCESS,
		Data:    id,
	})
}

func (h *HTTP) Token(c echo.Context) (err error) {
	request := new(models.TokenRequest)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	key, err := h.uc.TokenSign(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, models.Response{
		Code:    http.StatusCreated,
		Message: models.SUCCESS,
		Data:    key,
	})
}

func (h *HTTP) TokenHMAC(c echo.Context) (err error) {
	request := new(models.TokenRequest)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	key, err := h.uc.TokenHMAC(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, models.Response{
		Code:    http.StatusCreated,
		Message: models.SUCCESS,
		Data:    key,
	})
}

func (h *HTTP) Login(c echo.Context) (err error) {
	request := new(models.LoginRequest)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	data, err := h.uc.DoLogin(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: models.SUCCESS,
		Data:    data,
	})
}

func (h *HTTP) RefreshToken(c echo.Context) (err error) {
	request := new(models.RefreshTokenRequest)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	data, err := h.uc.RefreshToken(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: models.SUCCESS,
		Data:    data,
	})
}

func (h *HTTP) PingHandler(c echo.Context) (err error) {
	ping := models.Ping{
		Version: h.config.Version,
		Name:    h.config.AppName,
	}
	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: models.SUCCESS,
		Data:    ping,
	})
}

func (h *HTTP) GetProof(c echo.Context) (err error) {
	request := new(models.CustomerIdRequest)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var proof *models.ProofResponse
	switch algo := request.Algo; algo {
	case EllipticAlgorithm:
		proof, err = h.uc.GetEllipticProof(request.Id)
	case HashAlgorithm:
		proof, err = h.uc.GetHashProof(request.Id)
	case EddsaAlgorithm:
		proof, err = h.uc.GetEddsaProof(request.Id)
	default:
		err = fmt.Errorf("algorithm not found : %s", algo)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    proof,
	})
}

func (h *HTTP) VerifyProof(c echo.Context) (err error) {
	var request *models.ProofRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var valid bool
	switch algo := request.Algo; algo {
	case EllipticAlgorithm:
		_, valid = h.uc.VerifyEllipticProof(request.Proof)
	case HashAlgorithm:
		_, valid = h.uc.VerifyHashProof(request.Proof)
	case EddsaAlgorithm:
		_, valid = h.uc.VerifyEddsaProof(request.Proof)
	default:
		valid = false
	}

	if !valid {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: "Proof not knowledgeable",
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "Success",
	})
}

func (h *HTTP) PaymentTransaction(c echo.Context) (err error) {
	var request *models.PaymentTransactionRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	id, err := h.uc.PaymentTransaction(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: models.SUCCESS,
		Data:    id,
	})
}

func (h *HTTP) PaymentTransactionWithProof(c echo.Context) (err error) {
	var request *models.PaymentTransactionWithProofRequest
	if err = c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	var valid bool
	var userId string
	switch algo := request.Algo; algo {
	case EllipticAlgorithm:
		userId, valid = h.uc.VerifyEllipticProof(request.Proof)
	case HashAlgorithm:
		userId, valid = h.uc.VerifyHashProof(request.Proof)
	case EddsaAlgorithm:
		userId, valid = h.uc.VerifyEddsaProof(request.Proof)
	default:
		valid = false
	}

	if !valid {
		return c.JSON(http.StatusUnauthorized, models.Response{
			Code:    http.StatusUnauthorized,
			Message: errors.New("Proof not valid ").Error(),
		})
	}

	id, err := h.uc.PaymentTransaction(&models.PaymentTransactionRequest{
		PartnerReferenceNo: request.PartnerReferenceNo,
		CustomerId:         userId,
		Amount:             request.Amount,
		AdditionalInfo:     request.AdditionalInfo,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: models.SUCCESS,
		Data:    id,
	})
}
