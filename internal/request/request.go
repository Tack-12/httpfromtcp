package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// The New Parser State
type ParserState int

const (
	Initialized ParserState = iota
	Done
)

// The Request being Sent
type Request struct {
	RequestLine RequestLine
	State       ParserState
}

// The top line from the req (GET / HTTP/1.1)
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// Get the request from the Reader / network and return the final version of the Request as Struct
func RequestFromReader(reader io.Reader) (*Request, error) {

	var request *Request

	req, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	rl, _, err := parseRequestLine(req)

	request = &Request{}

	request.RequestLine = *rl

	if err != nil {
		return nil, err
	}

	return request, nil
}

// Parse the acutal request and create a new struct
func parseRequestLine(data []byte) (*RequestLine, int, error) {

	//Varibales:
	SEPERATOR := []byte("\r\n")

	//The index at which the Seperator or the Next line byte is found
	idx := bytes.Index(data, SEPERATOR)

	if idx == -1 {
		return nil, 0, nil
	}

	//The Actual ReqLine from the data
	reqLine := data[:idx]

	//Dividing the Req Line into 3 pars: Method - Request Target - HTTP VERSION (seperated by Spaces)
	parts := bytes.Split(reqLine, []byte(" "))

	if len(parts) != 3 {
		err := fmt.Errorf("the data contains less value")
		return nil, 0, err
	}

	//Changing all the slices of  bytes into String for easy comparison & maintianance
	method := string(parts[0])
	reqt := string(parts[1])
	h := strings.Split(string(parts[2]), "/")
	httpv := h[1]

	// Checking if the Method is all caps
	if method != strings.ToUpper(method) {
		err := errors.New("Not Capitalized String")
		return nil, 0, err
	}

	//Checking if the HTTP version is 1.1
	if httpv != "1.1" {
		err := errors.New("HTTP Version not 1.1")
		return nil, 0, err
	}

	//If Everything is fine Create a new Request with the Parsed attributes in the RequestLine
	req := &RequestLine{Method: method, RequestTarget: reqt, HttpVersion: httpv}

	return req, idx, nil
}

// New Parser Method in Request to use the Parser line to do the parsing:
func (rq *Request) Parser(data []byte) (int, error) {

	PARSING_ERROR := errors.New("PARSING ERROR")

	switch rq.State {
	case Initialized:
		reqLine, n, err := parseRequestLine(data)

		if err != nil {
			return 0, PARSING_ERROR
		}

		if n == 0 {
			return 0, nil
		}

		rq.RequestLine = *reqLine
		rq.State = Done

		return n, nil

	case Done:
		return 0, fmt.Errorf("ERROR: Trying to access the Parser in Done state,")

	default:
		return 0, fmt.Errorf("ERROR: Trying to access the Parser in Unkown state,")
	}

}
