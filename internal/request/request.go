package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

const bufSize = 8

const (
	initialized = iota
	done
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
	State       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func noError(err error, messages ...string) {
	if err != nil {
		messages = append(messages, err.Error())
		panic(strings.Join(messages, ","))
	}
}
func assertCond(condition bool, messages ...string) {
	if !condition {
		panic(strings.Join(messages, ","))
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, bufSize)
	request := &Request{
		State: initialized,
	}
	readToIndex := 0

	for request.State != done {
		bytesRead, err := reader.Read(buffer[readToIndex:])
		if err == io.EOF {
			request.State = done
		} else if err != nil {
			return request, err
		}
		readToIndex += bytesRead

		bytesConsumed, err := request.parse(buffer)
		if err != nil {
			return request, err
		}
		if bytesConsumed > 0 {
			copy(buffer, make([]byte, len(buffer)))
			readToIndex = 0
		}

		// Shift unparsed data and grow buffer if needed
		if readToIndex == len(buffer) {
			newBuffer := make([]byte, len(buffer)*2)
			copy(newBuffer, buffer)
			buffer = newBuffer
		}
	}
	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case initialized:
		line, bytesRead, err := parseRequestLine(data)
		if bytesRead != 0 {
			r.RequestLine = line
			r.State = done // TODO remove this as we parse more lines
		}
		return bytesRead, err
	}
	return 0, nil
}

func parseRequestLine(data []byte) (RequestLine, int, error) {
	if !strings.Contains(string(data), "\r\n") {
		return RequestLine{}, 0, nil
	}
	requestLine := strings.SplitAfter(string(data), "\r\n")[0]
	regex, err := regexp.Compile(
		`^(?<method>GET|HEAD|POST|PUT|DELETE|CONNECT|OPTIONS|TRACE)` +
			`\s+` +
			`(?<target>\S+)` + // TODO follow full RFC for the target
			`\s+` +
			`HTTP/(?<version>1\.1)\s*\r\n$`,
	)
	noError(err)

	submatches := regex.FindStringSubmatch(requestLine)
	if len(submatches) != 4 {
		return RequestLine{}, 0, fmt.Errorf("malformed request line expected: 'method' 'target' 'version'\nGot: %+v", submatches)
	}
	method := submatches[1]
	target := submatches[2]
	version := submatches[3]

	return RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, len(requestLine), nil
}
