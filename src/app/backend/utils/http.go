package utils

import (
	"encoding/base64"
	//"errors"
	//"fmt"
	//"github.com/emicklei/go-restful"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func urlEncoding(urlStr string) string {
	return template.URLQueryEscaper(urlStr)
}

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

func HttpRequest(client *http.Client, method, uri, bodyStr string, headers Headers) (int, string, error) {
	LogD("\n\n")
	LogD("HttpRequest: %s %s", method, uri)
	LogD("HttpRequest: Body=%s", bodyStr)
	payload := strings.NewReader(bodyStr)
	req, err := http.NewRequest(method, uri, payload)
	if err != nil {
		return -1, "", err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	//idmTransportUserName := "idmTransportUser"
	//idmTransportUserPassword := "idmTransportUser"
	//req.SetBasicAuth(idmTransportUserName, idmTransportUserPassword)
	res, err := client.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer res.Body.Close()
	//LogD("res.Status: %s", res.Status)
	//if res.StatusCode >= http.StatusBadRequest {
	//	return res.StatusCode, "", errors.New(fmt.Sprintf("HTTP status code: %s", res.Status))
	//}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		LogE("HTTP status code: %s; Read response's body error, %s", res.Status, err.Error())
		return -1, "", err
	}
	jsonStr := string(body)
	LogD("HTTP status: %s; Body: %s", res.Status, jsonStr)
	LogD("\n\n")
	return res.StatusCode, jsonStr, nil
}

//////////////////////////////

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
