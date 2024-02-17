package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pintoter/persons/services/command/internal/entity"
)

const (
	ageURL         = "https://api.agify.io/"
	genderURL      = "https://api.genderize.io/"
	nationalizeURL = "https://api.nationalize.io/"
	pathParamName  = "?name=%s"
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

func (c *Client) GenerateAge(ctx context.Context, name string) (int, error) {
	data, err := c.requestExecutor(ctx, fmt.Sprintf(ageURL+pathParamName, name))
	if err != nil {
		return 0, err
	}

	contract := new(getAgeContract)
	err = json.Unmarshal(data, contract)
	if err != nil {
		return 0, err
	}

	return contract.Age, nil
}

type getGenderContract struct {
	Gender string `json:"gender"`
}

func (c *Client) GenerateGender(ctx context.Context, name string) (string, error) {
	data, err := c.requestExecutor(ctx, fmt.Sprintf(genderURL+pathParamName, name))
	if err != nil {
		return "", err
	}

	contract := new(getGenderContract)
	err = json.Unmarshal(data, contract)
	if err != nil {
		return "", err
	}

	return contract.Gender, nil
}

type getNationalizeContract struct {
	Country []entity.Nationality `json:"country"`
}

func (c *Client) GenerateNationalize(ctx context.Context, name string) ([]entity.Nationality, error) {
	data, err := c.requestExecutor(ctx, fmt.Sprintf(nationalizeURL+pathParamName, name))
	if err != nil {
		return nil, err
	}

	contract := new(getNationalizeContract)
	err = json.Unmarshal(data, &contract)
	if err != nil {
		return nil, err
	}

	if len(contract.Country) == 0 {
		return nil, err
	}

	return contract.Country, nil
}

func (c *Client) requestExecutor(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
