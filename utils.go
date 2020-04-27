package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/eoscanada/eos-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func readOffset(offsetPath string) uint64 {
	log.Debug().Msg("reading offset from " + offsetPath)
	data, err := ioutil.ReadFile(offsetPath)
	if err != nil {
		log.Panic().Msg("couldn't read offset from file")
	}
	result, parseError := strconv.Atoi(strings.Trim(string(data), "\n"))
	if parseError != nil {
		log.Panic().Msgf("Failed to parse offset from %+v reason=%+v", offsetPath, parseError)
	}
	return uint64(result)
}

func writeOffset(offsetPath string, offset uint64) {
	log.Debug().Msg("writing offset to " + offsetPath)
	err := ioutil.WriteFile(offsetPath, []byte(strconv.Itoa(int(offset))), 0644)
	if err != nil {
		log.Error().Msgf("couldnt save offeset %+v", err.Error())
	}
}

func readWIF(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	wif := strings.TrimSpace(strings.TrimSuffix(string(content), "\n"))
	return wif
}

func readRsa(filename string) *rsa.PrivateKey {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	block, _ := pem.Decode(data)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	return key
}


func readConfigFile(cfg *Config, path string) {
	_, err  := toml.DecodeFile(path, &cfg)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Panic().Msg(err.Error())
	}
}

func getConfigPath(envVar, defaultValue string) string {
	configPath := flag.String("config", defaultValue, "config file path")
	flag.Parse()
	cfgPath, isSet := os.LookupEnv(envVar)
	if isSet {
		configPath = &cfgPath
	}
	return *configPath
}

func getAddr(port int) string {
	return ":" + strconv.Itoa(port)
}

func rsaSign(digest eos.Checksum256, key *rsa.PrivateKey) (string, error) {
	sign, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}

	// contract require base64 string
	return base64.StdEncoding.EncodeToString(sign), nil
}