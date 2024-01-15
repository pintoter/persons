package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ageURL         = "https://api.agify.io/?name=%s"
	genderURL      = "https://api.genderize.io/?name=%s"
	nationalizeURL = "https://api.nationalize.io/?name=%s"
)

type Client struct {
	httpClient *http.Client
}

type Config interface {
	GetTimeout() time.Duration
}

func New(cfg Config) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: cfg.GetTimeout(),
		},
	}
}

type getAgeContract struct {
	Age int `json:"age"`
}

func (c *Client) GetAge(name string) (int, error) {
	data, err := c.requestExecutor(fmt.Sprintf(ageURL, name))
	if err != nil {
		return 0, err
	}

	contract := new(getAgeContract)
	if err := json.Unmarshal(data, contract); err != nil {
		return 0, err
	}

	return contract.Age, nil
}

type getGenderContract struct {
	Gender string `json:"gender"`
}

func (c *Client) GetGender(name string) (string, error) {
	data, err := c.requestExecutor(fmt.Sprintf(genderURL, name))
	if err != nil {
		return "", err
	}

	contract := new(getGenderContract)
	if err := json.Unmarshal(data, contract); err != nil {
		return "", err
	}

	return contract.Gender, nil
}

type getNationalizeContract struct {
	Nationalize string `json:"nationalize"`
}

func (c *Client) GetNationalize(name string) (string, error) {
	data, err := c.requestExecutor(fmt.Sprintf(nationalizeURL, name))
	if err != nil {
		return "", err
	}

	contract := new(getNationalizeContract)
	if err := json.Unmarshal(data, contract); err != nil {
		return "", err
	}

	return contract.Nationalize, nil
}

func (c *Client) requestExecutor(url string) ([]byte, error) {
	// using context.Background() because Client already got parameter Timeout in constructor
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
