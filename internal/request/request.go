package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type States int

const (
	Initialized States = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	CurrState   States
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *Request) Parse(data []byte) (int, error) {

	switch r.CurrState {
	case Initialized:
		n, err := parseRequestLine(string(data), r)

		if err != nil {
			return -1, fmt.Errorf("Error Occured : %s", err)
		}

		if n == 0 {
			return 0, nil
		}

		r.CurrState = Done
		return n, nil
	case Done:
		return -1, errors.New("Error trying to read data from done State")

	default:
		return -1, errors.New("Error: Unkown State")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	req, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	r := &Request{CurrState: Initialized}

	_, err = parseRequestLine(string(req), r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func parseRequestLine(line string, req *Request) (int, error) {

	s := strings.Split(line, "\r\n")

	if len(s) == 1 {
		return 0, nil
	}

	rline := strings.Split(s[0], " ")

	if len(rline) != 3 {
		err := fmt.Errorf("the data contains less value")
		return -1, err
	}

	method := rline[0]
	reqt := rline[1]
	h := strings.Split(rline[2], "/")
	httpv := h[1]

	if method != strings.ToUpper(method) {
		err := errors.New("Not Capitalized String")
		return -1, err
	}

	if httpv != "1.1" {
		err := errors.New("HTTP Version not 1.1")
		return -1, err
	}

	req.RequestLine = RequestLine{Method: method, RequestTarget: reqt, HttpVersion: httpv}

	return len(s[0]), nil
}
