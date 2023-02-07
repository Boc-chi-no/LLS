package model

type ManageLinkReq struct {
	Hash      string   `json:"hash"`
	CAPTCHA   string   `json:"captcha"`
	Token     string   `json:"token"`
	Page      int64    `json:"page"`
	Size      int64    `json:"size"`
}
