package response

import (
	"fmt"
	"httpfromtcp/internal/header"
	"io"
	"strconv"
	"strings"
)

type StatusCode int

const (
	OK         StatusCode = 200
	BAD_REQ    StatusCode = 400
	SERVER_ERR StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	ok := []byte("HTTP/1.1 200 OK\n")
	bad_rq := []byte("HTTP/1.1 400 Bad Request\n")
	server_err := []byte("HTTP/1.1 500 Server Error\n")

	var err error
	switch statusCode {
	case OK:
		_, err = w.Write(ok)
	case BAD_REQ:
		_, err = w.Write(bad_rq)
	case SERVER_ERR:
		_, err = w.Write(server_err)

	default:
		err = fmt.Errorf("Not a valid Response Code for now")
	}

	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) header.Headers {
	headers := header.NewHeaders()
	length := strconv.Itoa(contentLen)
	headers["Content-Length"] = length
	headers["Connection"] = "close"
	headers["Content-Type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers header.Headers) error {

	var fullHeaders strings.Builder

	for key, values := range headers {
		header := fmt.Sprintf("%s:%s\n", key, values)
		fullHeaders.WriteString(header)
	}
	fullHeaders.WriteString("\n")

	result := fullHeaders.String()
	_, err := w.Write([]byte(result))

	if err != nil {
		return err
	}
	return nil
}
