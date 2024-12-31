package core

import (
	"bytes"
	"os"
	"net/url"
	"strings"

	"github.com/segmentio/ksuid"
)

// BuildString comment
func BuildString(parts ...string) string {
	var buf bytes.Buffer
	for _, val := range parts {
		buf.WriteString(val)
	}
	return buf.String()
}

// GetAppDir comment
func GetAppDir() (appDir string, err error) {
	pwd, err := os.Getwd()
	if err != nil {
		return pwd, err
	}
	return pwd, nil
}

// GenerateID comment
func GenerateID() string {
	id := ksuid.New()
	return id.String()
}

// NormalizeRawURLString comment
func NormalizeRawURL(input string) (*url.URL, error) {
	urlRouteURI, err := url.Parse(input)
	if err != nil {
		return urlRouteURI, err
	}
	oldRawQuery := urlRouteURI.RawQuery
	urlRouteURI.RawQuery = ""
	urlRouteURI, err = url.Parse(strings.TrimSuffix(urlRouteURI.String(), "/"))
	if err != nil {
		return urlRouteURI, err
	}
	urlRouteURI.RawQuery = oldRawQuery
	return urlRouteURI, nil
}

// NormalizeRawURLString comment
func NormalizeRawURLString(input string) (string, error) {
	urlRouteURI, err := NormalizeRawURL(input)
	return urlRouteURI.String(), err
}

// FilterMap comment
func FilterMap(vs map[string]interface{}, f func(interface{}, string) bool) map[string]interface{} {
	vsf := make(map[string]interface{})
	for k, v := range vs {
		if f(v, k) {
			vsf[k] = v
		}
	}
	return vsf
}

// IndexInStringSlice returns the first index of the target string t, or -1 if no match is found.
func IndexInStringSlice(t string, vs []string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// IncludeInStringSlice returns true if the target string t is in the slice.
func IncludeInStringSlice(t string, vs []string) bool {
	return IndexInStringSlice(t, vs) >= 0
}
