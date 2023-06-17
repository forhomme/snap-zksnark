package usecase

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	bn254 "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	eddsa2 "github.com/consensys/gnark-crypto/signature/eddsa"
	"github.com/consensys/gnark/backend/groth16"
	witness2 "github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/golang-jwt/jwt"
	"smart-contract-service/configuration"
	"smart-contract-service/internal"
	"smart-contract-service/models"
	models2 "smart-contract-service/models/circuit"
	"strings"
	"time"
)

type Usecase struct {
	redis RedisRepository
	db    DbRepository
	cfg   configuration.ConfigApp
}

func NewUsecase(redis RedisRepository, db DbRepository, cfg configuration.ConfigApp) *Usecase {
	return &Usecase{
		redis: redis,
		db:    db,
		cfg:   cfg,
	}
}

type InputPort interface {
	DoLogin(input *models.LoginRequest) (out *models.LoginResponse, err error)
	SignUp(input *models.LoginRequest) (id string, err error)
	RefreshToken(input *models.RefreshTokenRequest) (out *models.LoginResponse, err error)
	TokenSign(input *models.TokenRequest) (out string, err error)
	TokenHMAC(input *models.TokenRequest) (out string, err error)
	GetEllipticProof(id string) (data *models.ProofResponse, err error)
	GetHashProof(id string) (data *models.ProofResponse, err error)
	GetEddsaProof(id string) (data *models.ProofResponse, err error)
	VerifyEllipticProof(code string) (string, bool)
	VerifyHashProof(code string) (string, bool)
	VerifyEddsaProof(code string) (string, bool)
	PaymentTransaction(in *models.PaymentTransactionRequest) (id string, err error)
}

type DbRepository interface {
	GetCustomerData(id string) (data *models.Customer, err error)
	GetCustomerByAccount(account string) (data *models.Customer, err error)
	GetUserById(id string) (data *models.Partners, err error)
	GetUserByUsername(username string) (data *models.Partners, err error)
	GetUserByReferenceNo(referenceNo string) (data *models.Partners, err error)
	InsertUser(input *models.Partners) (id string, err error)
	InsertPayment(input *models.Payment) (id string, err error)
}

type RedisRepository interface {
	Set(key string, data string) (err error)
	Get(key string) (val string, err error)
}

func (u *Usecase) DoLogin(input *models.LoginRequest) (out *models.LoginResponse, err error) {
	if len(input.Username) == 0 {
		err = fmt.Errorf("please input email or username.")
		return
	}
	user, err := u.db.GetUserByUsername(input.Username)
	if err != nil {
		return
	}

	if !internal.CheckPasswordHash(input.Password, user.Password) {
		err = fmt.Errorf("password invalid.")
		return
	}
	out, err = u.generateToken(user)
	return
}

func (u *Usecase) SignUp(input *models.LoginRequest) (id string, err error) {
	if len(input.Username) == 0 {
		err = fmt.Errorf("please input username.")
		return
	}
	if len(input.Password) == 0 {
		err = fmt.Errorf("please input password.")
		return
	}

	user, err := u.db.GetUserByUsername(input.Username)
	if err != nil {
		return
	}
	if user.Username != "" {
		err = fmt.Errorf("%s: %s", "Username already taken", user.Username)
		return
	}
	hashPasword, _ := internal.HashPassword(input.Password)
	id, err = u.db.InsertUser(&models.Partners{
		Username: input.Username,
		Password: hashPasword,
	})
	return
}

func (u *Usecase) RefreshToken(input *models.RefreshTokenRequest) (out *models.LoginResponse, err error) {
	if len(input.Username) == 0 {
		err = fmt.Errorf("please input username.")
		return
	}

	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.cfg.Secret), nil
	})
	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user, errUser := u.db.GetUserById(claims["ID"].(string))
		if errUser != nil {
			return nil, errUser
		}
		out, err = u.generateToken(user)
		if err != nil {
			return nil, err
		}
	} else {
		return
	}
	return
}

func (u *Usecase) TokenSign(input *models.TokenRequest) (out string, err error) {
	_, err = time.Parse("2006-01-02T15:04:05.999TZ7", input.Timestamp)
	if err != nil {
		err = fmt.Errorf("%s : %s", "Wrong timestamp", err.Error())
		return "", err
	}

	privKey, err := internal.GeneratePrivateKey(u.cfg)
	if err != nil {
		return
	}
	dataRequest := []byte(input.BodyCompact)

	stringToSign := fmt.Sprintf("%s:%s:%s:%s", input.Method, input.Endpoint, strings.ToLower(internal.SignHMAC256(dataRequest,
		[]byte(u.cfg.Secret))), input.Timestamp)

	digest := sha256.Sum256([]byte(stringToSign))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, digest[:])
	if err != nil {
		return
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (u *Usecase) TokenHMAC(input *models.TokenRequest) (out string, err error) {
	_, err = time.Parse("2006-01-02T15:04:05.999TZ7", input.Timestamp)
	if err != nil {
		err = fmt.Errorf("%s : %s", "Wrong timestamp", err.Error())
		return "", err
	}

	dataRequest := []byte(input.BodyCompact)
	stringToSign := fmt.Sprintf("%s:%s:%s:%s", input.Method, input.Endpoint, strings.ToLower(internal.SignHMAC256(dataRequest,
		[]byte(u.cfg.Secret))), input.Timestamp)
	h := hmac.New(sha512.New, []byte(u.cfg.Secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func (u *Usecase) GetEllipticProof(id string) (data *models.ProofResponse, err error) {
	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	pk := groth16.NewProvingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csEllipticPath)
	internal.Deserialize(pk, internal.PkEllipticPath)

	cData, err := u.db.GetCustomerData(id)
	if err != nil {
		return nil, err
	}

	assignment := &models2.EllipticCurve{}
	x := len(cData.Name)
	assignment.X = frontend.Variable(x)
	assignment.Y = frontend.Variable(x*x*x + x + 5)

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return
	}

	var proofBuf bytes.Buffer
	proof.WriteTo(&proofBuf)

	publicWitness, _ := witness.Public()
	dataBin, _ := publicWitness.MarshalBinary()
	dataResponse := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s||%s", string(dataBin), cData.Id)))

	err = u.redis.Set(fmt.Sprintf("%s_%s", dataResponse, "proof"), proofBuf.String())
	if err != nil {
		return nil, err
	}
	err = u.redis.Set(fmt.Sprintf("%s_%s", dataResponse, "witness"), string(dataBin))
	if err != nil {
		return nil, err
	}
	data = &models.ProofResponse{Hash: dataResponse}

	return
}

func (u *Usecase) GetHashProof(id string) (data *models.ProofResponse, err error) {
	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	pk := groth16.NewProvingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csPath)
	internal.Deserialize(pk, internal.PkPath)

	cData, err := u.db.GetCustomerData(id)
	if err != nil {
		return nil, err
	}
	//marshalData := []byte(internal.StringWithCharset(len(cData.Id), cData.Id))
	marshalData, _ := json.Marshal(internal.StringWithCharset(len(cData.Id), cData.Id))

	assignment := &models2.Circuit{}
	b := make([]byte, 32)
	copy(b, marshalData)
	hash := internal.MimcHash(b)

	assignment.Secret = frontend.Variable(b)
	assignment.Hash = frontend.Variable(hash)

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return
	}

	var proofBuf bytes.Buffer
	proof.WriteTo(&proofBuf)

	publicWitness, _ := witness.Public()
	dataBin, _ := publicWitness.MarshalBinary()
	dataResponse := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s||%s", string(dataBin), cData.Id)))

	err = u.redis.Set(fmt.Sprintf("%s_%s", dataResponse, "proof"), proofBuf.String())
	if err != nil {
		return nil, err
	}
	err = u.redis.Set(fmt.Sprintf("%s_%s", dataResponse, "witness"), string(dataBin))
	if err != nil {
		return nil, err
	}
	data = &models.ProofResponse{Hash: dataResponse}

	return
}

func (u *Usecase) GetEddsaProof(id string) (data *models.ProofResponse, err error) {
	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	pk := groth16.NewProvingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csEddsaPath)
	internal.Deserialize(pk, internal.PkEddsaPath)

	cData, err := u.db.GetCustomerData(id)
	if err != nil {
		return nil, err
	}

	// instantiate hash function
	f := bn254.NewMiMC()

	// create a eddsa key pair
	privateKey, err := eddsa2.New(tedwards.BN254, rand.Reader)
	publicKey := privateKey.Public()

	marshalData, _ := json.Marshal(internal.StringWithCharset(len(cData.Id), cData.Id))
	b := make([]byte, 32)
	copy(b, marshalData)

	// sign the message
	signature, err := privateKey.Sign(b, f)

	// verifies signature
	isValid, err := publicKey.Verify(signature, b, f)
	if !isValid {
		return nil, errors.New("not valid")
	}
	// declare the witness
	assignment := &models2.EddsaCircuit{}

	// assign message value
	assignment.Message = b

	// public key bytes
	_publicKey := publicKey.Bytes()

	// assign public key values
	assignment.PublicKey.Assign(tedwards.BN254, _publicKey[:32])

	// assign signature values
	assignment.Signature.Assign(tedwards.BN254, signature)

	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return
	}

	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return
	}

	var proofBuf bytes.Buffer
	proof.WriteTo(&proofBuf)

	publicWitness, _ := witness.Public()
	dataBin, _ := publicWitness.MarshalBinary()
	dataResponse := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s||%s", string(dataBin), cData.Id)))

	err = u.redis.Set(fmt.Sprintf("%s_%s", dataResponse, "proof"), proofBuf.String())
	if err != nil {
		return nil, err
	}
	err = u.redis.Set(fmt.Sprintf("%s_%s", dataResponse, "witness"), string(dataBin))
	if err != nil {
		return nil, err
	}
	data = &models.ProofResponse{Hash: dataResponse}

	return
}

func (u *Usecase) VerifyEllipticProof(code string) (string, bool) {
	decodeString, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", false
	}
	if !strings.Contains(string(decodeString), "||") {
		return "", false
	}
	decodeArray := strings.Split(string(decodeString), "||")
	cData, err := u.db.GetCustomerData(decodeArray[1])
	if err != nil {
		return "", false
	}
	if cData.Id == "" {
		return "", false
	}

	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	vk := groth16.NewVerifyingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csEllipticPath)
	internal.Deserialize(vk, internal.VkEllipticPath)

	// get proof
	val, err := u.redis.Get(fmt.Sprintf("%s_%s", code, "proof"))
	if err != nil {
		return "", false
	}
	valReader := strings.NewReader(val)
	proof := groth16.NewProof(ecc.BN254)
	proof.ReadFrom(valReader)

	// get witness
	public, err := u.redis.Get(fmt.Sprintf("%s_%s", code, "witness"))
	if err != nil {
		return "", false
	}
	witness, _ := witness2.New(ecc.BN254.ScalarField())
	err = witness.UnmarshalBinary([]byte(public))
	if err != nil {
		return "", false
	}
	// verify the proof using witness
	err = groth16.Verify(proof, vk, witness)
	if err != nil {
		return "", false
	}
	return cData.Id, true
}

func (u *Usecase) VerifyHashProof(code string) (string, bool) {
	decodeString, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", false
	}
	if !strings.Contains(string(decodeString), "||") {
		return "", false
	}
	decodeArray := strings.Split(string(decodeString), "||")
	cData, err := u.db.GetCustomerData(decodeArray[1])
	if err != nil {
		return "", false
	}
	if cData.Id == "" {
		return "", false
	}

	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	vk := groth16.NewVerifyingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csPath)
	internal.Deserialize(vk, internal.VkPath)

	// get proof
	val, err := u.redis.Get(fmt.Sprintf("%s_%s", code, "proof"))
	if err != nil {
		return "", false
	}
	valReader := strings.NewReader(val)
	proof := groth16.NewProof(ecc.BN254)
	proof.ReadFrom(valReader)

	// get witness
	public, err := u.redis.Get(fmt.Sprintf("%s_%s", code, "witness"))
	if err != nil {
		return "", false
	}
	witness, _ := witness2.New(ecc.BN254.ScalarField())
	err = witness.UnmarshalBinary([]byte(public))
	if err != nil {
		return "", false
	}
	// verify the proof using witness
	err = groth16.Verify(proof, vk, witness)
	if err != nil {
		return "", false
	}
	return cData.Id, true
}

func (u *Usecase) VerifyEddsaProof(code string) (string, bool) {
	decodeString, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", false
	}
	if !strings.Contains(string(decodeString), "||") {
		return "", false
	}
	decodeArray := strings.Split(string(decodeString), "||")
	cData, err := u.db.GetCustomerData(decodeArray[1])
	if err != nil {
		return "", false
	}
	if cData.Id == "" {
		return "", false
	}

	// read R1CS, proving key and verifying keys
	ccs := groth16.NewCS(ecc.BN254)
	vk := groth16.NewVerifyingKey(ecc.BN254)
	internal.Deserialize(ccs, internal.R1csEddsaPath)
	internal.Deserialize(vk, internal.VkEddsaPath)

	// get proof
	val, err := u.redis.Get(fmt.Sprintf("%s_%s", code, "proof"))
	if err != nil {
		return "", false
	}
	valReader := strings.NewReader(val)
	proof := groth16.NewProof(ecc.BN254)
	proof.ReadFrom(valReader)

	// get witness
	public, err := u.redis.Get(fmt.Sprintf("%s_%s", code, "witness"))
	if err != nil {
		return "", false
	}
	witness, _ := witness2.New(ecc.BN254.ScalarField())
	err = witness.UnmarshalBinary([]byte(public))
	if err != nil {
		return "", false
	}
	// verify the proof using witness
	err = groth16.Verify(proof, vk, witness)
	if err != nil {
		return "", false
	}
	return cData.Id, true
}

func (u *Usecase) PaymentTransaction(in *models.PaymentTransactionRequest) (id string, err error) {
	data := &models.Customer{}
	if in.CustomerNumber != "" {
		data, err = u.db.GetCustomerByAccount(in.CustomerNumber)
		if err != nil {
			return
		}
	} else {
		data, err = u.db.GetCustomerData(in.CustomerId)
		if err != nil {
			return
		}
	}
	if data == (&models.Customer{}) {
		return "", fmt.Errorf("Customer not found ")
	}

	dataPartner, err := u.db.GetUserByReferenceNo(in.PartnerReferenceNo)
	if err != nil {
		return
	}
	if dataPartner == (&models.Partners{}) {
		return "", fmt.Errorf("Partner not found ")
	}

	addInfo := make(map[string]interface{})
	addInfo["deviceId"] = in.AdditionalInfo.DeviceId
	addInfo["channel"] = in.AdditionalInfo.Channel
	id, err = u.db.InsertPayment(&models.Payment{
		PartnerId:      dataPartner.Id,
		ConsumerId:     data.Id,
		Amount:         in.Amount.Value,
		Currency:       in.Amount.Currency,
		AdditionalInfo: addInfo,
	})
	return
}

func (u *Usecase) generateToken(user *models.Partners) (out *models.LoginResponse, err error) {
	expire := time.Now().Add(time.Hour * time.Duration(u.cfg.Expire))
	rtExpire := time.Now().Add(time.Hour * (time.Duration(u.cfg.Expire) * 24))

	//1. set the jwt token
	claims := &models.JwtCustomClaims{
		ID:       user.Id,
		Username: user.Username,
		LoginAt:  time.Now().Format(time.RFC3339),
		ExpireAt: expire.Format(time.RFC3339),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(u.cfg.Secret))
	if err != nil {
		return
	}

	//2. set refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = user.Id
	rtClaims["expireAt"] = rtExpire.Format(time.RFC3339)
	rt, err := refreshToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	//3. set to response api
	out = &models.LoginResponse{
		AccessToken:     t,
		AccessExpireAt:  expire.Format(time.RFC3339),
		RefreshToken:    rt,
		RefreshExpireAt: rtExpire.Format(time.RFC3339),
	}
	return
}
