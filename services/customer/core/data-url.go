package core

import (
	"errors"
	"strings"
)

// DataURI represents the parsed "data" URL
type DataURI struct {
	Type     string
	Subtype  string
	Params   map[string]string
	IsBase64 bool
	Data     string
}

// IsDataURL return whether the URL is "data" URL
func IsDataURL(url string) bool {
	return strings.HasPrefix(url, "data:")
}

// ParseDataURL parse the "data" URL into components.
func ParseDataURL(url string) (DataURI, error) {
	const (
		dataURIPrefix   = "data:"
		defaultType     = "text"
		defaultSubType  = "plain"
		defaultParam    = "charset=US-ASCII"
		base64Indicator = "base64"
	)

	if !IsDataURL(url) {
		return DataURI{}, errors.New("input URL is not correct Data URL")
	}

	data := url[len(dataURIPrefix):]
	if !strings.Contains(data, ",") {
		return DataURI{}, errors.New("data not found in Data URI")
	}
	// split propeties and actual encoded data
	comp := strings.SplitN(data, ",", 2)
	properties, encodedData := comp[0], comp[1]

	var result DataURI = DataURI{
		Data:   encodedData,
		Params: make(map[string]string),
	}

	for i, prop := range strings.Split(properties, ";") {
		if i == 0 {
			if strings.Contains(prop, "/") {
				appType := strings.SplitN(prop, "/", 2)
				result.Type, result.Subtype = appType[0], appType[1]
			} else {
				params := strings.Split(defaultParam, "=")
				result.Type, result.Subtype = defaultType, defaultSubType
				result.Params[params[0]] = params[1]
			}
		} else {
			if prop == base64Indicator {
				result.IsBase64 = true
			} else {
				// ignore if not valid properties assignment
				if strings.Contains(prop, "=") {
					propComponets := strings.SplitN(prop, "=", 2)
					result.Params[propComponets[0]] = propComponets[1]
				}
			}
		}
	}

	return result, nil
}
