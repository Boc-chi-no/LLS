package ip2location

import (
	"encoding/binary"
	"linkshortener/log"
	"linkshortener/model"
	"net"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	// IndexLen Index length
	IndexLen = 7
	// RedirectMode1 Type of country, pointing to another point
	RedirectMode1 = 0x01
	// RedirectMode2 Type of country, pointing to a point
	RedirectMode2 = 0x02
)

type fileData struct {
	Data  []byte
	IPNum int64
}

// QQwry ip library
// TODO:Replacement required
type QQwry struct {
	Data   *fileData
	Offset int64
}

// IPData Data from the IP library
var IPData fileData

// InitIPData Initialise ip library data into memory
func (f *fileData) InitIPData(data []byte) (rs interface{}) {
	var (
		tmpData []byte
		err     error
	)
	log.InfoPrint("Get Online IP Location...")
	tmpData, err = GetOnline()
	if err != nil {
		tmpData = data
		log.WarnPrint("Get Online IP Location failed: %s", err)
	}
	f.Data = tmpData

	buf := f.Data[0:8]
	start := binary.LittleEndian.Uint32(buf[:4])
	end := binary.LittleEndian.Uint32(buf[4:])

	f.IPNum = int64((end-start)/IndexLen + 1)

	return true
}

// NewQQwry New qqwry type
func NewQQwry() QQwry {
	return QQwry{
		Data: &IPData,
	}
}

// ReadData Reading data from a file
func (q *QQwry) ReadData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.SetOffset(offset[0])
	}
	nums := int64(num)
	end := q.Offset + nums
	dataNum := int64(len(q.Data.Data))
	if q.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = q.Data.Data[q.Offset:end]
	q.Offset = end
	return
}

// SetOffset Set the offset
func (q *QQwry) SetOffset(offset int64) {
	q.Offset = offset
}

// Find ip address search corresponding to attribution information
func (q *QQwry) Find(ip string) (res model.Location) {

	res = model.Location{}

	if strings.Count(ip, ".") != 3 {
		return res
	}
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return
	}

	var country []byte
	var area []byte

	mode := q.readMode(offset + 4)
	if mode == RedirectMode1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == RedirectMode2 {
			c := q.readUInt24()
			country = q.readString(c)
			countryOffset += 4
		} else {
			country = q.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = q.readArea(countryOffset)
	} else if mode == RedirectMode2 {
		countryOffset := q.readUInt24()
		country = q.readString(countryOffset)
		area = q.readArea(offset + 8)
	} else {
		country = q.readString(offset + 4)
		area = q.readArea(offset + uint32(5+len(country)))
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	res.Country, _ = enc.String(string(country))
	res.Area, _ = enc.String(string(area))

	return
}

// readMode Get offset value type
func (q *QQwry) readMode(offset uint32) byte {
	mode := q.ReadData(1, int64(offset))
	return mode[0]
}

// readArea Reading area
func (q *QQwry) readArea(offset uint32) []byte {
	mode := q.readMode(offset)
	if mode == RedirectMode1 || mode == RedirectMode2 {
		areaOffset := q.readUInt24()
		if areaOffset == 0 {
			return []byte("")
		}
		return q.readString(areaOffset)
	}
	return q.readString(offset)
}

// readString Get String
func (q *QQwry) readString(offset uint32) []byte {
	q.SetOffset(int64(offset))
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = q.ReadData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

// searchIndex Find an index location
func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.ReadData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, IndexLen)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = q.getMiddleOffset(start, end)
		buf = q.ReadData(IndexLen, int64(mid))
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == IndexLen {
			offset := byteToUInt32(buf[4:])
			buf = q.ReadData(IndexLen)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		// Found larger, moving forward
		if _ip > ip {
			end = mid
		} else if _ip < ip { // Found smaller, moving backwards
			start = mid
		} else if _ip == ip {
			return byteToUInt32(buf[4:])
		}
	}
}

// readUInt24
func (q *QQwry) readUInt24() uint32 {
	buf := q.ReadData(3)
	return byteToUInt32(buf)
}

// getMiddleOffset
func (q *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / IndexLen) >> 1
	return start + records*IndexLen
}

// byteToUInt32 Converts byte to uint32
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
