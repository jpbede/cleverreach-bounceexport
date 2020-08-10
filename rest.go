package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// InvalidateToken invalidates given token at REST API
func InvalidateToken(client *resty.Client) error {
	resp, httperr := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Delete("https://rest.cleverreach.com/v3/oauth/token")

	if httperr != nil {
		return httperr
	}
	if resp.StatusCode() != http.StatusOK {
		return errors.New(fmt.Sprintf("Invalidating returned with none 200: %d %s", resp.StatusCode(), string(resp.Body())))
	}
	return nil
}

func GetTokenForAccount(clientID string, client *resty.Client) (string, error) {
	resp, httperr := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("https://rest.cleverreach.com/v3/clients/%s/token", clientID))
	if httperr != nil {
		HandleError(httperr)
	}
	if resp.StatusCode() == http.StatusOK {
		return strings.Trim(string(resp.Body()), "\""), nil
	}
	log.Debugf("Impersonating response body: %s", string(resp.Body()))
	return "", errors.New(fmt.Sprintf("Impersonating request returned with a none 200 status code: %d", resp.StatusCode()))
}
