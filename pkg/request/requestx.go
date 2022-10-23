package requestx

import (
	"CloudPlatform/internal/handle"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"

	JSON = "json"
	FORM = "form"
)

var (
	ErrMethodNotSupported = errors.New("method is not supported")
	ErrMIMENotSupported   = errors.New("mime is not supported")
)

// make request which contains uploading file
func MakeFileRequest(method, api, fileName, fieldName string, param interface{}) (request *http.Request, err error) {
	method = strings.ToUpper(method)
	if method != POST && method != PUT {
		err = ErrMethodNotSupported
		return
	}

	// create form file
	buf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buf)
	fileWriter, err := bodyWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return
	}

	// read the file
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	// read the file to the fileWriter
	length, err := fileWriter.Write(fileBytes)
	if err != nil {
		return
	}

	bodyWriter.Close()

	// make request
	queryStr := MakeQueryStrFrom(param)
	if queryStr != "" {
		api += "?" + queryStr
	}
	request, err = http.NewRequest(string(method), api, buf)
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	err = request.ParseMultipartForm(int64(length))
	return
}

// make request
func MakeRequest(method, mime, api string, param interface{}) (request *http.Request, err error) {
	method = strings.ToUpper(method)
	mime = strings.ToLower(mime)

	switch mime {
	case JSON:
		var (
			contentBuffer *bytes.Buffer
			jsonBytes     []byte
		)
		jsonBytes, err = json.Marshal(param)
		if err != nil {
			return
		}
		contentBuffer = bytes.NewBuffer(jsonBytes)
		request, err = http.NewRequest(string(method), api, contentBuffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	case FORM:
		queryStr := MakeQueryStrFrom(param)
		var buffer io.Reader

		if (method == DELETE || method == GET) && queryStr != "" {
			api += "?" + queryStr
		} else {
			buffer = bytes.NewReader([]byte(queryStr))
		}

		request, err = http.NewRequest(string(method), api, buffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	default:
		err = ErrMIMENotSupported
		return
	}
	return
}

// make query string from params
func MakeQueryStrFrom(params interface{}) (result string) {
	if params == nil {
		return
	}
	value := reflect.ValueOf(params)

	switch value.Kind() {
	case reflect.Struct:
		var formName string
		for i := 0; i < value.NumField(); i++ {
			if formName = value.Type().Field(i).Tag.Get("form"); formName == "" {
				// don't tag the form name, use camel name
				formName = GetCamelNameFrom(value.Type().Field(i).Name)
			}
			result += "&" + formName + "=" + fmt.Sprintf("%v", value.Field(i).Interface())
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			result += "&" + fmt.Sprintf("%v", key.Interface()) + "=" + fmt.Sprintf("%v", value.MapIndex(key).Interface())
		}
	default:
		return
	}

	if result != "" {
		result = result[1:]
	}
	return
}

// get the Camel name of the original name
func GetCamelNameFrom(name string) string {
	result := ""
	i := 0
	j := 0
	r := []rune(name)
	for m, v := range r {
		// if the char is the capital
		if v >= 'A' && v < 'a' {
			// if the prior is the lower-case || if the prior is the capital and the latter is the lower-case
			if (m != 0 && r[m-1] >= 'a') || ((m != 0 && r[m-1] >= 'A' && r[m-1] < 'a') && (m != len(r)-1 && r[m+1] >= 'a')) {
				i = j
				j = m
				result += name[i:j] + "_"
			}
		}
	}

	result += name[j:]
	return strings.ToLower(result)
}

func Request(r http.Handler, req *http.Request, resp *handle.JsonMsgResult) error {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// extract the response from the response record
	result := w.Result()
	defer result.Body.Close()

	// extract response body
	bodyByte, err := io.ReadAll(result.Body)
	json.Unmarshal(bodyByte, &resp)
	return err
}
