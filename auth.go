package gokhipu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func setAuth(params map[string]string, method, uri, secret, receicerID string) string {
	var buff bytes.Buffer
	buff.WriteString(url.QueryEscape(method))
	buff.WriteString("&")
	buff.WriteString(url.QueryEscape(uri))
	buff.WriteString(makeParams(params))

	sign := hmac.New(sha256.New, []byte(secret))
	sign.Write(buff.Bytes())
	return receicerID + ":" + hex.EncodeToString(sign.Sum(nil))
}

func makeParams(params map[string]string) string {
	if params == nil {
		return ""
	}

	v := sortKeys(params)
	u := url.Values{}

	for i := range v {
		if params[v[i]] == "" {
			continue
		}
		u.Add(v[i], params[v[i]])
	}

	var urlParams *url.URL
	urlParams, err := url.Parse("")
	if err != nil {
		return ""
	}
	urlParams.RawQuery = u.Encode()
	return "&" + strings.Replace(urlParams.String(), "+", "%20", -1)[1:]
}

func sortKeys(m map[string]string) []string {
	sortedKeys := make([]string, len(m))

	i := 0
	for k := range m {
		sortedKeys[i] = k
		i++
	}

	sort.Strings(sortedKeys)

	return sortedKeys
}
