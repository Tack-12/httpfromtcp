package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/header"
	"io"
	"strings"
)

var BUFFSIZE int = 8

// The New Parser State
type ParserState int

const (
	Initialized ParserState = iota
	RequestParsingHeader
	Done
)

// The Request being Sent
type Request struct {
	RequestLine RequestLine
	State       ParserState
	Headers     header.Headers
}

// The top line from the req (GET / HTTP/1.1)
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// New Parser Method in Request to use the Parser line to do the parsing:
func (rq *Request) Parser(data []byte) (int, error) {

	readuntil := 0

outer:
	for {
		switch rq.State {
		case Initialized:
			reqLine, n, err := parseRequestLine(data[readuntil:])

			if err != nil {
				return readuntil, fmt.Errorf("Error occured %s", err)
			}

			if n == 0 {
				break outer
			}

			rq.RequestLine = *reqLine
			rq.State = Done

			readuntil += n

		case RequestParsingHeader:
			n, done, err := rq.Headers.Parse(data)

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			if done {
				rq.State = Done
				readuntil += n
			}

			readuntil += n
		case Done:
			break outer
		default:
			return 0, fmt.Errorf("ERROR: Trying to access the Parser in Unkown state,")
		}
	}

	return readuntil, nil
}

// Get the request from the Reader / network and return the final version of the Request as Struct
func RequestFromReader(reader io.Reader) (*Request, error) {

	READUNTIL := 0

	buf := make([]byte, BUFFSIZE)

	var request *Request

	request = &Request{}
	request.State = Initialized
	request.Headers = header.NewHeaders()

	for request.State != Done {

		//Instead of checking for free-ness check for how much of the data in buff is already read
		//If all data in buff is read than double if not keep it stagnant.
		if len(buf) == READUNTIL {
			BUFFSIZE = BUFFSIZE * 2
			newbuf := make([]byte, BUFFSIZE)
			copy(newbuf, buf)
			buf = newbuf
		}

		nRead, err := reader.Read(buf[READUNTIL:])

		if err != nil {
			if err == io.EOF {
				request.State = Done
			}
			return nil, fmt.Errorf("error occured : %s", err)
		}

		READUNTIL += nRead
		nParsed, err := request.Parser(buf[:READUNTIL])

		if err != nil {
			return nil, fmt.Errorf("There was an Error Parsing %s", err)
		}

		copy(buf, buf[nParsed:READUNTIL])
		READUNTIL -= nParsed

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

	read := idx + len(SEPERATOR)

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

	return req, read, nil
}
