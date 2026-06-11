package header

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Header Type header name : value
type Headers map[string]string

// Function to create a header
func NewHeaders() Headers {
	return make(Headers)
}

// Method to Parse the header into the datatype
func (h Headers) Parse(data []byte) (int, bool, error) {

	//Seperators and ending CRLF
	CRLF := []byte("\r\n")
	read := 0
	done := false

outer:
	for {

		//  Checking for the index of both occurance
		crlf_idx := bytes.Index(data[read:], CRLF)

		if crlf_idx == 0 {
			read += len(CRLF)
			done = true
			break outer
		}

		if crlf_idx == -1 {
			break outer
		}

		//Check validitiy if hostname contains spaces it is invalid
		fName, fval, err := GetFields(data[read : read+crlf_idx])

		if err != nil {
			return 0, false, fmt.Errorf("Theres an error :%s", err)
		}

		//Updating the read to parse until the current crlf
		read += crlf_idx + len(CRLF)

		//Parse the data into the data type
		if h[fName] != "" {
			temp := h[fName]
			h[fName] = ""
			h[fName] = temp + ", " + fval
		} else {
			h[fName] = fval
		}
	}
	return read, done, nil
}

func GetFields(data []byte) (string, string, error) {
	fields := bytes.SplitN(data, []byte(":"), 2)

	if len(fields) != 2 {
		return "", "", fmt.Errorf("Has a malformed Header line")
	}

	fieldN, err := validateHeader(string(fields[0]))
	if err != nil {
		return "", "", fmt.Errorf("Malformed Header Name: %s", err)
	}
	fieldV := strings.TrimSpace(string(fields[1]))

	return fieldN, fieldV, nil
}

func validateHeader(h string) (string, error) {

	if strings.TrimSpace(h) != h {
		return "", errors.New("There are spaces in host name")
	}

	if len(h) <= 1 {
		return "", errors.New("Not enough length..")
	}

	valid, err := regexp.MatchString(`^[A-Za-z0-9!#$%&'*+\-._\|~^]+$`, h)

	if err != nil {
		return "", fmt.Errorf("Error Parsing string: %s \n", err)
	}

	if valid {
		strings.ReplaceAll(h, " ", "")
		return strings.ToLower(h), nil
	}

	return "", errors.New("Error Parsing string \n")

}

func (h Headers) Get(key string) (string, bool) {

	key = strings.ToLower(key)

	value, ok := h[key]

	return value, ok
}
