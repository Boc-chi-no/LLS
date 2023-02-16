package model

// InsertLinkReq This struct represents the payload to be posted to the shortener link
type InsertLinkReq struct {
	URL     string `json:"link" binding:"required,url"`
	CAPTCHA string `json:"captcha" binding:"required,alphanum"`
}
