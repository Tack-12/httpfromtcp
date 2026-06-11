package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/header"
	"io"
	"strconv"
	"strings"
)

var BUFFSIZE int = 8

// The New Parser State
type ParserState int

const (
	Initialized ParserState = iota
	RequestParsingHeader
	ParsingBody
	Done
)

// The Request being Sent
type Request struct {
	RequestLine RequestLine
	State       ParserState
	Headers     header.Headers
	Body        string
}

// The top line from the req (GET / HTTP/1.1)
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func CreateRequest() *Request {
	return &Request{
		State:   Initialized,
		Headers: header.NewHeaders(),
		Body:    "",
	}
}

func (rq *Request) hasBody() bool {
	value, exists := rq.Headers.Get("Content-Length")

	if !exists {
		value = "0"
	}
	num, err := strconv.Atoi(value)

	if err != nil {
		num = 0
	}

	return num > 0

}

// New Parser Method in Request to use the Parser line to do the parsing:
func (rq *Request) Parser(data []byte) (int, error) {

	read := 0

outer:
	for {
		currentData := data[read:]

		if len(currentData) == 0 {
			break outer
		}

		switch rq.State {
		case Initialized:
			reqLine, n, err := parseRequestLine(currentData)

			if err != nil {
				return read, fmt.Errorf("Error occured %s", err)
			}

			if n == 0 {
				break outer
			}

			rq.RequestLine = *reqLine
			rq.State = RequestParsingHeader

			read += n

		case RequestParsingHeader:
			n, done, err := rq.Headers.Parse(currentData)

			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			if done {

				if rq.hasBody() {
					rq.State = ParsingBody
				} else {
					rq.State = Done
				}

			}

			read += n

		case ParsingBody:

			value, exists := rq.Headers.Get("Content-Length")

			if !exists {
				rq.State = Done
				break outer
			}

			if value == "0" {
				rq.State = Done
				break outer
			}

			num, err := strconv.Atoi(value)

			if err != nil {
				return 0, fmt.Errorf("Error converting the value")
			}

			remainder := min(num-len(rq.Body), len(currentData))
			rq.Body += string(currentData[:remainder])
			read += remainder

			if len(rq.Body) == num {
				rq.State = Done
			}

		case Done:
			break outer
		default:
			return 0, fmt.Errorf("ERROR: Trying to access the Parser in Unkown state,")
		}
	}

	return read, nil
}

// Get the request from the Reader / network and return the final version of the Request as Struct
func RequestFromReader(reader io.Reader) (*Request, error) {

	read := 0

	buf := make([]byte, BUFFSIZE)

	var request *Request

	request = CreateRequest()

	for request.State != Done {

		//Instead of checking for free-ness check for how much of the data in buff is already read
		//If all data in buff is read than double if not keep it stagnant.
		if len(buf) == read {
			BUFFSIZE = BUFFSIZE * 2
			newbuf := make([]byte, BUFFSIZE)
			copy(newbuf, buf)
			buf = newbuf
		}

		nRead, err := reader.Read(buf[read:])

		if err != nil {
			if err == io.EOF {
				request.State = Done
			}
			return nil, fmt.Errorf("error occured : %s", err)
		}

		read += nRead
		nParsed, err := request.Parser(buf[:read])

		if err != nil {
			return nil, fmt.Errorf("There was an Error Parsing %s", err)
		}

		copy(buf, buf[nParsed:read])
		read -= nParsed

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
