package base62

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	mathrand "math/rand"
	"strings"
	"testing"
)

func Test_MethodEncodeDecode(t *testing.T) {
	for _, c := range []string{
		"",
		"Hello, 世界！",
	} {
		src := []byte(c)
		enc := make([]byte, len(src)*2)
		n1 := StdEncoding.Encode(enc, src)

		got := make([]byte, len(src)*2)
		n2, err := StdEncoding.Decode(got, enc[:n1])
		if err != nil {
			t.Fatalf("%s failed decode, err = %v", c, err)
		}
		if n2 != len(src) {
			t.Fatalf("%s failed decode, got = %d, want = %d", c, n2, len(src))
		}
		if !bytes.Equal(src, got[:n2]) {
			t.Fatalf("%s failed decode, got = %v, want = %v", c, got, src)
		}
	}
}

func Test_EncodeDecode(t *testing.T) {
	for _, c := range []string{
		"",
		"Hello, 世界！",
	} {
		src := []byte(c)
		dst := EncodeToBytes(src)
		got, err := DecodeBytes(dst)
		if err != nil {
			t.Fatalf("failed decode, err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}
	}
}

func Test_EncodeDecodeString(t *testing.T) {
	for _, c := range []string{
		"",
		"Hello, 世界！",
	} {
		src := []byte(c)
		dstStr := EncodeToString(src)
		t.Logf("Base62 vs Base64:\n%s\n%s\n", dstStr, base64.StdEncoding.EncodeToString(src))
		got, _ := DecodeString(dstStr)
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode string, got = %v, want = %v", got, src)
		}
	}
}

func Test_EncodeDecode2(t *testing.T) {
	src := []byte("http://our-uploads.s3.amazonaws.com/file-export/stuff-1427217700-12.csv?AWSAccessKeyId=AAAAIIG7MAQRTHTD7CLP&Expires=1427304113&Signature=VQsRAhgamiw1RVtbrCXOsMu%2BgFo")
	dst := EncodeToBytes(src)
	got, err := DecodeBytes(dst)
	if err != nil {
		t.Fatalf("failed decode, err = %v", err)
	}
	if !bytes.Equal(src, got) {
		t.Fatalf("failed decode, got = %v, want = %v", got, src)
	}

	dstStr := EncodeToString(src)
	t.Logf("Base62 vs Base64:\n%s\n%s\n", dstStr, base64.StdEncoding.EncodeToString(src))
	got, _ = DecodeString(dstStr)
	if !bytes.Equal(src, got) {
		t.Fatalf("failed decode string, got = %v, want = %v", got, src)
	}
}

func Test_EncodeDecode_Zeros(t *testing.T) {
	for i := 0; i < 1000; i++ {
		src := make([]byte, i)
		dst := EncodeToBytes(src)
		got, err := DecodeBytes(dst)
		if err != nil {
			t.Fatalf("failed decode: err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}
	}
}

func Test_EncodeDecode_0xFF(t *testing.T) {
	for i := 0; i < 1000; i++ {
		src := make([]byte, i)
		for i := range src {
			src[i] = 0xff
		}
		dst := EncodeToBytes(src)
		got, err := DecodeBytes(dst)
		if err != nil {
			t.Fatalf("failed decode: err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}
	}
}

func Test_EncodeDecode_RandomBytes(t *testing.T) {
	for i := 0; i < 1000000; i++ {
		src := make([]byte, 32+mathrand.Intn(32))
		_, _ = rand.Read(src)
		dst := EncodeToBytes(src)
		got, err := DecodeBytes(dst)
		if err != nil {
			t.Fatalf("failed decode, err = %v", err)
		}
		if !bytes.Equal(src, got) {
			t.Fatalf("failed decode, got = %v, want = %v", got, src)
		}
	}
}

// ----------

func Test_encoder_next(t *testing.T) {
	src := []byte{123, 234, 255}
	enc := newEncoder(src)

	//for _, w := range src {
	//	fmt.Printf("%08b", w)
	//}
	//fmt.Println()

	x, hasMore := enc.next()
	for {
		_ = x
		if !hasMore {
			break
		}
		x, hasMore = enc.next()
	}
}

// ----------

func Test_NewEncoding_panic(t *testing.T) {
	func() {
		encoder := "abcdef"
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("NewEncoding did not panic with encoder %q", encoder)
			}
		}()
		_ = NewEncoding(encoder)
	}()

	func() {
		encoder := []byte(encodeStd)
		encoder[1] = '\n'
		defer func() {
			if r := recover(); r == nil {
				t.Error("NewEncoding did not panic with encoder contains \\n")
			}
		}()
		_ = NewEncoding(string(encoder))
	}()

	func() {
		encoder := []byte(encodeStd)
		encoder[1] = '\r'
		defer func() {
			if r := recover(); r == nil {
				t.Error("NewEncoding did not panic with encoder contains \\r")
			}
		}()
		_ = NewEncoding(string(encoder))
	}()
}

func Test_Decode_CorruptInputError(t *testing.T) {
	src := make([]byte, 256)
	for i := range src {
		src[i] = byte(i)
	}
	_, err := StdEncoding.Decode(make([]byte, 512), src)
	if err == nil || !strings.Contains(err.Error(), "illegal base62 data at input byte") {
		t.Fatal("decoding invalid data did not return CorruptInputError")
	}
}
