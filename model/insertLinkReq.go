package model

// InsertLinkReq This struct represents the payload to be posted to the shortener link
type InsertLinkReq struct {
	URL      string `json:"link"    binding:"required,url"`
	CAPTCHA  string `json:"captcha" binding:"required,alphanum,max=6"`
	PASSWORD string `json:"pwd"     binding:"omitempty,alphanum,max=8"`
	EXPIRE   int64  `json:"expire"  binding:"omitempty,numeric"`
	MEMO     string `json:"memo"    binding:"omitempty,max=32"`
}
