package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

func GenerateSHA256(inputs ...string) string {
	trimmed := make([]string, len(inputs))
	for i, str := range inputs {
		trimmed[i] = strings.TrimSpace(str)
	}
	combined := strings.Join(trimmed, ":")
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func NowStrNano() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
