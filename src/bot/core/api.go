package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (b *Bot) Request(method string, path string, body any, out any, args ...any) error {
	path = fmt.Sprintf("http://%s%s", os.Getenv("API_BIND"), fmt.Sprintf(path, args...))

	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return err
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(data)
		req.Body = io.NopCloser(buf)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", os.Getenv("API_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return err
		}
	}

	return nil
}
