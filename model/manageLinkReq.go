package model

type ManageLinkReq struct {
	Hash    string `json:"hash" binding:"required,alphanum"`
	CAPTCHA string `json:"captcha" binding:"required,alphanum"`
	Token   string `json:"token" binding:"required,alphanum"`
	Page    int64  `json:"page" binding:"numeric"`
	Size    int64  `json:"size" binding:"numeric"`
}
