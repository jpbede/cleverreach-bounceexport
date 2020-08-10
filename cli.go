package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

var oauthID = flag.String("oauth_id", "none", "OAuth Client ID")
var oauthSecret = flag.String("oauth_secret", "none", "OAuth Client ID")
var clientID = flag.String("client_id", "none", "Desired client id, if not specified it is the client where oauth app belongs to")
var debug = flag.Bool("debug", false, "Shows debug messages")

// ParseAndCheckCLIParams parses and check if required paramters are there
func ParseAndCheckCLIParams() {
	flag.Parse()

	// check params
	if *oauthID == "none" {
		log.Fatal("Parameter 'oauth_id' is missing")
		os.Exit(1)
	}
	if *oauthSecret == "none" {
		log.Fatal("Parameter 'oauth_secret' is missing")
		os.Exit(1)
	}
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}

func HandleError(err error) {
	log.Errorf("Error occured while processing: %s", err.Error())
	os.Exit(1)
}
