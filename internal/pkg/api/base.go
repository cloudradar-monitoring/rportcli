package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/breathbath/go_utils/utils/env"

	"github.com/sirupsen/logrus"
)

const (
	maxValidResponseCode = 399
	connectionTimeoutSec = 10
)

type Rport struct {
	BaseURL string
	Auth    Auth
}

func New(baseURL string, a Auth) *Rport {
	return &Rport{BaseURL: baseURL, Auth: a}
}

// BaseClient responsible for calling rport API
type BaseClient struct {
	auth Auth
}

func (c *BaseClient) WithAuth(a Auth) {
	c.auth = a
}

func (c *BaseClient) buildClient() *http.Client {
	connectionTimeout := env.ReadEnvInt("CONN_TIMEOUT_SEC", connectionTimeoutSec)
	transport := &http.Transport{
		DisableKeepAlives:     true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		ResponseHeaderTimeout: time.Duration(connectionTimeout) * time.Second,
	}
	cl := &http.Client{Transport: transport}

	return cl
}

func (c *BaseClient) Call(req *http.Request, target interface{}) (resp *http.Response, err error) {
	cl := c.buildClient()
	dump, _ := httputil.DumpRequest(req, true)
	logrus.Debugf("raw request: %s", string(dump))

	if c.auth != nil {
		err = c.auth.AuthRequest(req)
		if err != nil {
			return nil, err
		}
	}

	resp, err = cl.Do(req)
	if err != nil {
		return resp, fmt.Errorf("operation failed with an error: %v", err)
	}
	var respBodyBytes []byte
	if resp.StatusCode > maxValidResponseCode {
		respBodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Warnf("failed to read response body: %v", err)
			e := c.convertResponseCodeToError(resp.StatusCode, nil)
			return resp, e
		}

		var errResp = new(ErrorResp)
		err = json.Unmarshal(respBodyBytes, errResp)
		if err != nil {
			logrus.Warnf("cannot unmarshal error response %s: %v", string(respBodyBytes), err)
		}
		return resp, errResp
	}

	if target == nil {
		return resp, nil
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err == io.EOF {
		return resp, errors.New("no data received from command execution")
	}
	if err != nil {
		return resp, fmt.Errorf("can't parse data from command execution: %v", err)
	}

	logrus.Debugf("Got response: '%s', status code: '%d'", string(respBody), resp.StatusCode)

	err = json.Unmarshal(respBody, target)
	if err != nil {
		return resp, fmt.Errorf("can't parse data from command execution: %v", err)
	}

	return resp, nil
}

func (c *BaseClient) convertResponseCodeToError(respCode int, errResp *ErrorResp) (err error) {
	if respCode == http.StatusNotFound {
		err = errors.New("the specified item doesn't exist")
	} else if respCode == http.StatusInternalServerError {
		err = fmt.Errorf("operation failed %s", errResp.Error())
	} else if respCode == http.StatusBadRequest {
		err = fmt.Errorf("invalid input provided %s", errResp.Error())
	} else {
		err = fmt.Errorf("unknown error %s", errResp.Error())
	}

	return err
}

func closeRespBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	closeErr := resp.Body.Close()
	if closeErr != nil {
		logrus.Error(closeErr)
	}
}
