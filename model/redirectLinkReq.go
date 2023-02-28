package model

type RedirectLinkReq struct {
	Hash string `uri:"hash" binding:"required,alphanum"`
}
