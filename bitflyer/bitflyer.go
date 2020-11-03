package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseURL = "https://api.bitflyer.com/"

type APIClient struct {
	key        string
	secret     string
	httpClient *http.Client
}

func New(key, secret string) *APIClient {
	return &APIClient{key, secret, &http.Client{}}
}

func (api APIClient) header(method, endpoint string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	log.Println(timestamp)
	message := timestamp + method + endpoint + string(body)

	mac := hmac.New(sha256.New, []byte(api.secret))
	mac.Write([]byte(message))
	sign := hex.EncodeToString(mac.Sum(nil))
	return map[string]string{
		"ACCESS-KEY":       api.key,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-type":     "application/json",
	}
}

func (api *APIClient) doRequest(method, urlPath string, query map[string]string, data []byte) (body []byte, err error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		return
	}
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}

	endPoint := baseURL.ResolveReference(apiURL).String()
	fmt.Println(endPoint)

	req, err := http.NewRequest(method, endPoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	for key, value := range api.header(method, req.URL.RequestURI(), data) {
		req.Header.Add(key, value)
	}
	//ここでリクエストを送っている
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}

type Balance struct {
	CurrentCode string  `json:"currency_code"`
	Amount      float64 `json:"amount"`
	Available   string  `json:"available"`
}

func (api *APIClient) GetBalance() ([]Balance, error) {
	url := "v1/me/getbalance"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}
	var balance []Balance
	err = json.Unmarshal(resp, &balance)
	if err != nil {
		log.Printf("action=GetBalance err=%s", err.Error())
		return nil, err
	}
	return balance, nil
}

//自分で認証の必要ないAPIにアクセスしてみる
//"v1/markets"にしてみる

type Market struct {
	PRODUCT_CODE string `json:"product_code"`
	MARKET_TYPE string `json:"market_type"`
}

func (api *APIClient) GetMarket() ([]Market, error) {
	url := "v1/markets"
	resp, err := api.doRequest("GET", url, map[string]string{}, nil)
	if err != nil {
		log.Printf("action=GetMarket err=%s", err.Error())
		return nil, err
	}
	fmt.Println(string(resp))
	var Market []Market
	err = json.Unmarshal(resp, &Market)
	if err != nil {
		log.Printf("action=GetMarket err=%s", err.Error())
		return nil, err
	}
	return Market, nil
}










