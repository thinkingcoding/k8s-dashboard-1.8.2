package utils

import (
	"encoding/base64"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func SetCookie(w http.ResponseWriter, k, v string) {
	cookie := http.Cookie{
		Name:     k,
		Value:    v,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   1800}
	http.SetCookie(w, &cookie)
}

func GetCookie(r *http.Request, k string) string {
	cookie, err := r.Cookie(k)
	if err == nil {
		return cookie.Value
	} else {
		LogW("Read cookie[%s] error: %v", k, err)
		return ""
	}
}

func DeleteCookie(w http.ResponseWriter, k string) {
	cookie := http.Cookie{Name: k, Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)
}

func Curl(client *http.Client, method, uri, bodyStr string, headers Headers) (int, string, error) {
	LogD("HttpRequest: %s %s", method, uri)
	LogD("HttpRequest: body=%s", bodyStr)
	payload := strings.NewReader(bodyStr)
	req, err := http.NewRequest(method, uri, payload)
	if err != nil {
		return -1, "", err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	res, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LogE("HTTP status code: %s; Read response's body error, %s", res.Status, err.Error())
		return -1, "", err
	}
	jsonStr := string(body)
	LogD("HttpRequest: status=%s\n", res.Status)
	return res.StatusCode, jsonStr, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func urlEncoding(urlStr string) string {
	return template.URLQueryEscaper(urlStr)
}
