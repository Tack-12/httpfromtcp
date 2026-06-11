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
<<<<<<< HEAD
	readUpto := 0
	if len(data[readUpto:]) == 0 {
		return readUpto, false, nil
=======
	read := 0
	done := false

outer:
	for {

		//  Checking for the index of both occurance
		crlf_idx := bytes.Index(data[read:], CRLF)
		fmt.Printf("CD:%s , read:%v CRLF_idx: %v\n", string(data[read:]), read, crlf_idx)

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

		fmt.Printf("Field Name:%s , Field Value:%s \n", fName, fval)

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

	fmt.Printf("The Fields Are: %s ", fields)

	if len(fields) != 2 {
		return "", "", fmt.Errorf("Has a malformed Header line")
>>>>>>> Checking
	}
outer:
	for {
		//  Checking for the index of both occurance
		crlf_idx := bytes.Index(data[readUpto:], CRLF)

<<<<<<< HEAD
		if crlf_idx == 0 {
			readUpto += len(CRLF)
			break outer
		}

		if crlf_idx == -1 {
			return 0, false, nil
		}

		//Get Value upto the first sep
		fN, fV, err := ParseHeader(data[readUpto : readUpto+crlf_idx])

		if err != nil {
			return 0, false, fmt.Errorf("Error Parsing the Header %s", err)
		}

		readUpto += crlf_idx + len(CRLF)

		if h[fN] != "" {
			h[fN] = h[fN] + "," + fV
		} else {
			h[fN] = fV
		}
	}
	return readUpto, false, nil
=======
	fieldN, err := validateHeader(string(fields[0]))
	if err != nil {
		return "", "", fmt.Errorf("Malformed Header Name: %s", err)
	}
	fieldV := strings.TrimSpace(string(fields[1]))

	return fieldN, fieldV, nil
>>>>>>> Checking
}

func ParseHeader(data []byte) (string, string, error) {

<<<<<<< HEAD
	fieldLine := bytes.SplitN(data, []byte(":"), 2)

	if len(fieldLine) != 2 {
		return "", "", fmt.Errorf("Not a valid field Line")
=======
	if strings.TrimSpace(h) != h {
		return "", errors.New("There are spaces in host name")
>>>>>>> Checking
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

<<<<<<< HEAD
	if !valid {
		return "", "", fmt.Errorf("Invalid Field Name , Contains weird things")
=======
	if valid {
		strings.ReplaceAll(h, " ", "")
		return strings.ToLower(h), nil
>>>>>>> Checking
	}

	return strings.ToLower(fName), triFVal, nil
}
