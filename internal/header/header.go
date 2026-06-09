package header

import (
	"bytes"
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
	readUpto := 0

	for {
		//  Checking for the index of both occurance
		crlf_idx := bytes.Index(data[readUpto:], CRLF)

		if crlf_idx == 0 {
			readUpto += len(CRLF)
			break
		}

		if crlf_idx == -1 {
			break
		}

		//Get Value upto the first sep
		fN, fV, err := ParseHeader(data[readUpto : readUpto+crlf_idx])

		if err != nil {
			return 0, false, fmt.Errorf("Error Parsing the Header %s", err)
		}

		fmt.Printf("FN: %s, Fv:%s \n", fN, fV)

		readUpto += crlf_idx + len(CRLF)

		if h[fN] != "" {
			h[fN] = h[fN] + "," + fV
		} else {
			h[fN] = fV
		}
	}
	return readUpto, true, nil
}

func ParseHeader(data []byte) (string, string, error) {

	fieldLine := bytes.SplitN(data, []byte(":"), 2)

	if len(fieldLine) != 2 {
		return "", "", fmt.Errorf("Not a valid field Line")
	}

	fName := string(fieldLine[0])
	fVal := fieldLine[1]

	triFVal := string(bytes.TrimSpace(fVal))

	if fName != strings.TrimSpace(fName) {
		return "", "", fmt.Errorf("Invalid Field Name , Contains Spaces")
	}

	valid, err := regexp.MatchString("^[A-Za-z0-9!#$%&'*+.^_`|~-]+$", fName)

	if err != nil {

		return "", "", fmt.Errorf("Error: %s", err)
	}

	if !valid {
		return "", "", fmt.Errorf("Invalid Field Name , Contains weird things")
	}

	return strings.ToLower(fName), triFVal, nil
}
