package grpcx

import (
	"encoding/pem"
	"strings"
	"testing"
)

func TestSHA1Fingerprint_Empty(t *testing.T) {
	tests := []struct {
		s    SHA1Fingerprint
		want bool
	}{
		{"", true},
		{"abc", false},
	}
	for _, tt := range tests {
		t.Run(string(tt.s), func(t *testing.T) {
			if got := tt.s.Empty(); got != tt.want {
				t.Errorf("SHA1Fingerprint.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1Fingerprint_Match(t *testing.T) {
	block, _ := pem.Decode([]byte(`-----BEGIN CERTIFICATE-----
MIIDFTCCAf0CFDM1zq1rlbnEsf+EQhynC2C8WBOuMA0GCSqGSIb3DQEBCwUAMCcx
CzAJBgNVBAYTAlVTMRgwFgYDVQQDDA9FeGFtcGxlLVJvb3QtQ0EwHhcNMjEwNTI5
MDA0NDI4WhcNMjEwNTI5MDA0NDI4WjBnMQswCQYDVQQGEwJVUzESMBAGA1UECAwJ
WW91clN0YXRlMREwDwYDVQQHDAhZb3VyQ2l0eTEdMBsGA1UECgwURXhhbXBsZS1D
ZXJ0aWZpY2F0ZXMxEjAQBgNVBAMMCWxvY2FsaG9zdDCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBAL7gDonH4kZR8J9cOhYCTIKS05exiKGAJ25GgI32Lzon
K8Mpcu/1XSorbJJaOIleJ0RtSMbSbdH9ao3/Ez6xDNuUguzjIgFUs0CgeYWFl5YE
k2IvCiWvd3nwSHCvZDkShhjFnJjy4bTX2JI18MfhaKfdW15Bp+So9tOIie6Uefix
zTDP6U2Lno4v9uGhQfTYUuZD+3JuUw9R4TXjNd+GDW7LTiFkg63MdDWyoYzTmLx6
emnC0umMgXwClESQ1BDHYqmEfk+Zwy9l1/BBEK0Z2XdjKcp5UBgNE3GiSg8MQCLS
aD48U9zTDeyX4ClE8VCip+7kXrYxL6Jj7iuCoeL10ekCAwEAATANBgkqhkiG9w0B
AQsFAAOCAQEAXC79DsBi/ocVsTcqBpICj2C1DE49Lw9Pdw+6mVGmog5pPY5wCr5b
gsHxl9NS9syagCnlmBqSOIqFGp8z2uI6sQZOfmERxSeyodpA0gDJNHueMhBtJPEw
6nER5h2zbObOUmF5iAdiA+EDZmZgwc/n8R05gNTKd9d94hK8DCoGb1ZLJWj1thB2
gufZN88kU+8db+Kvz3HV+r4kvfy9TxKxlHsEZPnV4Gd6/42D+3Zt8O/A8LjmKAiZ
v3Nm6z6FpuVnqbMCim0NsBu925GSknaXPO6gZm+MFbj7FrmsESpY+0Cybrv/Rl1o
Y/6FWKJ0CfLKg2ykX3tqE2GK3uJ0c30ayg==
-----END CERTIFICATE-----`))

	tests := []struct {
		s    SHA1Fingerprint
		in   []byte
		want bool
	}{
		{"", []byte(""), false},
		{"", []byte("invalid"), false},
		{"7E:12:49:9C:EC:EC:22:DE:53:78:71:79:BF:28:D4:51:2D:66:23:96", block.Bytes, true},
		{"7E12499CECEC22DE53787179BF28D4512D662396", block.Bytes, true},
		{SHA1Fingerprint(strings.ToLower("7E12499CECEC22DE53787179BF28D4512D662396")), block.Bytes, true},
	}
	for _, tt := range tests {
		t.Run(string(tt.s), func(t *testing.T) {
			if got := tt.s.Match(tt.in); got != tt.want {
				t.Errorf("SHA1Fingerprint.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
