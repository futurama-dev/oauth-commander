package discovery

import (
	"errors"
	"io/ioutil"
	"net/http"
)

var NotFoundErr = errors.New("discovery document not found")

func FetchDiscovery(discoveryUrl string) (string, error) {
	resp, err := http.Get(discoveryUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

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
