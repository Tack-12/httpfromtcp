package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var bufSize int = 8

// Enum Type
type States int

const (
	Initialized States = iota
	Done
)

// --- REQUEST STRUCT AND METHODS ---
type Request struct {
	RequestLine RequestLine
	CurrState   States
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

//Method to Parse the Data until reaches the last byte, If the
//current state is not Initialized return Error

func (r *Request) Parse(data []byte) (int, error) {

	switch r.CurrState {
	case Initialized:
		n, err := parseRequestLine(string(data), r)

		if err != nil {
			return 0, fmt.Errorf("Error Occured : %s", err)
		}

		if n == 0 {
			return 0, nil
		}
		return n, nil
	case Done:
		return 0, errors.New("Error trying to read data from done State")

	default:
		return 0, errors.New("Error: Unkown State")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	var r *Request = &Request{}

	var readtoIndex int = 0

	r.CurrState = Initialized

	buf := make([]byte, bufSize)

	for {

		n, err := reader.Read(buf[readtoIndex:])

		if err != nil {
			if err == io.EOF {
				r.CurrState = Done
				break
			}
		}

		if cap(buf) == bufSize {
			bufSize *= 2
			oldbuf := buf

			buf = make([]byte, bufSize)
			copy(buf, oldbuf)
		}

		if r.CurrState == Done {
			break
		}

		readtoIndex += n

		n, err = r.Parse(buf)

		if err != nil {
			return nil, err
		}

		nbuff := make([]byte, bufSize-n)
		copy(nbuff, buf[readtoIndex:])

		readtoIndex -= n

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
		return 0, err
	}

	method := rline[0]
	reqt := rline[1]
	h := strings.Split(rline[2], "/")
	httpv := h[1]

	if method != strings.ToUpper(method) {
		err := errors.New("Not Capitalized String")
		return 0, err
	}

	if httpv != "1.1" {
		err := errors.New("HTTP Version not 1.1")
		return 0, err
	}

	req.RequestLine = RequestLine{Method: method, RequestTarget: reqt, HttpVersion: httpv}

	return len(s[0]), nil
}
