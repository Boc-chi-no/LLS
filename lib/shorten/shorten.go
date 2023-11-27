package shorten

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/spaolacci/murmur3"
	"linkshortener/lib/tool"
	"linkshortener/model"
	"linkshortener/setting"
	"strconv"
	"time"
)

// GenerateShortenLink This method generates the hash
func GenerateShortenLink(req model.InsertLinkReq) model.Link {
	now := time.Now()
	sec := now.Unix()
	nanoSecStr := strconv.FormatInt(now.UnixNano(), 16)

	count := tool.GlobalCounterSafeAdd(1)
	countStr := strconv.FormatUint(count, 16)

	murmurHash := murmur3.Sum32WithSeed([]byte(tool.ConcatStrings(req.URL, ":", nanoSecStr, ":", countStr)), setting.Cfg.Seed)
	hex62Hash := tool.Uint32ToBase62String(murmurHash)

	var link model.Link
	link.Created = sec
	link.Token, _ = tool.GetToken(16)
	link.ShortHash = hex62Hash
	link.URL = req.URL
	link.Memo = req.MEMO
	link.Expire = req.EXPIRE
	if req.PASSWORD != "" {
		passwordHash := sha256.Sum256([]byte(tool.ConcatStrings(link.ShortHash, req.PASSWORD, tool.Uint32ToBase62String(setting.Cfg.Seed))))
		link.Password = hex.EncodeToString(passwordHash[:])
	}
	link.Delete = false

	return link
}
