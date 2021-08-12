package codec

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"
)

type BitSize uint8

const (
	Bit8  BitSize = 1
	Bit16 BitSize = 2
	Bit32 BitSize = 4
	Bit64 BitSize = 8
)

// HexReverse 16进制大小端反转
//
// 左侧补零，每两个字符进行反转
//
// "1234", 4 => "00001234" => "34120000"
func HexReverse(hex string, byteSize int) string {
	size := byteSize * 2
	strLen := len(hex)
	index := 0
	var c byte

	b := make([]byte, size)
	// 左侧补 0
	for i := 0; i < size; i++ {
		if size > strLen+i {
			b[i] = '0'
		} else {
			b[i] = hex[i+strLen-size]
		}
	}

	for i := 0; i < size/2; i++ {
		index = size - i
		if i%2 == 0 {
			index -= 2
		}
		c = b[i]
		b[i] = b[index]
		b[index] = c
	}

	return string(b)
}

// PadLeft 左侧补全
func PadLeft(s string, size int, b byte) string {
	strLen := len(s)
	if strLen >= size {
		return s
	}

	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		if i < size-strLen {
			buf[i] = b
		} else {
			buf[i] = s[strLen+i-size]
		}
	}

	return string(buf)
}

func Uint2Hex(i uint64, size BitSize) string {
	s := strconv.FormatUint(i, 16)

	return HexReverse(s, int(size))
}

func Uint8ToHex(i uint8) string {
	return Uint2Hex(uint64(i), Bit8)
}

func Uint16ToHex(i uint16) string {
	return Uint2Hex(uint64(i), Bit16)
}

func Uint32ToHex(i uint32) string {
	return Uint2Hex(uint64(i), Bit32)
}

func Uint64ToHex(i uint64) string {
	return Uint2Hex(i, Bit64)
}

func Hex2Uint(hex string, size BitSize) uint64 {
	s := int(size)
	i, err := strconv.ParseUint(HexReverse(hex, s), 16, s*8)
	if err != nil {
		return 0
	}
	return i
}

func Hex2Uint8(hex string) uint8 {
	return uint8(Hex2Uint(hex, Bit8))
}

func Hex2Uint16(hex string) uint16 {
	return uint16(Hex2Uint(hex, Bit16))
}

func Hex2Uint32(hex string) uint32 {
	return uint32(Hex2Uint(hex, Bit32))
}

func Hex2Uint64(hex string) uint64 {
	return uint64(Hex2Uint(hex, Bit64))
}

func BCD2Uint(s string, size BitSize) uint64 {
	i, err := strconv.ParseUint(s, 10, int(size)*8)
	if err != nil {
		return 0
	}
	return i
}

func BCD2Uint8(s string) uint8 {
	return uint8(BCD2Uint(s, Bit8))
}

func BCD2Uint16(s string) uint16 {
	return uint16(BCD2Uint(s, Bit16))
}

func BCD2Uint32(s string) uint32 {
	return uint32(BCD2Uint(s, Bit32))
}

func BCD2Uint64(s string) uint64 {
	return uint64(BCD2Uint(s, Bit64))
}

func Uint2BCD(i uint64, byteSize int) string {
	if byteSize == 0 {
		return strconv.FormatUint(i, 10)
	}

	return PadLeft(strconv.FormatUint(i, 10), byteSize*2, '0')
}

func Uint8ToBCD(i uint8) string {
	return Uint2BCD(uint64(i), int(Bit8))
}

// CP56Time2a time to CP56Time2a
func CP56Time2a(t time.Time) string {
	msec := t.Nanosecond()/int(time.Millisecond) + t.Second()*1000
	bytes := []byte{byte(msec), byte(msec >> 8), byte(t.Minute()), byte(t.Hour()),
		byte(t.Weekday()<<5) | byte(t.Day()), byte(t.Month()), byte(t.Year() - 2000)}
	return hex.EncodeToString(bytes)
}

// ParseCP56Time2a 7个八位位组二进制时间，读7个字节，返回时间
// The year is assumed to be in the 20th century.
// See IEC 60870-5-4 § 6.8 and IEC 60870-5-101 second edition § 7.2.6.18.
func ParseCP56Time2a(s string) time.Time {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return time.Unix(0, 0)
	}

	if len(bytes) < 7 || bytes[2]&0x80 == 0x80 {
		return time.Time{}
	}

	x := int(binary.LittleEndian.Uint16(bytes))
	msec := x % 1000
	sec := x / 1000
	min := int(bytes[2] & 0x3f)
	hour := int(bytes[3] & 0x1f)
	day := int(bytes[4] & 0x1f)
	month := time.Month(bytes[5] & 0x0f)
	year := 2000 + int(bytes[6]&0x7f)

	nsec := msec * int(time.Millisecond)
	return time.Date(year, month, day, hour, min, sec, nsec, time.Local)
}

// Crc16 modbus
func Crc16(s string) string {
	return Uint16ToHex(ChecksumMBus([]byte(s)))
}
