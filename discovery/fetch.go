package discovery

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var NotFoundErr = errors.New("discovery document not found")

func FetchDiscovery(discoveryUrl string) (string, error) {
	resp, err := http.Get(discoveryUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	contentLength := resp.ContentLength

	// TODO in config allow to change size of unreasonable length
	if contentLength > 15000 {
		return "", errors.New("content too long")
	}

	contentType := resp.Header.Get("content-type")
	slices := strings.Split(contentType, ";")

	if slices[0] != "application/json" {
		return "", errors.New("Invalid header type: " + contentType)
	}

	if resp.StatusCode == 404 {
		return "", NotFoundErr
	}

	if resp.StatusCode != 200 {
		return "", errors.New("failed with status: " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
