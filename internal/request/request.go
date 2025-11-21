package request

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
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
	parseRequestLine := func(content string) (RequestLine, error) {
		lines := strings.Split(string(content), "\r\n")
		if len(lines) == 0 {
			return RequestLine{}, fmt.Errorf("no content read from io.Reader")
		}

		requestLine := lines[0]
		regex, err := regexp.Compile(
			`^(?<method>GET|HEAD|POST|PUT|DELETE|CONNECT|OPTIONS|TRACE)` +
				`\s+` +
				`(?<target>\S+)` + // TODO follow full RFC for the target
				`\s+` +
				`HTTP/(?<version>1\.1)$`,
		)
		noError(err)

		submatches := regex.FindStringSubmatch(requestLine)
		if len(submatches) != 4 {
			return RequestLine{}, fmt.Errorf("malformed request line expected: 'method' 'target' 'version'\nGot: %+v", submatches)
		}
		method := submatches[1]
		target := submatches[2]
		version := submatches[3]

		return RequestLine{
			HttpVersion:   version,
			RequestTarget: target,
			Method:        method,
		}, nil

	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return &Request{}, err
	}

	reqLine, err := parseRequestLine(string(content))
	return &Request{
		RequestLine: reqLine,
		Headers:     map[string]string{},
		Body:        []byte{},
	}, err
}
