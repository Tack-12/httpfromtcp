package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	req, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	rl, err := parseRequestLine(string(req))

	if err != nil {
		return nil, err
	}

	return rl, nil
}

func parseRequestLine(line string) (*Request, error) {

	s := strings.Split(line, "\r\n")

	rline := strings.Split(s[0], " ")

	if len(rline) != 3 {
		err := fmt.Errorf("the data contains less value")
		return nil, err
	}

	method := rline[0]
	reqt := rline[1]
	h := strings.Split(rline[2], "/")
	httpv := h[1]

	if method != strings.ToUpper(method) {
		err := errors.New("Not Capitalized String")
		return nil, err
	}

	if httpv != "1.1" {
		err := errors.New("HTTP Version not 1.1")
		return nil, err
	}

	req := &Request{RequestLine{Method: method, RequestTarget: reqt, HttpVersion: httpv}}

	return req, nil
}
