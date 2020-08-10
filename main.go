package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"os"
)

func main() {
	// parse and check cli parameters
	ParseAndCheckCLIParams()

	// get token of account to which the oauth app belongs to
	ctx := context.Background()
	config := clientcredentials.Config{
		ClientID:     *oauthID,
		ClientSecret: *oauthSecret,
		TokenURL:     "https://rest.cleverreach.com/oauth/token.php",
	}

	tk, err := config.Token(ctx)
	if err != nil {
		HandleError(err)
	}
	log.Info("Got token for account")
	token := tk.AccessToken

	// do we want the bounces of a sub-account of us ?
	if *clientID != "none" {
		log.Debugf("Impersonating client %s", *clientID)

		client := resty.NewWithClient(config.Client(ctx))

		resp, httperr := client.R().
			SetHeader("Accept", "application/json").
			SetHeader("Content-Type", "application/json").
			Get(fmt.Sprintf("https://rest.cleverreach.com/v3/clients/%s/token", *clientID))
		if httperr != nil {
			HandleError(httperr)
		}
		if resp.StatusCode() == 200 {
			token = string(resp.Body())
			log.Infof("Got token for client ID %s", *clientID)
		} else {
			log.Debugf("Impersonating response body: %s", string(resp.Body()))
			HandleError(errors.New(fmt.Sprintf("Impersonating request returned with a none 200 status code: %d", resp.StatusCode())))
		}
	}

	// now get all bounces
	page := 0
	var bounceList []*Bounce
	for {
		gotBounces, httpErr := GetBounces(page, token)
		bounceList = append(bounceList, gotBounces...)

		if httpErr != nil {
			HandleError(httpErr)
		}

		// if we've 500 bounces there are maybe more, so head over to the next page
		if len(gotBounces) == 500 {
			page++
		} else {
			log.Debug("Page size is below 500, so there are no more bounces. Continuing writing of csv")
			break
		}
	}

	// Then write it to a semicolon separated file because it can contain a comma in bounce message
	fileName := "bounces"
	if *clientID != "none" {
		fileName += "_" + *clientID
	}

	bouncesFile, openErr := os.OpenFile(fileName+".csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if openErr != nil {
		HandleError(openErr)
	}
	defer bouncesFile.Close()

	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = ';'
		return gocsv.NewSafeCSVWriter(writer)
	})

	csvMarshalErr := gocsv.MarshalFile(bounceList, bouncesFile)
	if csvMarshalErr != nil {
		HandleError(csvMarshalErr)
	}
	log.Infof("%d bounces written to CSV", len(bounceList))
}
