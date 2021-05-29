package grpcx

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

type SHA1Fingerprint string

// supported the format with colon and simple string
func normalHex(in string) ([]byte, error) {
	if strings.Contains(in, ":") {
		in = strings.Join(strings.Split(in, ":"), "")
	}
	inBytes, err := hex.DecodeString(in)
	if err != nil {
		return nil, err
	}
	// dst := make([]byte, hex.DecodedLen(len(in)))
	// _, err := hex.Decode(dst, []byte(in))
	// if err != nil {
	// 	return nil, err
	// }
	return inBytes, nil
}

func (s SHA1Fingerprint) Empty() bool {
	return len(s) == 0
}

func (s SHA1Fingerprint) Match(in []byte) bool {
	// validation
	want, err := normalHex(string(s))
	if err != nil {
		return false
	}

	if len(want) != sha1.Size {
		return false
	}

	// prepare
	got := make([]byte, sha1.Size)
	for i, b := range sha1.Sum(in) {
		got[i] = b
	}

	// matching
	return bytes.Equal(want, got)
}
