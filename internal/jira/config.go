package jira

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Config struct {
	URL          string         `toml:"url"`
	AccessToken  string         `toml:"access_token"`
	AccessSecret string         `toml:"access_secret"`
	ConsumerKey  string         `toml:"consumer_key"`
	KeyFile      string         `toml:"key_file"`
	Projects     map[string]int `toml:"projects"`
	privKey      *rsa.PrivateKey
}

func (cfg *Config) ensureKeyLoaded() {
	if cfg.privKey != nil {
		return
	}

	fileBytes, err := os.ReadFile(cfg.KeyFile)
	maybeDie(err != nil, "could not read %s: %s", cfg.KeyFile, err)

	keyDERBlock, _ := pem.Decode(bytes.TrimSpace(fileBytes))
	maybeDie(keyDERBlock == nil, "unable to decode key PEM block")

	privateKey, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	maybeDie(err != nil, "unable to parse privkey: %v", err)

	cfg.privKey = privateKey
}

func maybeDie(predicate bool, msg string, args ...any) {
	if !predicate {
		return
	}

	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
