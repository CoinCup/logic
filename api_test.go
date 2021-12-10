package logic

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
)

func TestApi_GenerateDecimalWrongMarshal(t *testing.T) {
	marshal = func(v interface{}) ([]byte, error) {
		return nil, errors.New("test error")
	}

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateDecimal(context.Background(), 3)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	marshal = json.Marshal
}

func TestApi_GenerateDecimalWrongMarshal1(t *testing.T) {
	marshal1 = func(v interface{}) ([]byte, error) {
		return nil, errors.New("test error")
	}

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateDecimal(context.Background(), 3)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	marshal1 = json.Marshal
}

func TestApi_GenerateDecimalWrongUnmarshal(t *testing.T) {
	unmarshal = func(data []byte, v interface{}) error {
		return errors.New("test error")
	}

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateDecimal(context.Background(), 3)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	unmarshal = json.Unmarshal
}

func TestApi_GenerateDecimalWrongApiUrl(t *testing.T) {
	url := ApiUrl
	ApiUrl = "wrong"

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateDecimal(context.Background(), 3)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	ApiUrl = url
}

func TestApi_GenerateDecimalWrongApiMethod(t *testing.T) {
	method := ApiMethod
	ApiMethod = "вронг"

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateDecimal(context.Background(), 3)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	ApiMethod = method
}

func TestApi_GenerateDecimalWrongPlaces(t *testing.T) {
	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateDecimal(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApi_GenerateDecimal(t *testing.T) {
	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateDecimal(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApi_GenerateIntegerWrongMarshal(t *testing.T) {
	marshal = func(_ interface{}) ([]byte, error) {
		return nil, errors.New("test error")
	}

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateInteger(context.Background(), 0, 53)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	marshal = json.Marshal
}

func TestApi_GenerateIntegerWrongMarshal1(t *testing.T) {
	marshal1 = func(_ interface{}) ([]byte, error) {
		return nil, errors.New("test error")
	}

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateInteger(context.Background(), 0, 53)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	marshal1 = json.Marshal
}

func TestApi_GenerateIntegerWrongUnmarshal(t *testing.T) {
	unmarshal = func(_ []byte, _ interface{}) error {
		return errors.New("test error")
	}

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateInteger(context.Background(), 0, 53)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	unmarshal = json.Unmarshal
}

func TestApi_GenerateIntegerWrongApiUrl(t *testing.T) {
	url := ApiUrl
	ApiUrl = "wrong"

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateInteger(context.Background(), 0, 53)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	ApiUrl = url
}

func TestApi_GenerateIntegerWrongApiMethod(t *testing.T) {
	method := ApiMethod
	ApiMethod = "вронг"

	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateInteger(context.Background(), 0, 53)
	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	ApiMethod = method
}

func TestApi_GenerateInteger(t *testing.T) {
	api := NewApi(os.Getenv("API_KEY"))
	_, err := api.GenerateInteger(context.Background(), 0, 53)
	if err != nil {
		t.Fatal(err)
	}
}
