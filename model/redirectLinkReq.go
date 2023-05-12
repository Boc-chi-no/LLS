package model

type RedirectLinkReq struct {
	Hash     string `uri:"hash"     binding:"required,alphanum"`
	Password string `form:"pwd"     binding:"omitempty,alphanum,max=8"`
	Soft     bool   `form:"soft"    binding:"omitempty"`
	Detect   bool   `form:"detect"  binding:"omitempty"`
}
