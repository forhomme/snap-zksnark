package models

type CustomerIdRequest struct {
	Algo string `query:"algo"`
	Id   string `query:"id"`
}

type ProofRequest struct {
	Algo  string `json:"algo"`
	Proof string `json:"proof"`
}

type RequestHeader struct {
	Authorization string `json:"-"`
	Signature     string `json:"-"`
	PartnerId     string `json:"-"`
	ChannelId     string `json:"-"`
	DeviceId      string `json:"-"`
	Timestamp     string `json:"-"`
	ExternalId    string `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	Username     string `json:"username"`
	RefreshToken string `json:"refreshToken"`
}
type TokenRequest struct {
	Method      string `json:"method"`
	Endpoint    string `json:"endpoint"`
	Timestamp   string `json:"timestamp,omitempty"` //YYYY:MM:ddThh:mm:ss.SSSTZD
	BodyCompact string `json:"bodyCompact"`
}

type PaymentTransactionProofRequest struct {
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	Proof              string `json:"proof"`
	Amount             struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	AdditionalInfo struct {
		DeviceId string `json:"deviceId,omitempty"`
		Channel  string `json:"channel,omitempty"`
	} `json:"additionalInfo,omitempty"`
}

type PaymentTransactionRequest struct {
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	CustomerNumber     string `json:"customerNumber"`
	CustomerId         string `json:"customerId"`
	Amount             struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	AdditionalInfo struct {
		DeviceId string `json:"deviceId,omitempty"`
		Channel  string `json:"channel,omitempty"`
	} `json:"additionalInfo,omitempty"`
}

type PaymentTransactionWithProofRequest struct {
	PartnerReferenceNo string `json:"partnerReferenceNo"`
	Algo               string `json:"algo"`
	Proof              string `json:"proof"`
	Amount             struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	AdditionalInfo struct {
		DeviceId string `json:"deviceId,omitempty"`
		Channel  string `json:"channel,omitempty"`
	} `json:"additionalInfo,omitempty"`
}
