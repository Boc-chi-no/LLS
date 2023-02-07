package model

// This struct represents the payload to be posted to the shortener link
type InsertLinkReq struct {
	URL     string `json:"link"`
	CAPTCHA string `json:"captcha"`
}
