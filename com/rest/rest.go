package rest

import (
	"bj-pfd2/com/log"
	"io"
	"io/ioutil"
	"net/http"
)

func Client(url string, method string, body io.Reader, header http.Header) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Close body [%v] error: %v", url, err)
		}
	}(resp.Body)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Response Info:")
	//fmt.Println("  Error      :", err)
	//fmt.Println("  Status Code:", resp.StatusCode)
	//fmt.Println("  Status     :", resp.Status)
	//fmt.Println("  Proto      :", resp.Proto)
	//fmt.Println("  Body       :\n", string(respBytes))
	return respBytes, nil
}
