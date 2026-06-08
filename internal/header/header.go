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
func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	//Seperators and ending CRLF
	SEPERATOR := []byte(":")
	CRLF := []byte("\r\n")

	NO_SEP := errors.New("Not A valid header")

	//  Checking for the index of both occurance
	sep_idx := bytes.Index(data, SEPERATOR)
	crlf_idx := bytes.Index(data, CRLF)

	if crlf_idx == 0 {
		return 0, true, nil
	}

	if sep_idx == -1 {
		return 0, false, NO_SEP
	}

	if crlf_idx == -1 {
		return 0, false, nil
	}

	//Dividing the data into Host name and value
	host := string(data[:sep_idx])
	fval := string(data[sep_idx+1:])

	//Check validitiy if hostname contains spaces it is invalid
	fName, err := validateHeader(host)

	if err != nil {
		return 0, false, fmt.Errorf("Theres an error :%s", err)
	}

	//Remove whitespaces and crlf from the field value
	fval = strings.ReplaceAll(fval, " ", "")
	fval = strings.ReplaceAll(fval, "\r\n", "")

	//Parse the data into the data type
	if h[fName] != "" {
		h[fName] = h[fName] + ", " + fval
	} else {
		h[fName] = fval
	}

	return crlf_idx + len(CRLF), false, nil
}

func validateHeader(h string) (string, error) {

	if strings.Contains(h, " ") {
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
		return strings.ToLower(h), nil
	}

	return "", errors.New("Error Parsing string \n")

}
