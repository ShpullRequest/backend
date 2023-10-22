package ip

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type ipResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func GetGeoByIP(ip string) (string, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	response, err := resty.New().R().
		SetResult(&ipResponse{}).
		Get(url)
	if err != nil {
		return "", err
	}

	result := response.Result().(*ipResponse)
	if result.Lat == 0 && result.Lon == 0 {
		return "", errors.New("invalid geo")
	}

	return fmt.Sprintf("%f, %f", result.Lat, result.Lon), nil
}
