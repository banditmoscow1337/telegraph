package telegraph

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/goccy/go-json"
)

func (t *Telegraph) call(uri string, querys []Query, body []byte) (ret []byte, err error) {
	var (
		header = http.Header{}
		resp   *http.Response
		data   []byte
		r      Result
	)

	method := "GET"

	if len(body) > 0 {
		header.Add("Content-Type", "application/json")
		method = "POST"
	} else {
		uri += "?access_token=" + t.a.AccessToken + "&"
		for _, q := range querys {
			uri += q.Name + "=" + url.QueryEscape(q.Value) + "&"
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, "https://api.telegra.ph/"+uri, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		return
	}

	if data, err = io.ReadAll(resp.Body); err != nil {
		return
	}

	err = json.Unmarshal(data, &r)
	if err != nil {
		return
	}

	if !r.Ok {
		err = errors.New(r.Error)
	}

	ret = r.Result

	return
}
