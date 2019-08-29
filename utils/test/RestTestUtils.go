package test

import bytes2 "bytes"
import "encoding/json"
import "fmt"
import "github.com/gin-gonic/gin"
import "github.com/stretchr/testify/assert"
import "io"
import "net/http"
import "net/http/httptest"
import "net/http/httputil"
import "testing"

type HttpTester struct {
	Router *gin.Engine
}

func (this *HttpTester) T_GET(t *testing.T, code int, path string, body interface{}) *httptest.ResponseRecorder {
	return this.T_(t, "GET", code, path, body)
}

func (this *HttpTester) T_PUT(t *testing.T, code int, path string, body interface{}) *httptest.ResponseRecorder {
	return this.T_(t, "PUT", code, path, body)
}

func (this *HttpTester) T_DELETE(t *testing.T, code int, path string, body interface{}) *httptest.ResponseRecorder {
	return this.T_(t, "DELETE", code, path, body)
}

func (this *HttpTester) T_POST(t *testing.T, code int, path string, body interface{}) *httptest.ResponseRecorder {
	return this.T_(t, "POST", code, path, body)
}

// test reguest
func (this *HttpTester) T_(
	t *testing.T,
	method string,
	code int,
	path string,
	body interface{},
) *httptest.ResponseRecorder {

	var reader io.Reader = nil
	if body != nil {
		bytes, _ := json.Marshal(body)
		reader = bytes2.NewReader(bytes)
	}

	req, _ := http.NewRequest(method, path, reader)

	w := httptest.NewRecorder()
	dump(req)
	this.Router.ServeHTTP(w, req)

	assert.Equal(t, code, w.Code)
	return w
}

func (this *HttpTester) R_BODY(
	t *testing.T,
	method string,
	path string,
	code int,
) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	this.Router.ServeHTTP(w, req)
	assert.Equal(t, code, w.Code)
	return w
}

func dump(r *http.Request) {
	output, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println(string(output))
}
