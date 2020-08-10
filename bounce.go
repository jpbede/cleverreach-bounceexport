package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// Bounce represents a CleverReach bounce
type Bounce struct {
	EMail         string `json:"email" csv:"email"`
	Category      string `json:"category" csv:"category"`
	Occurrences   int    `json:"occurences" csv:"occurences"`
	LastUpdate    int    `json:"last_update" csv:"last_update"`
	LastUpdateGMT string `json:"last_update_gmt" csv:"last_update_gmt"`
	ExpiresBy     int    `json:"expires_by" csv:"expires_by"`
	ExpiresByGMT  string `json:"expires_by_gmt" csv:"expires_by_gmt"`
	BounceMessage string `json:"bounce_message" csv:"bounce_message"`
	Type          string `json:"type" csv:"type"`
	TypeID        string `json:"type_id" csv:"type_id"`
}

func GetBounces(page int, client *resty.Client) ([]*Bounce, error) {
	resp, httperr := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		Get(fmt.Sprintf("https://rest.cleverreach.com/v3/bounces?page=%d&pagesize=500", page))
	if httperr != nil {
		return nil, httperr
	}
	if resp.StatusCode() != 200 {
		log.Debugf("Bounce request returned with body %s", string(resp.Body()))
	}
	log.Debugf("Bounce request returned with HTTP %d", resp.StatusCode())

	var bounces []*Bounce
	marshalErr := json.Unmarshal(resp.Body(), &bounces)
	if marshalErr != nil {
		return nil, marshalErr
	}

	log.Infof("Got %d bounces for page %d", len(bounces), page)

	return bounces, nil
}
