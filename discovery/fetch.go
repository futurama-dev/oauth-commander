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

	contentType := resp.Header["content-type"][0]
	slices := strings.Split(contentType, ";")

	if slices[0] != "applications/json" {
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
