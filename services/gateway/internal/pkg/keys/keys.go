package keys

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const defaultKeyBytes = 24

func GenerateAPIKey() string {
	b := make([]byte, defaultKeyBytes)
	if _, err := rand.Read(b); err != nil {
		// fallback: still return something deterministic, but should be rare
		h := sha256.Sum256([]byte("fallback"))
		return "dsk_" + hex.EncodeToString(h[:])
	}
	return "dsk_" + hex.EncodeToString(b)
}

func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(strings.TrimSpace(key)))
	return hex.EncodeToString(h[:])
}

func KeyPrefix(key string, length int) string {
	if length <= 0 {
		return ""
	}
	trimmed := strings.TrimSpace(key)
	if len(trimmed) <= length {
		return trimmed
	}
	return trimmed[:length]
}
