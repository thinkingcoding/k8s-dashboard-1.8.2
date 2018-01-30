package utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Headers map[string]string

type RequestTokenReq struct {
	ReturnUri string `json:"returnUri"`
	Tenant    string `json:"tenant"`
}

type RequestTokenRes struct {
	IssuedAt       string `json:"issuedAt"`
	IssuedAtAsDate int    `json:"issuedAtAsDate"`
	Tenant         string `json:"tenant"`
	Expires        string `json:"expires"`
	Id             string `json:"id"`
}

type AccessTokenReq struct {
	RequestToken string `json:"requestToken"`
	ReturnUri    string `json:"returnUri"`
	Tenant       string `json:"tenant"`
}

type AccessTokenRes struct {
	SessionState string `json:"sessionState"`
	Token        struct {
		Expires  string `json:"expires"`
		ID       string `json:"id"`
		IssuedAt string `json:"issuedAt"`
		Tenant   struct {
			Enabled bool   `json:"enabled"`
			ID      string `json:"id"`
			Name    string `json:"name"`
		} `json:"tenant"`
	} `json:"token"`
	User struct {
		Groups []struct {
			DisplayName string `json:"displayName"`
			GroupInfo   string `json:"groupInfo"`
			ID          string `json:"id"`
			Metadata    struct {
				Description string `json:"description"`
			} `json:"metadata"`
			Name  string        `json:"name"`
			Roles []interface{} `json:"roles"`
		} `json:"groups"`
		ID      string `json:"id"`
		Name    string `json:"name"`
		Profile struct {
			CommonName  string `json:"common_name"`
			DisplayName string `json:"displayName"`
			Email       string `json:"email"`
			UserEmail   string `json:"userEmail"`
			Username    string `json:"username"`
		} `json:"profile"`
		Roles []struct {
			Application string `json:"application"`
			Description string `json:"description"`
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			Integration string `json:"integration"`
			Locked      bool   `json:"locked"`
			Name        string `json:"name"`
			Type        string `json:"type"`
		} `json:"roles"`
		Username string `json:"username"`
	} `json:"user"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRes struct {
	RefreshToken string `json:"refreshToken"`
	Token        struct {
		Expires  string `json:"expires"`
		ID       string `json:"id"`
		IssuedAt string `json:"issuedAt"`
		Tenant   struct {
			Description string `json:"description"`
			Enabled     bool   `json:"enabled"`
			ID          string `json:"id"`
			Name        string `json:"name"`
		} `json:"tenant"`
	} `json:"token"`
	User struct {
		Groups []struct {
			DisplayName string        `json:"displayName"`
			ID          string        `json:"id"`
			Metadata    struct{}      `json:"metadata"`
			Name        string        `json:"name"`
			Roles       []interface{} `json:"roles"`
		} `json:"groups"`
		ID      string `json:"id"`
		Name    string `json:"name"`
		Profile struct {
			Username string `json:"username"`
		} `json:"profile"`
		Roles []struct {
			Application string `json:"application"`
			Description string `json:"description"`
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			Locked      bool   `json:"locked"`
			Name        string `json:"name"`
			Type        string `json:"type"`
		} `json:"roles"`
		Username string `json:"username"`
	} `json:"user"`
}

var (
	CDF_DEBUG           string = getEnvOrDefault("JUST_CDF_DEBUG_AND_USER_DONOT_SET", "")
	CLIENT_REDIRECT_URI string = getEnvOrDefault("CLIENT_REDIRECT_URI", "https://localhost:9099/loading.html")
	TOKEN_COOKIE_NAME   string = getEnvOrDefault("TOKEN_COOKIE_NAME", "X-CDF-K8S-TOKEN")
	CDF_API_SERVER      string = getEnvOrDefault("CDF_API_SERVER", "https://shclitvm0682.hpeswlab.net:5443")
	IDM_API_SERVER      string = getEnvOrDefault("IDM_API_SERVER", "https://shclitvm0682.hpeswlab.net:5443")
)

var Client = &http.Client{
	// TODO add a timeout, either the goroutines will deadlock if the
	// connection takes too long to resolve.
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func CheckRedirectPage(w http.ResponseWriter, r *http.Request) bool {
	if CDF_DEBUG == "CDF_DEBUG" {
		return false
	}
	if r.URL.EscapedPath() == "/" || r.URL.EscapedPath() == "/index.html" {
		token := GetTokenFromCookie(r)
		if token == "" {
			redirectLoginPage(w, r)
			return true
		} else {
			if _, err := ValidateAccessToken(token); err != nil {
				redirectLoginPage(w, r)
				return true
			}
		}
	} else if r.URL.EscapedPath() == "/loading.html" {
		requestToken := r.FormValue("token")
		if requestToken != "" {
			token, err := AccessToken(requestToken)
			if err != nil {
				LogE("GetAccessTokenFromIDM Error: %v", err)
			} else {
				setTokenFromCookie(w, token)
				redirectIndexPage(w, r)
				return true
			}
		} else {
			//token := GetTokenFromCookie(r)
			//if token == "" {
			//	redirectLoginPage(w, r)
			//	return true
			//} else {
			//	if _, err := ValidateAccessToken(token); err != nil {
			//		redirectLoginPage(w, r)
			//		return true
			//	}
			//	redirectIndexPage(w, r)
			//	return true
			//}
		}
	} else if r.URL.EscapedPath() == "/logout" {
		deleteTokenFromCookie(w)
		RedirectLogoutPage(w, r)
		return true
	}
	return false
}

func CheckApi(w http.ResponseWriter, r *http.Request) error {
	if CDF_DEBUG == "CDF_DEBUG" {
		return nil
	}
	token := GetTokenFromCookie(r)
	_, err := ValidateAccessToken(token)
	if err != nil {
		http.Redirect(w, r, "/loading.html", http.StatusMovedPermanently)
		return err
	}
	return nil
}

func redirectLoginPage(w http.ResponseWriter, r *http.Request) {
	deleteTokenFromCookie(w)
	requestToken, err := requestToken()
	if err != nil {
		LogE(err.Error())
		return
	}
	loginUri := IDM_API_SERVER + "/idm-service/idm/v0/login?tenant=Provider&token=" + urlEncoding(requestToken)
	LogD("Login URI: %s", loginUri)
	http.Redirect(w, r, loginUri, http.StatusMovedPermanently)
}

func redirectIndexPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/#!/login", http.StatusMovedPermanently)
}

func RedirectLogoutPage(w http.ResponseWriter, r *http.Request) {
	deleteTokenFromCookie(w)
	requestToken, err := requestToken()
	if err != nil {
		LogE(err.Error())
		return
	}
	logoutUri := IDM_API_SERVER + "/idm-service/idm/v0/logout?tenant=Provider&token=" + urlEncoding(requestToken)
	LogD("Logout URI: %s", logoutUri)
	http.Redirect(w, r, logoutUri, http.StatusMovedPermanently)
}

func AccessToken(requestToken string) (string, error) {
	body := &AccessTokenReq{
		RequestToken: requestToken,
		ReturnUri:    "",
		Tenant:       "",
	}
	bodyByte, _ := json.Marshal(body)
	code, bodyStr, err := Curl(Client,
		"POST",
		CDF_API_SERVER+"/suiteInstaller/urest/v1.1/tokens/access",
		string(bodyByte), Headers{
			"Accept":       "application/json",
			"content-type": "application/json",
		})
	if err != nil {
		return "", err
	}
	if code == 201 {
		var accessToken AccessTokenRes
		if err := json.Unmarshal([]byte(bodyStr), &accessToken); err == nil {
			return accessToken.Token.ID, nil
		} else {
			return "", err
		}
	} else if code == 401 {
		return "", errors.New("Unauthorized")
	} else if code == 403 {
		return "", errors.New("Forbidden")
	} else if code == 404 {
		return "", errors.New("Not Found")
	} else {
		return "", errors.New(bodyStr)
	}
}

func ValidateAccessToken(token string) (string, error) {
	code, bodyStr, err := Curl(Client,
		"GET",
		CDF_API_SERVER+"/suiteInstaller/urest/v1.1/tokens/"+token,
		"", Headers{
			"Accept":       "application/json",
			"content-type": "application/json",
		})
	if err != nil {
		return "", err
	}
	if code == 200 {
		return "", nil
	} else if code == 401 {
		return "", errors.New("Unauthorized")
	} else if code == 403 {
		return "", errors.New("Forbidden")
	} else if code == 404 {
		return "", errors.New("Not Found")
	} else {
		return "", errors.New(bodyStr)
	}
}

func requestToken() (string, error) {
	body := &RequestTokenReq{
		ReturnUri: CLIENT_REDIRECT_URI,
		Tenant:    "Provider",
	}
	bodyByte, _ := json.Marshal(body)
	code, bodyStr, err := Curl(Client,
		"POST",
		CDF_API_SERVER+"/suiteInstaller/urest/v1.1/tokens/request",
		string(bodyByte), Headers{
			"Accept":       "application/json",
			"content-type": "application/json",
		})
	if err != nil {
		return "", err
	}
	if code == 201 {
		var requestToken RequestTokenRes
		if err := json.Unmarshal([]byte(bodyStr), &requestToken); err == nil {
			return requestToken.Id, nil
		} else {
			return "", err
		}
	} else if code == 401 {
		return "", errors.New("Unauthorized")
	} else if code == 404 {
		return "", errors.New("Not Found")
	} else if code == 403 {
		return "", errors.New("Forbidden")
	} else {
		return "", errors.New(bodyStr)
	}
}

func RefreshAccessToken(token string) (string, error) {
	body := &RefreshTokenReq{
		RefreshToken: token,
	}
	bodyByte, _ := json.Marshal(body)
	code, bodyStr, err := Curl(Client,
		"POST",
		CDF_API_SERVER+"/suiteInstaller/urest/v1.1/tokens/refreshToken",
		string(bodyByte), Headers{
			"Accept":       "application/json",
			"content-type": "application/json",
		})
	if err != nil {
		return "", err
	}
	if code == 201 {
		var requestToken RequestTokenRes
		if err := json.Unmarshal([]byte(bodyStr), &requestToken); err == nil {
			return requestToken.Id, nil
		} else {
			return "", err
		}
	} else if code == 401 {
		return "", errors.New("Unauthorized")
	} else if code == 404 {
		return "", errors.New("Not Found")
	} else if code == 403 {
		return "", errors.New("Forbidden")
	} else {
		return "", errors.New(bodyStr)
	}
}

func GetTokenFromCookie(r *http.Request) string {
	return GetCookie(r, TOKEN_COOKIE_NAME)
}

func setTokenFromCookie(w http.ResponseWriter, token string) {
	SetCookie(w, TOKEN_COOKIE_NAME, token)
}

func deleteTokenFromCookie(w http.ResponseWriter) {
	DeleteCookie(w, TOKEN_COOKIE_NAME)
}

func getEnvOrDefault(envVar, defaultValue string) string {
	v := os.Getenv(envVar)
	if v == "" {
		LogW("%s=%s", envVar, defaultValue)
		return defaultValue
	}
	return v
}