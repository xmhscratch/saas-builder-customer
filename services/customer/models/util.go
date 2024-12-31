package models

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	// "log"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	matchFirstCap = regexp.MustCompile("([A-Z])([A-Z][a-z])")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// ParseString comment
func ParseString(rawValue interface{}) (string, error) {
	var (
		err      error
		newValue string
	)

	valueKind := reflect.ValueOf(rawValue).Kind()
	switch valueKind {
	case reflect.Float64:
		{
			newValue = string(strconv.FormatFloat(rawValue.(float64), 'E', -1, 64))
			break
		}
	case reflect.Float32:
		{
			newValue = string(strconv.FormatFloat(float64(rawValue.(float32)), 'E', -1, 32))
			break
		}
	case reflect.Int64:
		{
			newValue = string(strconv.FormatInt(rawValue.(int64), 10))
			break
		}
	case reflect.Int32:
		{
			newValue = string(strconv.FormatInt(int64(rawValue.(int32)), 10))
			break
		}
	case reflect.String:
		{
			newValue = string(rawValue.(string))
			break
		}
	case reflect.Bool:
		{
			newValue = strconv.FormatBool(rawValue.(bool))
			break
		}
	case reflect.Slice:
		{
			typeOf := reflect.TypeOf(rawValue).Elem()
			if typeOf.Name() == "uint8" {
				cValue := []byte(rawValue.([]uint8))
				newValue = string(cValue[:])
			}
			break
		}
	default:
		{
			newValue = string("")
			break
		}
	}
	return newValue, err
}

// ParseInt comment
func ParseInt(rawValue interface{}) (int64, error) {
	var (
		err      error
		newValue int64
	)

	valueKind := reflect.ValueOf(rawValue).Kind()
	switch valueKind {
	case reflect.Float64:
		{
			newValue = int64(rawValue.(float64))
			break
		}
	case reflect.Float32:
		{
			newValue = int64(rawValue.(float32))
			break
		}
	case reflect.Int64:
		{
			newValue = int64(rawValue.(int64))
			break
		}
	case reflect.Int32:
		{
			newValue = int64(rawValue.(int32))
			break
		}
	case reflect.String:
		{
			cValue, err := strconv.ParseInt(rawValue.(string), 10, 64)
			if err != nil {
				newValue = int64(math.NaN())
				break
			}
			newValue = int64(cValue)
			break
		}
	case reflect.Bool:
		{
			if bool(rawValue.(bool)) {
				newValue = int64(1)
			} else {
				newValue = int64(0)
			}
			break
		}
	case reflect.Slice:
		{
			typeOf := reflect.TypeOf(rawValue).Elem()
			if typeOf.Name() == "uint8" {
				vStr := string([]byte(rawValue.([]uint8))[:])
				cValue, err := strconv.ParseInt(vStr, 10, 64)
				if err != nil {
					newValue = int64(math.NaN())
					break
				}
				newValue = int64(cValue)
			}
			break
		}
	default:
		{
			newValue = int64(0)
			break
		}
	}

	return newValue, err
}

// ParseFloat comment
func ParseFloat(rawValue interface{}) (float64, error) {
	var (
		err      error
		newValue float64
	)

	valueKind := reflect.ValueOf(rawValue).Kind()
	switch valueKind {
	case reflect.Float64:
		{
			newValue = float64(rawValue.(float64))
			break
		}
	case reflect.Float32:
		{
			newValue = float64(rawValue.(float32))
			break
		}
	case reflect.Int64:
		{
			newValue = float64(rawValue.(int64))
			break
		}
	case reflect.Int32:
		{
			newValue = float64(rawValue.(int32))
			break
		}
	case reflect.String:
		{
			cValue, err := strconv.ParseFloat(rawValue.(string), 64)
			if err != nil {
				newValue = float64(math.NaN())
				break
			}
			newValue = float64(cValue)
			break
		}
	case reflect.Bool:
		{
			if bool(rawValue.(bool)) {
				newValue = float64(1)
			} else {
				newValue = float64(0)
			}
			break
		}
	case reflect.Slice:
		{
			typeOf := reflect.TypeOf(rawValue).Elem()
			if typeOf.Name() == "uint8" {
				vStr := string([]byte(rawValue.([]uint8))[:])
				cValue, err := strconv.ParseFloat(vStr, 64)
				if err != nil {
					newValue = float64(math.NaN())
					break
				}
				newValue = float64(cValue)
			}
			break
		}
	default:
		{
			newValue = float64(0.00)
			break
		}
	}

	return newValue, err
}

// ParseTime comment
func ParseTime(rawValue interface{}) (*time.Time, error) {
	var (
		err      error
		newValue time.Time
	)

	strValue, err := ParseString(rawValue)
	if err != nil {
		return nil, err
	}

	valueKind := reflect.ValueOf(rawValue).Kind()
	switch valueKind {
	case reflect.Float64:
		{
			cValue, err := strconv.ParseFloat(strValue, 64)
			if err != nil {
				return nil, err
			}
			newValue = time.Unix(int64(cValue), 0)
			break
		}
	case reflect.Float32:
		{
			cValue, err := strconv.ParseFloat(strValue, 32)
			if err != nil {
				return nil, err
			}
			newValue = time.Unix(int64(cValue), 0)
			break
		}
	case reflect.Int64:
		{
			cValue, err := strconv.ParseInt(strValue, 0, 64)
			if err != nil {
				return nil, err
			}
			newValue = time.Unix(int64(cValue), 0)
			break
		}
	case reflect.Int32:
		{
			cValue, err := strconv.ParseInt(strValue, 0, 32)
			if err != nil {
				return nil, err
			}
			newValue = time.Unix(int64(cValue), 0)
			break
		}
	case reflect.Bool:
		{
			// newValue = nil
			break
		}
	default:
		{
			newValue, err = time.Parse(time.RFC3339, strValue)
			if err != nil {
				return nil, err
			}
			break
		}
	}

	return &newValue, err
}

// ParseBool comment
func ParseBool(rawValue interface{}) (bool, error) {
	strValue, err := ParseString(rawValue)
	if err != nil {
		return false, err
	}
	strValue = strings.ToLower(strValue)
	return regexp.MatchString("^(true|1|yes|accept|on)", strValue)
}

// Min returns the smaller of a and b.
func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of a and b.
func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// Clamp returns a value restricted between lo and hi.
func Clamp(v, lo, hi int64) int64 {
	return Min(Max(v, lo), hi)
}

// // SimplifyFractions comment
// func SimplifyFractions(numerator float64, denominator float64) (float64, float64) {
// 	for i := math.Floor(math.Max(numerator, denominator)); i > 1; i-- {
// 		if math.Mod(numerator, i) == 0 && math.Mod(denominator, i) == 0 {
// 			numerator /= i
// 			denominator /= i
// 		}
// 	}
// 	return numerator, denominator
// }

// FetchDetailInfo comment
func FetchDetailInfo(db *gorm.DB, obj interface{}) (map[string]interface{}, error) {
	var details map[string]interface{}

	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &details)
	if err != nil {
		return nil, err
	}
	return details, err
}

// StringToSnakeCase comment
func StringToSnakeCase(input string) string {
	output := matchFirstCap.ReplaceAllString(input, "${1}_${2}")
	output = matchAllCap.ReplaceAllString(output, "${1}_${2}")
	output = strings.ReplaceAll(output, "-", "_")

	return strings.ToLower(output)
}

// ConvertJSONToKeySnakeCase comment
func ConvertJSONToKeySnakeCase(j json.RawMessage) json.RawMessage {
	m := make(map[string]json.RawMessage)
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		// Not a JSON object
		return j
	}

	for k, v := range m {
		fixed := StringToSnakeCase(k)
		delete(m, k)
		m[fixed] = ConvertJSONToKeySnakeCase(v)
	}

	b, err := json.Marshal(m)
	if err != nil {
		return j
	}

	return json.RawMessage(b)
}

// ConvertMapToKeySnakeCase comment
func ConvertMapToKeySnakeCase(input map[string]interface{}) (map[string]interface{}, error) {
	var (
		err        error
		infoValues map[string]interface{}
	)

	infoValues = make(map[string]interface{})
	stringInfoValues, err := json.Marshal(input)
	if err != nil {
		return input, err
	}

	rawInfoValues := json.RawMessage(string(stringInfoValues))
	err = json.Unmarshal(ConvertJSONToKeySnakeCase(rawInfoValues), &infoValues)
	if err != nil {
		return input, err
	}
	return infoValues, err
}

func IsValidURL(input string) bool {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return false
	}

	u, err := url.Parse(input)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// ZipSourceDirectory comment
func ZipSourceDirectory(sourcePath string) (*bytes.Buffer, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	writer := zip.NewWriter(buf)
	defer writer.Close()

	// 2. Go through all the files of the source
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(sourcePath), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}
		header.Name = strings.TrimPrefix(strings.Replace(header.Name, strings.TrimPrefix(strings.Replace(sourcePath, filepath.Dir(sourcePath), "", 1), "/"), "", 1), "/")
		if header.Name == "" {
			return nil
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})

	return buf, err
}

// GZipCompress comment
func GZipCompress(data []byte) ([]byte, error) {
	// Create a buffer to write our archive to.
	var (
		err error
		gz  *gzip.Writer
		bf  *bytes.Buffer = new(bytes.Buffer)
	)

	if gz, err = gzip.NewWriterLevel(bf, gzip.BestCompression); err != nil {
		return data, err
	}
	if _, err := gz.Write(data); err != nil {
		return data, err
	}
	if err := gz.Close(); err != nil {
		return data, err
	}

	return bf.Bytes(), err
}

// GZipDecompress comment
func GZipDecompress(data []byte) ([]byte, error) {
	// Create a buffer to write our archive to.
	var (
		err error
		gz  *gzip.Reader
		bf  *bytes.Buffer = new(bytes.Buffer)
	)

	if gz, err = gzip.NewReader(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}
	if _, err := io.Copy(bf, gz); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return bf.Bytes(), err
}
