package bitflyer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var apiKey string
var apiSecret string

var logger = log.New(os.Stdout, "main: ", log.Lshortfile)

func init() {
	apiKey = os.Getenv("BITFLYER_APIKEY")
	apiSecret = os.Getenv("BITFLYER_APISECRET")
	if apiKey == "" && apiSecret == "" {
		logger.Fatal("BITFLYER_APIKEY or BITFLYER_APISECRET environment variable is not set ")
	}
}

// Executions is the response json format for getexecution API
type Executions []struct {
	ID                         int     `json:"id"`
	Side                       string  `json:"side"`
	Price                      float64 `json:"price"`
	Size                       float64 `json:"size"`
	ExecDate                   string  `json:"exec_date"`
	BuyChildOrderAcceptanceID  string  `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string  `json:"sell_child_order_acceptance_id"`
}

// Balance is responce json format for getbalance API
type Balance []struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

// Markets is responce json format for markets API
type Markets []struct {
	ProductCode string `json:"product_code"`
	Alias       string `json:"alias,omitempty"`
}

// Client is the http client for bitflyer
type Client struct {
	*http.Client
	baseURL string
	version string
}

// New create default Client for bitflyer
func New() *Client {
	return &Client{
		Client:  &http.Client{},
		baseURL: "https://api.bitflyer.com/",
		version: "v1",
	}
}

// TestPrint is just for test
func (c *Client) TestPrint() {
	fmt.Printf("test")
}

func createSign(method string, body string, path string) (timestamp string,
	sign string) {
	timestamp = fmt.Sprint(int32(time.Now().Unix()))

	text := timestamp + method + path + body
	h := hmac.New(sha256.New, []byte(apiSecret))
	h.Write([]byte(text))
	sign = hex.EncodeToString(h.Sum(nil))

	return timestamp, sign
}

// GetBalance send request to getbalance API
func (c *Client) GetBalance() (Balance, error) {
	var balance Balance
	path := c.getPath("me/getbalance")
	method := "GET"
	body := ""
	req, err := http.NewRequest(method, c.baseURL+path, nil)
	if err != nil {
		fmt.Println(err)
		return balance, err
	}

	req.Header.Add("ACCESS-KEY", apiKey)
	timestamp, sign := createSign(method, body, path)
	req.Header.Add("ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("ACCESS-SIGN", sign)

	if err != nil {
		fmt.Println(err)
		return balance, err
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return balance, err
	}
	fmt.Println(resp.Status)

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&balance); err != nil {
		fmt.Println(err)
		return balance, err
	}

	return balance, nil

}

func (c *Client) getPath(path string) string {
	return "/" + c.version + "/" + path
}

// GetExecutions get execution history
func (c *Client) GetExecutions() (Executions, error) {
	path := c.getPath("getexecutions")
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	var executions Executions

	if err != nil {
		fmt.Println(err)
		return executions, err
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return executions, err
	}
	fmt.Println(resp.Status)

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&executions); err != nil {
		fmt.Println(err)
		return executions, err
	}

	return executions, nil
}

// GetMarkets get markets infomation
func (c *Client) GetMarkets() (Markets, error) {
	path := c.getPath("getmarkets")
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	var markets Markets

	if err != nil {
		fmt.Println(err)
		return markets, err
	}

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return markets, err
	}
	fmt.Println(resp.Status)

	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&markets); err != nil {
		fmt.Println(err)
		return markets, err
	}

	return markets, nil

}
