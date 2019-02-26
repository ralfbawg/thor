/******************************************************************************
***** file: error.go
***** description: wrapper of error, for more convenient
***** author: Zhen Liu
***** date: 2016.09.20
***** history: N/A
*****
******************************************************************************/
package werror

import (
	"bytes"
	"common/logging"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// interface WError: short of wrapped error,
//                   It support more features than error built-in golang
type WError struct {
	errCode int
	errMsg  string
}

func (err *WError) Code() int {
	return err.errCode
}

func (err *WError) Error() string {
	return fmt.Sprintf("Code: %d, Msg: %s", err.errCode, err.errMsg)
}

func (err *WError) ToString() []byte {
	return []byte(fmt.Sprintf("Error[%s]", err.Error()))
}

func NewError(code int, format string, args ...interface{}) *WError {
	msg := format
	if len(format) > 0 && len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	err := &WError{
		errCode: code,
		errMsg:  msg,
	}

	return err
}

type HttpResponseData struct {
	Result int         `json:"errCode"`
	Msg    string      `json:"msg"`
	MsgID  string      `json:"msgid"`
	Data   interface{} `json:"data"`
}

type LogRequestInfo struct {
	Method       string
	URL          string
	User_id      string
	X_Auth_Token string
	Form         interface{}
	Body         interface{}
}

func GetResponseInfo(err *WError, data interface{}) []byte {
	response := &HttpResponseData{}
	if err == nil {
		response.Result = 0
		response.Msg = "ok"
	} else {
		response.Result = err.Code()
		response.Msg = string(err.ToString())
	}
	response.Data = data

	ret, ee := json.Marshal(response)
	if ee != nil {
		logging.Error("[GetResponseInfo] Failed, %s, data = %v", string(err.ToString()), data)
		panic("[GetResponseInfo] Failed, " + ee.Error())
	}

	return ret
}

func LogGetReponseInfo(req *http.Request, err *WError, data interface{}) []byte {
	ret := GetResponseInfo(err, data)

	body := ""
	b_body, ee := ioutil.ReadAll(req.Body)
	if ee == nil {
		body = string(b_body)
	}

	logReq := LogRequestInfo{
		Method: req.Method,
		URL:    req.RequestURI,
		Form:   req.Form,
		Body:   body,
	}
	logging.Info("HANDLE_LOG: request = %+v, response = %s", logReq, string(ret))
	return ret
}

func GetRequestBody(req *http.Request, v interface{}) (interface{}, *WError) {
	buff := bytes.NewBufferString("")
	_, err := io.Copy(buff, req.Body)
	if err != nil {
		logging.Error("[GetRequestBody] copy req body error, request=%v, error=%s", req.Body, err.Error())
		return nil, NewError(1003, err.Error())
	}
	text := buff.String()
	req.Body = ioutil.NopCloser(strings.NewReader(text))

	err = json.Unmarshal([]byte(text), v)
	if err != nil {
		logging.Error("[GetRequestBody] decode json body error, text=%s, error=%s", text, err.Error())
		return nil, NewError(1003, err.Error())
	}

	return v, nil
}
