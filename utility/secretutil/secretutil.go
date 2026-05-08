package secretutil

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const (
	secretPrefix      = "HC1"
	secretCharset     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	secretChecksumLen = 6
	secretMinBodyLen  = 12
)

func GenRandomSecret(length int) string {
	minLen := len(secretPrefix) + secretMinBodyLen + secretChecksumLen
	if length < minLen {
		length = minLen
	}
	if length > 256 {
		length = 256
	}
	bodyLen := length - len(secretPrefix) - secretChecksumLen

	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower := "abcdefghijklmnopqrstuvwxyz"
	digit := "0123456789"

	var body string
	for attempt := 0; attempt < 10; attempt++ {
		raw := randomBase62(bodyLen)
		must := string([]byte{
			upper[secureRandInt(len(upper))],
			lower[secureRandInt(len(lower))],
			digit[secureRandInt(len(digit))],
		})
		merged := (must + raw)[:bodyLen]
		arr := []byte(merged)
		for i := len(arr) - 1; i > 0; i-- {
			j := secureRandInt(i + 1)
			arr[i], arr[j] = arr[j], arr[i]
		}
		merged = string(arr)
		if !containsUpper(merged) || !containsLower(merged) || !containsDigit(merged) {
			continue
		}
		body = merged
		break
	}
	if body == "" {
		body = randomBase62(bodyLen)
	}

	prefixAndBody := secretPrefix + body
	checksum := checksumForSecret(prefixAndBody)
	return prefixAndBody + checksum
}

func CheckSecret(secret string) bool {
	s := strings.TrimSpace(secret)
	if s == "" {
		return false
	}
	minLen := len(secretPrefix) + secretMinBodyLen + secretChecksumLen
	if len(s) < minLen || len(s) > 256 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
			return false
		}
	}
	if !strings.HasPrefix(s, secretPrefix) {
		return false
	}
	body := s[len(secretPrefix) : len(s)-secretChecksumLen]
	checksum := s[len(s)-secretChecksumLen:]
	if len(body) < secretMinBodyLen {
		return false
	}
	if !containsUpper(body) || !containsLower(body) || !containsDigit(body) {
		return false
	}
	expected := checksumForSecret(secretPrefix + body)
	return checksum == expected
}

func fnv1a32(s string) uint32 {
	var hash uint32 = 0x811c9dc5
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= 0x01000193
	}
	return hash
}

func toBase62Fixed(num uint32, length int) string {
	base := uint32(len(secretCharset))
	out := make([]byte, length)
	n := num
	for i := length - 1; i >= 0; i-- {
		out[i] = secretCharset[n%base]
		n = n / base
	}
	return string(out)
}

func checksumForSecret(prefixAndBody string) string {
	h := fnv1a32(prefixAndBody)
	return toBase62Fixed(h, secretChecksumLen)
}

func randomBase62(length int) string {
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		out[i] = secretCharset[secureRandInt(len(secretCharset))]
	}
	return string(out)
}

func secureRandInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}

func containsUpper(s string) bool {
	for _, c := range s {
		if c >= 'A' && c <= 'Z' {
			return true
		}
	}
	return false
}

func containsLower(s string) bool {
	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, c := range s {
		if c >= '0' && c <= '9' {
			return true
		}
	}
	return false
}