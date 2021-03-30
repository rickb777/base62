package base62

import (
	"math/bits"
	"reflect"
	"strconv"
	"unsafe"
)

const (
	base        = 62
	compactMask = 0x1E // 00011110
	mask5bits   = 0x1F // 00011111
)

const encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// An Encoding is a radix 62 encoding/decoding scheme, defined by a 62-character alphabet.
type Encoding struct {
	encode    [base]byte
	decodeMap [256]byte
}

// NewEncoding returns a new padded Encoding defined by the given alphabet. It will
// panic if the alphabet is not 62 bytes long.
func NewEncoding(encoder string) *Encoding {
	if len(encoder) != base {
		panic("encoding alphabet is not 62-bytes long")
	}
	for i := 0; i < len(encoder); i++ {
		if encoder[i] == '\n' || encoder[i] == '\r' {
			panic("encoding alphabet contains newline character")
		}
	}

	e := new(Encoding)
	copy(e.encode[:], encoder)
	for i := 0; i < len(e.decodeMap); i++ {
		e.decodeMap[i] = 0xFF
	}
	for i := 0; i < len(encoder); i++ {
		e.decodeMap[encoder[i]] = byte(i)
	}
	return e
}

// StdEncoding is an encoding/decoding scheme consisting of [A-Za-z0-9].
var StdEncoding = NewEncoding(encodeStd)

// Encode encodes src using the encoding enc.
func (enc *Encoding) Encode(dest, src []byte) int {
	if len(src) == 0 {
		return 0
	}
	encoder := newEncoder(src)
	return encoder.encode(dest, enc.encode[:])
}

// EncodeToString returns the encoding of src.
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := allocEncodeBuffer(len(src))
	n := enc.Encode(buf, src)
	return b2s(buf[:n])
}

type encoder struct {
	b   []byte
	pos int
}

func newEncoder(src []byte) *encoder {
	return &encoder{
		b:   src,
		pos: len(src) * 8,
	}
}

func (enc *encoder) next() (byte, bool) {
	var i, pos int
	var j, blen byte
	pos = enc.pos - 6
	if pos <= 0 {
		pos = 0
		blen = byte(enc.pos)
	} else {
		i = pos / 8
		j = byte(pos % 8)
		blen = byte((i+1)*8 - pos)
		if blen > 6 {
			blen = 6
		}
	}
	shift := 8 - j - blen
	b := enc.b[i] >> shift & (1<<blen - 1)

	if blen < 6 && pos > 0 {
		blen1 := 6 - blen
		b = b<<blen1 | enc.b[i+1]>>(8-blen1)
	}
	if b&compactMask == compactMask {
		if pos > 0 || b > mask5bits {
			pos++
		}
		b &= mask5bits
	}
	enc.pos = pos

	return b, pos > 0
}

func (enc *encoder) encode(ret, encTable []byte) int {
	i := 0
	x, hasMore := enc.next()
	for {
		ret[i] = encTable[x]
		i++
		if !hasMore {
			break
		}
		x, hasMore = enc.next()
	}
	return i
}

func allocEncodeBuffer(srcLen int) []byte {
	return make([]byte, srcLen*8/5+2)
}

//-------------------------------------------------------------------------------------------------

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base62 data at input byte " + strconv.FormatInt(int64(e), 10)
}

//-------------------------------------------------------------------------------------------------

// Decode decodes src using the encoding enc.
func (enc *Encoding) Decode(dest, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	n, err := decode(dest, src, enc.decodeMap[:])
	if err == nil && n < len(dest) {
		copy(dest, dest[len(dest)-n:])
	}
	return n, err
}

// DecodeString returns the bytes represented by the base62 string s.
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	return enc.decodeBytes(s2b(s))
}

func (enc *Encoding) decodeBytes(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return []byte{}, nil
	}
	buf := allocDecodeBuffer(len(src))
	n, err := decode(buf, src, enc.decodeMap[:])
	return buf[len(buf)-n:], err
}

func decode(ret, src, decTable []byte) (n int, err error) {
	idx := len(ret)
	pos := byte(0)
	b := 0
	for i, c := range src {
		x := decTable[c]
		if x == 0xFF {
			return 0, CorruptInputError(i)
		}
		if i == len(src)-1 {
			b |= int(x) << pos
			pos += byte(bits.Len8(x))
		} else if x&compactMask == compactMask {
			b |= int(x) << pos
			pos += 5
		} else {
			b |= int(x) << pos
			pos += 6
		}
		if pos >= 8 {
			idx--
			ret[idx] = byte(b)
			n++
			pos %= 8
			b >>= 8
		}
	}
	if pos > 0 {
		idx--
		ret[idx] = byte(b)
		n++
	}

	return n, nil
}

func allocDecodeBuffer(srcLen int) []byte {
	return make([]byte, srcLen*6/8+1)
}

//-------------------------------------------------------------------------------------------------

// EncodeToBytes encodes src using StdEncoding.
func EncodeToBytes(src []byte) []byte {
	buf := allocEncodeBuffer(len(src))
	n := StdEncoding.Encode(buf, src)
	return buf[:n]
}

// EncodeToString encodes src using StdEncoding.
func EncodeToString(src []byte) string {
	return StdEncoding.EncodeToString(src)
}

// Decode decodes src using StdEncoding.
func DecodeBytes(src []byte) ([]byte, error) {
	return StdEncoding.decodeBytes(src)
}

// Decode decodes src using StdEncoding.
func DecodeString(src string) ([]byte, error) {
	return StdEncoding.DecodeString(src)
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := &reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(bh))
}
