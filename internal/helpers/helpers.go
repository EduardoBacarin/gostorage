package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

func Slugify(s string) string {
	str := strings.ToLower(strings.TrimSpace(s))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	str, _, _ = transform.String(t, str)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	str = reg.ReplaceAllString(str, "-")
	str = strings.Trim(str, "-")

	return str
}
