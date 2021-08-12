package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHexReverse(t *testing.T) {
	cases := []struct {
		input string
		size  int
		want  string
	}{
		{"", 2, "0000"},
		{"123", 2, "2301"},
		{"fafb", 2, "fbfa"},
		{"123456", 4, "56341200"},
	}

	for _, c := range cases {
		got := HexReverse(c.input, c.size)
		assert.Equal(t, c.want, got)
	}
}

func BenchmarkHexReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = HexReverse("123456", 4)
	}
}

func TestPadLeft(t *testing.T) {
	cases := []struct {
		input string
		size  int
		b     byte
		want  string
	}{
		{"", 2, '0', "00"},
		{"123", 4, '0', "0123"},
		{"fafb", 2, '0', "fafb"},
		{"a", 3, '+', "++a"},
	}

	for _, c := range cases {
		got := PadLeft(c.input, c.size, c.b)
		assert.Equal(t, c.want, got)
	}
}

func BenchmarkPadLeft(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = PadLeft("1234", 32, '0')
	}
}

func TestUint2Hex(t *testing.T) {
	cases := []struct {
		i    int
		size BitSize
		want string
	}{
		{0, Bit8, "00"},
		{1, Bit8, "01"},
		{10, Bit8, "0a"},
		{258, Bit16, "0201"},
		{258, Bit32, "02010000"},
		{16909060, Bit64, "0403020100000000"},
	}

	for _, c := range cases {
		got := Uint2Hex(uint64(c.i), c.size)
		assert.Equal(t, c.want, got)
	}
}

func BenchmarkUint2Hex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Uint2Hex(uint64(256), Bit64)
	}
}

func TestHex2Uint(t *testing.T) {
	cases := []struct {
		input string
		size  BitSize
		want  uint64
	}{
		{"", Bit8, 0},
		{"0a", Bit8, 10},
		{"0201", Bit8, 1},
		{"0201", Bit16, 258},
		{"02010000", Bit32, 258},
		{"0403020100000000", Bit64, 16909060},
	}

	for _, c := range cases {
		got := Hex2Uint(c.input, c.size)
		assert.Equal(t, c.want, got)
	}
}

func BenchmarkHex2Uint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Hex2Uint("12345678", Bit64)
	}
}

func TestBCD2Uint(t *testing.T) {
	cases := []struct {
		input string
		size  BitSize
		want  uint64
	}{
		{"", Bit8, 0},
		{"01", Bit8, 1},
		{"0257", Bit8, 0},
		{"0201", Bit16, 201},
	}

	for _, c := range cases {
		got := BCD2Uint(c.input, c.size)
		assert.Equal(t, c.want, got, c.input)
	}
}

func TestUint2BCD(t *testing.T) {
	cases := []struct {
		input uint64
		size  int
		want  string
	}{
		{0, 0, "0"},
		{0, 1, "00"},
		{12, 1, "12"},
		{99, 4, "00000099"},
	}

	for _, c := range cases {
		got := Uint2BCD(c.input, c.size)
		assert.Equal(t, c.want, got, c)
	}
}

func BenchmarkUint2BCD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Uint2BCD(1234567890, 10)
	}
}

func TestUint8ToBCD(t *testing.T) {
	cases := []struct {
		input uint8
		want  string
	}{
		{0, "00"},
		{1, "01"},
		{12, "12"},
	}

	for _, c := range cases {
		got := Uint8ToBCD(c.input)
		assert.Equal(t, c.want, got, c)
	}
}

func BenchmarkUint8ToBCD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Uint8ToBCD(1)
	}
}
