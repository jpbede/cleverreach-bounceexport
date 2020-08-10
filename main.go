package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"os"
)

func main() {
	// parse and check cli parameters
	ParseAndCheckCLIParams()

	// open file for writing
	fileName := "bounces"
	if *clientID != "none" {
		fileName += "_" + *clientID
	}

	bouncesFile, openErr := os.OpenFile(fileName+".csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if openErr != nil {
		HandleError(openErr)
	}
	defer bouncesFile.Close()

	fileWriter := bufio.NewWriter(bouncesFile)
	writer := csv.NewWriter(fileWriter)
	writer.Comma = ';'
	safeWriter := gocsv.NewSafeCSVWriter(writer)

	// get token of account to which the oauth app belongs to
	ctx := context.Background()
	config := clientcredentials.Config{
		ClientID:     *oauthID,
		ClientSecret: *oauthSecret,
		TokenURL:     "https://rest.cleverreach.com/oauth/token.php",
	}
	client := resty.NewWithClient(config.Client(ctx))

	hasSubAccountToken := false
	// do we want the bounces of a sub-account of us ?
	if *clientID != "none" {
		log.Debugf("Impersonating client %s", *clientID)

		subAccountToken, impersonatingErr := GetTokenForAccount(*clientID, client)
		if impersonatingErr == nil && subAccountToken != "" {
			log.Infof("Got token for client ID %s", *clientID)

			if err := InvalidateToken(client); err != nil {
				log.Errorf("Error invalidating agency token: %s", err.Error())
			} else {
				log.Debug("Successfully invalidated agency token")
			}

			// init new client with sub-account token
			client = resty.New()
			client.SetAuthToken(subAccountToken)
			hasSubAccountToken = true
		} else {
			HandleError(impersonatingErr)
		}
	}

	// now get all bounces
	page := 0
	cnt := 0
	for {
		gotBounces, httpErr := GetBounces(page, client)

		if httpErr != nil {
			HandleError(httpErr)
		}

		gocsv.MarshalCSV(gotBounces, safeWriter)
		cnt += len(gotBounces)

		// if we've 500 bounces there are maybe more, so head over to the next page
		if len(gotBounces) == 500 {
			page++
		} else {
			log.Debug("Page size is below 500, so there are no more bounces.")
			break
		}
	}
	log.Infof("%d bounces written to CSV", cnt)

	// just delete main account tokens
	// sub-account tokens gotten via API are temporarily
	if !hasSubAccountToken {
		if err := InvalidateToken(client); err != nil {
			log.Errorf("Error invalidating token: %s", err.Error())
		} else {
			log.Info("Successfully invalidated token")
		}
	}
}
