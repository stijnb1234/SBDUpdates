package main

import (
	"crypto/rand"
	"encoding/base32"
	"github.com/dgryski/dgoogauth"
	"net/url"
	"rsc.io/qr"
)

func generate(account string, issuer string) (string, []byte) {
	secret := make([]byte, 64)
	_, err := rand.Read(secret)
	checkErr(err)

	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	URL, err := url.Parse("otpauth://totp")
	checkErr(err)

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)

	params := url.Values{}
	params.Add("secret", secretBase32)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()

	code, err := qr.Encode(URL.String(), qr.Q)
	checkErr(err)
	b := code.PNG()

	return secretBase32, b
}

func validate(secretStr string, token string) bool {
	otpc := &dgoogauth.OTPConfig{
		Secret:      secretStr,
		WindowSize:  3,
		HotpCounter: 0,
	}

	val, err := otpc.Authenticate(token)
	checkErr(err)
	return val
}
