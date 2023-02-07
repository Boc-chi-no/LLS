package shorten

import (
	"github.com/spaolacci/murmur3"
	"linkshortener/lib/tool"
	"linkshortener/model"
	"strconv"
	"time"
)

var seed uint32 = uint32(10011011)
var Hex62Map map[int]string = map[int]string{
	0: "0",
	1: "1",
	2: "2",
	3: "3",
	4: "4",
	5: "5",
	6: "6",
	7: "7",
	8: "8",
	9: "9",
	10: "a",
	11: "b",
	12: "c",
	13: "d",
	14: "e",
	15: "f",
	16: "g",
	17: "h",
	18: "i",
	19: "j",
	20: "k",
	21: "l",
	22: "m",
	23: "n",
	24: "o",
	25: "p",
	26: "q",
	27: "r",
	28: "s",
	29: "t",
	30: "u",
	31: "v",
	32: "w",
	33: "x",
	34: "y",
	35: "z",
	36: "A",
	37: "B",
	38: "C",
	39: "D",
	40: "E",
	41: "F",
	42: "G",
	43: "H",
	44: "I",
	45: "J",
	46: "K",
	47: "L",
	48: "M",
	49: "N",
	50: "O",
	51: "P",
	52: "Q",
	53: "R",
	54: "S",
	55: "T",
	56: "U",
	57: "V",
	58: "W",
	59: "X",
	60: "Y",
	61: "Z"}

func uint32ToHex62(uitNum uint32) string {
	hex62Str := ""
	num := int(uitNum)
	var remainder int
	var remainderString string
	for num != 0 {
		remainder = num % 62
		if 76 > remainder && remainder > 9 {
			remainderString = Hex62Map[remainder]
		} else {
			remainderString = strconv.Itoa(remainder)
		}
		hex62Str = remainderString + hex62Str
		num = num / 62
	}
	return hex62Str
}
// This method generates the hash
func ShortenLink(url model.InsertLinkReq) model.Link {
	now := time.Now()
	sec := now.Unix()
	nsecStr := strconv.FormatInt(now.UnixNano(), 16)
	murmurHash := murmur3.Sum32WithSeed([]byte(url.URL + nsecStr), seed)
	hex62Hash := uint32ToHex62(murmurHash)


	var link model.Link
	link.Created = sec
	link.Token = tool.GetToken(16)
	link.ShortHash = hex62Hash
	link.URL = url.URL
	link.Delete = false

	return link
}
