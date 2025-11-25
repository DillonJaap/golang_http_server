package request

import (
	"flag"
	uslices "http_server/internal/util/slices"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testCase = flag.String("test-case", "", "the test case to run")
)

func TestRequestLineParse(t *testing.T) {
	flag.Parse()

	type test struct {
		name        string
		request     io.Reader
		wantRequest Request
		wantError   bool
	}

	tests := [][]test{
		// strings
		{
			{
				name:    "Good GET Request line",
				request: strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
				wantRequest: Request{
					RequestLine: RequestLine{
						HttpVersion:   "1.1",
						RequestTarget: "/",
						Method:        "GET",
					},
					Headers: map[string]string{},
					Body:    []byte{},
				},
			},
			{
				name:    "Good GET Request line - line with path",
				request: strings.NewReader("GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
				wantRequest: Request{
					RequestLine: RequestLine{
						HttpVersion:   "1.1",
						RequestTarget: "/coffee",
						Method:        "GET",
					},
					Headers: map[string]string{},
					Body:    []byte{},
				},
			},
			{
				name:    "Good POST Request line - line with path",
				request: strings.NewReader("POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
				wantRequest: Request{
					RequestLine: RequestLine{
						HttpVersion:   "1.1",
						RequestTarget: "/coffee",
						Method:        "POST",
					},
					Headers: map[string]string{},
					Body:    []byte{},
				},
			},
			{
				name:      "Invalid number of parts in request line",
				request:   strings.NewReader("/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
				wantError: true,
			},
			{
				name:      "Invalid method (out of order) Request Line",
				request:   strings.NewReader("/coffee POST HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
				wantError: true,
			},
			{
				name:      "Invalid version in Request line",
				request:   strings.NewReader("POST /coffee HTTP/5.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"),
				wantError: true,
			},
		},
		// chunk
		{
			{
				name: "Good GET Request line - chunk reader",
				request: &chunkReader{
					data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
					numBytesPerRead: 3,
				},
				wantRequest: Request{
					RequestLine: RequestLine{
						HttpVersion:   "1.1",
						RequestTarget: "/",
						Method:        "GET",
					},
					Headers: map[string]string{},
					Body:    []byte{},
				},
			},
			{
				name: "Good GET Request line - line with path - chunk",
				request: &chunkReader{
					data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
					numBytesPerRead: 3,
				},
				wantRequest: Request{
					RequestLine: RequestLine{
						HttpVersion:   "1.1",
						RequestTarget: "/coffee",
						Method:        "GET",
					},
					Headers: map[string]string{},
					Body:    []byte{},
				},
			},
		},
	}

	for _, tt := range uslices.Flatten(tests) {
		t.Run(tt.name, func(t *testing.T) {
			if testCase != nil && *testCase != "" && tt.name != *testCase {
				t.Skip()
			}

			r, err := RequestFromReader(tt.request)
			if !tt.wantError {
				require.NoError(t, err)
				require.NotNil(t, r)
				assert.Equal(t, tt.wantRequest.RequestLine.Method, r.RequestLine.Method)
				assert.Equal(t, tt.wantRequest.RequestLine.RequestTarget, r.RequestLine.RequestTarget)
				assert.Equal(t, tt.wantRequest.RequestLine.HttpVersion, r.RequestLine.HttpVersion)
			} else {
				require.Error(t, err)
			}
		})
	}

}

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

// Read reads up to len(p) or numBytesPerRead bytes from the string per call
// its useful for simulating reading a variable number of bytes per chunk from a network connection
func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := cr.pos + cr.numBytesPerRead
	if endIndex > len(cr.data) {
		endIndex = len(cr.data)
	}
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n

	return n, nil
}
