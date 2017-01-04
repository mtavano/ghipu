package gokhipu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
)

func setAuth(params map[string]string, method, uri, secret, receicerID string) string {
	toSign := url.QueryEscape(method) + "&" + url.QueryEscape(uri) + setParams(params)
	sign := hmac.New(sha256.New, []byte(secret))
	sign.Write([]byte(toSign))
	return receicerID + ":" + hex.EncodeToString(sign.Sum(nil))
}

func setParams(params map[string]string) string {
	var qs string

	if params == nil {
		return qs
	}

	m := sortKeys(params)

	for k, v := range m {
		qs += "&" + k + "=" + v
	}
	return qs
}

func sortKeys(m map[string]string) map[string]string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	mm := make(map[string]string)
	for _, k := range keys {
		mm[k] = m[k]
	}
	return mm
}
