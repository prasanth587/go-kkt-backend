package dtos

type ResStruct struct {
	Status   string `json:"status" example:"SUCCESS" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"200" example:"500"`
	Message  string `json:"message" example:"pong" example:"could not connect to db"`
}

type Res500Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"500"`
	Message  string `json:"message" example:"could not connect to db"`
}

type Res400Struct struct {
	Status   string `json:"status" example:"FAILED"`
	HTTPCode int    `json:"httpCode" example:"400"`
	Message  string `json:"message" example:"Invalid param"`
}

type UserActionConfirmationReq struct {
	DeviceUUID          string `json:"deviceUUID"`
	TimeZone            string `json:"timeZone"`
	EmailID             string `json:"emailID"`
	WarningLabelRead    bool   `json:"warningLabelRead"`
	AccessCodeVerified  bool   `json:"accessCodeVerified"`
	TermsAndPrivacyRead bool   `json:"termsAndPrivacyRead"`
}

type UserActionConfirmationResponse struct {
	Message string `json:"message"`
}
