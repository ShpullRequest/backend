package maps

import (
	"errors"
	"fmt"
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/go-resty/resty/v2"
	"net/url"
	"strings"
)

type vkMaps struct {
	apiKey string
	client *resty.Client
}

func New(cfg config.NodeConfig) *vkMaps {
	return &vkMaps{
		apiKey: cfg.VkMapsAPIKey,
		client: resty.New().SetBaseURL("https://maps.vk.com/api/"),
	}
}

type SearchResponse struct {
	Request string `json:"request"`
	Results []struct {
		AddressDetails struct {
			Building   string `json:"building,omitempty"`
			Country    string `json:"country"`
			Locality   string `json:"locality"`
			Region     string `json:"region"`
			Street     string `json:"street"`
			Suburb     string `json:"suburb,omitempty"`
			PostalCode string `json:"postal_code,omitempty"`
		} `json:"address_details"`
		Name string    `json:"name,omitempty"`
		Pin  []float64 `json:"pin"`
		Type string    `json:"type"`
	} `json:"results"`
}

func (m *vkMaps) GetAddressByGeo(lng, lat float64) (string, error) {
	params := url.Values{}
	params.Set("api_key", m.apiKey)
	params.Set("q", fmt.Sprintf("%f,%f", lng, lat))

	response, err := m.client.R().
		SetQueryParamsFromValues(params).
		SetResult(&SearchResponse{}).
		Get("search")
	if err != nil {
		return "", err
	}

	searchResult := response.Result().(*SearchResponse)
	if len(searchResult.Results) < 1 {
		return "", errors.New("no address")
	}

	return m.parseAddress(searchResult), nil
}
func (m *vkMaps) GetAddressByGeoQ(q string) (string, error) {
	params := url.Values{}
	params.Set("api_key", m.apiKey)
	params.Set("q", q)

	response, err := m.client.R().
		SetQueryParamsFromValues(params).
		SetResult(&SearchResponse{}).
		Get("search")
	if err != nil {
		return "", err
	}

	searchResult := response.Result().(*SearchResponse)
	if len(searchResult.Results) < 1 {
		return "", errors.New("no address")
	}

	return m.parseAddress(searchResult), nil
}

func (m *vkMaps) parseAddress(result *SearchResponse) string {
	addressDetails := result.Results[0].AddressDetails

	var texts []string
	if addressDetails.Country != "" {
		texts = append(texts, addressDetails.Country)
	}
	if addressDetails.Region != "" {
		texts = append(texts, addressDetails.Region)
	}
	if addressDetails.Suburb != "" {
		texts = append(texts, addressDetails.Suburb)
	}
	if addressDetails.Street != "" {
		texts = append(texts, addressDetails.Street)
	}
	if addressDetails.Building != "" {
		texts = append(texts, addressDetails.Building)
	}

	return strings.Join(texts, ", ")
}
