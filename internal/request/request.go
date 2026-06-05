package request

import (
	"bytes"
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
		n, err := parseRequestLine(data, r)

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

		readN, err := r.Parse(buf[:readtoIndex])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufSize])

		readtoIndex -= readN
	}
	return r, nil

}

func parseRequestLine(line []byte, req *Request) (int, error) {
	SEPERATOR := []byte("\r\n")
	idx := bytes.Index(line, SEPERATOR)

	if idx == -1 {
		return 0, nil
	}

	startline := line[:idx]
	readL := idx + len(SEPERATOR)

	rline := bytes.Split(startline, []byte(" "))

	if len(rline) != 3 {
		return 0, fmt.Errorf("Not Enough Data")

	}

	method := string(rline[0])
	reqt := string(rline[1])
	h := strings.Split(string(rline[2]), "/")
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

	return readL, nil
}
