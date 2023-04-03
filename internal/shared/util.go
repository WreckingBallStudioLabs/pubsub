//////
// Shared utils.
//////

package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/thalesfsp/customerror"
)

// GenerateUUID generates a RFC4122 UUID and DCE 1.1: Authentication and
// Security Services.
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateID generates MD5 hash (content-based) based on content.
// Good to be used to avoid duplicated messages.
func GenerateID(ct string) string {
	// Convert the string to bytes
	data := []byte(strings.Trim(ct, "\f\t\r\n "))

	// Create a new SHA-256 hash
	h := sha256.New()

	// Write the bytes to the hash
	h.Write(data)

	// Get the hash sum as a byte slice
	hashSum := h.Sum(nil)

	// Convert the hash sum to a hex string
	hashString := hex.EncodeToString(hashSum)

	return hashString
}

// SliceContains returns true if the slice contains the string.
//
// NOTE: It's case insensitive.
//
// NOTE: @andres moved to here.
func SliceContains(source []string, text string) bool {
	for _, s := range source {
		if strings.EqualFold(s, text) {
			return true
		}
	}

	return false
}

// Unmarshal with custom error.
func Unmarshal(data []byte, v any) error {
	if err := json.Unmarshal(data, &v); err != nil {
		return customerror.NewFailedToError("to unmarshal",
			customerror.WithError(err),
		)
	}

	return nil
}

// Marshal with custom error.
func Marshal(v any) ([]byte, error) {
	data, err := json.Marshal(&v)
	if err != nil {
		return nil, customerror.NewFailedToError("to marshal",
			customerror.WithError(err),
		)
	}

	return data, nil
}

// Decode process stream `r` into `v` and returns an error if any.
func Decode(r io.Reader, v any) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return customerror.NewFailedToError("decode",
			customerror.WithError(err),
		)
	}

	return nil
}

// ReadAll reads all the data from `r` and returns an error if any.
func ReadAll(r io.Reader) ([]byte, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, customerror.NewFailedToError("read response body", customerror.WithError(err))
	}

	return b, nil
}

// HasLetter checks string for an alphabetic letter.
func HasLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}

	return false
}

// LettersToPhoneNumberKeys translates a given string, e.g.: `DOCTOR` to the
// right telephone number keys, e.g.: `362867`.
func LettersToPhoneNumberKeys(value string) string {
	letterMap := map[string][]string{
		"2": {"a", "b", "c"},
		"3": {"d", "e", "f"},
		"4": {"g", "h", "i"},
		"5": {"j", "k", "l"},
		"6": {"m", "n", "o"},
		"7": {"p", "q", "r", "s"},
		"8": {"t", "u", "v"},
		"9": {"w", "x", "y", "z"},
		"*": {"*", "?"},
	}

	translatedString := ""

	for _, n := range value {
		lowerCase := strings.ToLower(string(n))

		for groupNumber, letterGroup := range letterMap {
			for _, letter := range letterGroup {
				if letter != lowerCase {
					continue
				}

				translatedString += groupNumber
			}
		}
	}

	return translatedString
}

// PrintErrorMessages prints the concatenated error messages.
func PrintErrorMessages(errors ...error) string {
	finalErrMsg := ""

	for _, err := range errors {
		finalErrMsg += err.Error() + ". "
	}

	// Trim the last dot.
	finalErrMsg = strings.TrimSuffix(finalErrMsg, ". ")

	return finalErrMsg
}

// TargetName returns the provided target name, or the configured one. A target,
// depending on the storage, is a collection, a table, a bucket, etc.
// For ElasticSearch - as it doesn't have a concept of a database - the target
// is the index.
func TargetName(name, alternative string) (string, error) {
	if name != "" {
		return name, nil
	}

	if alternative != "" {
		return alternative, nil
	}

	return "", customerror.NewMissingError("target name")
}

// MarshalIndent will marshal `v` with `prefix` and `indent` and returns an error
// if any.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	data, err := json.MarshalIndent(&v, prefix, indent)
	if err != nil {
		return nil, customerror.NewFailedToError("to marshal",
			customerror.WithError(err),
		)
	}

	return data, nil
}
