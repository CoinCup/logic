package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ApiUrl    = "https://api.random.org/json-rpc/4/invoke"
	ApiMethod = http.MethodPost
)

type Api struct {
	key string
}

func NewApi(apiKey string) *Api {
	return &Api{key: apiKey}
}

var (
	marshal   = json.Marshal
	marshal1  = json.Marshal
	unmarshal = json.Unmarshal
)

type ApiError struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func (err *ApiError) Error() string {
	return fmt.Sprintf("error %d: %s", err.Code, err.Message)
}

type Decimal struct {
	Value        float64 `json:"value"`
	Random       string  `json:"random"`
	Signature    string  `json:"signature"`
	SerialNumber uint64  `json:"serial_number"`
}

type decimalRequestParams struct {
	ApiKey        string `json:"apiKey"`
	N             uint   `json:"n"`
	DecimalPlaces uint   `json:"decimalPlaces"`
}

type decimalRequest struct {
	JsonRPC string               `json:"jsonrpc"`
	Method  string               `json:"method"`
	Params  decimalRequestParams `json:"params"`
	ID      int                  `json:"id"`
}

type decimalResponseRandom struct {
	Method                    string                 `json:"method"`
	HashedApiKey              string                 `json:"hashedApiKey"`
	N                         int                    `json:"n"`
	DecimalPlaces             int                    `json:"decimalPlaces"`
	Replacement               bool                   `json:"replacement"`
	PregeneratedRandomization map[string]interface{} `json:"pregeneratedRandomization"`
	Data                      []float64              `json:"data"`
	License                   map[string]interface{} `json:"license"`
	LicenseData               map[string]interface{} `json:"licenseData"`
	UserData                  map[string]interface{} `json:"userData"`
	TicketData                map[string]interface{} `json:"ticketData"`
	CompletionTime            string                 `json:"completionTime"`
	SerialNumber              uint64                 `json:"serialNumber"`
}

type decimalResponseResult struct {
	Random        decimalResponseRandom `json:"random"`
	Signature     string                `json:"signature"`
	Cost          float64               `json:"cost"`
	BitsUsed      int                   `json:"bits_used"`
	BitsLeft      int                   `json:"bits_left"`
	RequestsLeft  int                   `json:"requests_left"`
	AdvisoryDelay int                   `json:"advisory_delay"`
}

type decimalResponse struct {
	JsonRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Result  *decimalResponseResult `json:"result"`
	Error   *ApiError              `json:"error"`
	ID      int                    `json:"id"`
}

func (api *Api) GenerateDecimal(ctx context.Context, decimalPlaces uint) (*Decimal, error) {
	if decimalPlaces == 0 || decimalPlaces > 8 {
		decimalPlaces = 8
	}

	requestData := decimalRequest{
		JsonRPC: "2.0",
		Method:  "generateSignedDecimalFractions",
		Params: decimalRequestParams{
			ApiKey:        api.key,
			N:             1,
			DecimalPlaces: decimalPlaces,
		},
		ID: 1337,
	}

	requestBytes, err := marshal(requestData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		ApiMethod,
		ApiUrl,
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	responseData := &decimalResponse{}
	err = unmarshal(body, responseData)
	if err != nil {
		return nil, err
	}

	if responseData.Error != nil {
		return nil, responseData.Error
	}

	randomBytes, err := marshal1(responseData.Result.Random)
	if err != nil {
		return nil, err
	}

	if responseData.Result == nil || len(responseData.Result.Random.Data) == 0 {
		return nil, errors.New("wrong response")
	}

	decimal := &Decimal{
		Value:        responseData.Result.Random.Data[0],
		Random:       string(randomBytes),
		Signature:    responseData.Result.Signature,
		SerialNumber: responseData.Result.Random.SerialNumber,
	}

	return decimal, nil
}

type Integer struct {
	Value        int    `json:"value"`
	Random       string `json:"random"`
	Signature    string `json:"signature"`
	SerialNumber uint64 `json:"serial_number"`
}

type integerRequestParams struct {
	ApiKey string `json:"apiKey"`
	N      int    `json:"n"`
	Min    int    `json:"min"`
	Max    int    `json:"max"`
}

type integerRequest struct {
	JsonRPC string               `json:"jsonrpc"`
	Method  string               `json:"method"`
	Params  integerRequestParams `json:"params"`
	ID      int                  `json:"id"`
}

type integerResponseRandom struct {
	Method                    string                 `json:"method"`
	HashedApiKey              string                 `json:"hashedApiKey"`
	N                         int                    `json:"n"`
	Min                       int                    `json:"min"`
	Max                       int                    `json:"max"`
	Replacement               bool                   `json:"replacement"`
	Base                      int                    `json:"base"`
	PregeneratedRandomization map[string]interface{} `json:"pregeneratedRandomization"`
	Data                      []int                  `json:"data"`
	License                   map[string]interface{} `json:"license"`
	LicenseData               map[string]interface{} `json:"licenseData"`
	UserData                  map[string]interface{} `json:"userData"`
	TicketData                map[string]interface{} `json:"ticketData"`
	CompletionTime            string                 `json:"completionTime"`
	SerialNumber              uint64                 `json:"serialNumber"`
}

type integerResponseResult struct {
	Random        integerResponseRandom `json:"random"`
	Signature     string                `json:"signature"`
	Cost          float64               `json:"cost"`
	BitsUsed      int                   `json:"bitsUsed"`
	BitsLeft      int                   `json:"bitsLeft"`
	RequestsLeft  int                   `json:"requestsLeft"`
	AdvisoryDelay int                   `json:"advisoryDelay"`
}

type integerResponse struct {
	JsonRPC string                 `json:"jsonrpc"`
	Result  *integerResponseResult `json:"result"`
	Error   *ApiError              `json:"error"`
	ID      int                    `json:"id"`
}

func (api *Api) GenerateInteger(ctx context.Context, min int, max int) (*Integer, error) {
	requestData := &integerRequest{
		JsonRPC: "2.0",
		Method:  "generateSignedIntegers",
		Params: integerRequestParams{
			ApiKey: api.key,
			N:      1,
			Min:    min,
			Max:    max,
		},
	}

	requestBytes, err := marshal(requestData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		ApiMethod,
		ApiUrl,
		bytes.NewBuffer(requestBytes),
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	responseData := &integerResponse{}
	err = unmarshal(body, responseData)
	if err != nil {
		return nil, err
	}

	if responseData.Error != nil {
		return nil, responseData.Error
	}

	randomBytes, err := marshal1(responseData.Result.Random)
	if err != nil {
		return nil, err
	}

	if responseData.Result == nil || len(responseData.Result.Random.Data) == 0 {
		return nil, errors.New("wrong response")
	}

	integer := &Integer{
		Value:        responseData.Result.Random.Data[0],
		Random:       string(randomBytes),
		Signature:    responseData.Result.Signature,
		SerialNumber: responseData.Result.Random.SerialNumber,
	}

	return integer, nil
}
