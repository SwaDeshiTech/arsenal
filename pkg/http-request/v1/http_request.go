package v1

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type RESTStructure struct {
	Method      string            `json:"methodType"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	ContentType string            `json:"contentType"`
	Timeout     time.Duration     `json:"timeout"`
	Body        interface{}       `json:"body"`
}

func (rest *RESTStructure) Do() ([]byte, error) {

	payload, err := createBodyPayload(rest.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(rest.Method, rest.URL, payload)
	if err != nil {
		return nil, err
	}

	addHeadersToRequest(req, rest.Headers)

	client := &http.Client{Timeout: rest.Timeout}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func addHeadersToRequest(req *http.Request, headers map[string]string) {
	for id, value := range headers {
		req.Header.Add(id, value)
	}
}

func createBodyPayload(body interface{}) (io.Reader, error) {
	if body != nil {
		d, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(d), nil
	}
	return nil, nil
}
