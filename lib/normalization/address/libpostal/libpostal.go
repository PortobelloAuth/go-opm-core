package libpostal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PortobelloAuth/go-opm-core/lib/normalization/address/libpostal/addressinfo"
)

type LibpostalAddressPart struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

const (
	// from https://github.com/openvenues/libpostal/blob/master/src/address_parser.h
	LIBPOSTAL_LABEL_HOUSE          = "house"
	LIBPOSTAL_LABEL_HOUSE_NUMBER   = "house_number"
	LIBPOSTAL_LABEL_PO_BOX         = "po_box"
	LIBPOSTAL_LABEL_BUILDING       = "building"
	LIBPOSTAL_LABEL_ENTRANCE       = "entrance"
	LIBPOSTAL_LABEL_STAIRCASE      = "staircase"
	LIBPOSTAL_LABEL_LEVEL          = "level"
	LIBPOSTAL_LABEL_UNIT           = "unit"
	LIBPOSTAL_LABEL_ROAD           = "road"
	LIBPOSTAL_LABEL_METRO_STATION  = "metro_station"
	LIBPOSTAL_LABEL_SUBURB         = "suburb"
	LIBPOSTAL_LABEL_CITY_DISTRICT  = "city_district"
	LIBPOSTAL_LABEL_CITY           = "city"
	LIBPOSTAL_LABEL_STATE_DISTRICT = "state_district"
	LIBPOSTAL_LABEL_ISLAND         = "island"
	LIBPOSTAL_LABEL_STATE          = "state"
	LIBPOSTAL_LABEL_POSTAL_CODE    = "postcode"
	LIBPOSTAL_LABEL_COUNTRY_REGION = "country_region"
	LIBPOSTAL_LABEL_COUNTRY        = "country"
	LIBPOSTAL_LABEL_WORLD_REGION   = "world_region"

	LIBPOSTAL_LABEL_WEBSITE   = "website"
	LIBPOSTAL_LABEL_TELEPHONE = "phone"
)

type LibpostalHttpService struct {
	BaseURL *url.URL
	Timeout time.Duration
}

// NOTE: libpostal does not implement the Project US@ specification!! This is an approximate
// implementation intended to begin to excercise the go-opm-core address normalization API

func NewService(base string, timeout time.Duration) (*LibpostalHttpService, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse base URL for service")
	}

	return &LibpostalHttpService{
		BaseURL: u,
		Timeout: timeout,
	}, nil
}

func (l *LibpostalHttpService) call(endpoint string, address string) ([]byte, error) {
	client := &http.Client{Timeout: l.Timeout}

	p := endpoint
	if len(l.BaseURL.Path) > 0 {
		bp, _ := strings.CutPrefix(l.BaseURL.Path, "/")
		bp, _ = strings.CutSuffix(l.BaseURL.Path, "/")
		p = fmt.Sprintf("%s/%s", bp, endpoint)
	}

	u := &url.URL{
		Scheme: l.BaseURL.Scheme,
		Host:   l.BaseURL.Host,
		Path:   p,
	}

	q := u.Query()
	q.Set("address", address)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Request failed: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read response body: %w", err)
	}

	return body, nil
}

func (l *LibpostalHttpService) Parse(address string) ([]LibpostalAddressPart, error) {
	body, err := l.call("parse", address)
	if err != nil {
		return nil, err
	}

	var parts []LibpostalAddressPart
	err = json.Unmarshal(body, &parts)
	if err != nil {
		return nil, fmt.Errorf("Unable to process parse endpoint response: %w", err)
	}

	return parts, nil
}

func (l *LibpostalHttpService) Expand(address string) ([]string, error) {
	body, err := l.call("expand", address)
	if err != nil {
		return nil, err
	}

	var expanded []string
	err = json.Unmarshal(body, &expanded)
	if err != nil {
		return nil, fmt.Errorf("Unable to process expand endpoint response: %w", err)
	}

	return expanded, nil
}

func (l *LibpostalHttpService) Apply(input string) (string, error) {
	parsed, err := l.Parse(input)
	if err != nil {
		return "", err
	}

	a := &addressinfo.Address{}
	// FIXME: this is an imperfect mapping!!
	for _, v := range parsed {
		switch v.Label {
		case LIBPOSTAL_LABEL_HOUSE:
			a.BuisinessName = v.Value
		case LIBPOSTAL_LABEL_HOUSE_NUMBER:
			a.PrimaryNumber = v.Value
		case LIBPOSTAL_LABEL_UNIT:
			a.SecondaryDesignator = "#"
			a.SecondaryNumber = v.Value
		case LIBPOSTAL_LABEL_ROAD:
			a.StreetName = v.Value
		case LIBPOSTAL_LABEL_CITY:
			a.City = v.Value
		case LIBPOSTAL_LABEL_STATE:
			a.Region = v.Value
		case LIBPOSTAL_LABEL_POSTAL_CODE:
			a.Postal = v.Value
		case LIBPOSTAL_LABEL_COUNTRY:
			a.Country = v.Value

		// case LIBPOSTAL_LABEL_PO_BOX:
		// case LIBPOSTAL_LABEL_BUILDING:
		// case LIBPOSTAL_LABEL_ENTRANCE:
		// case LIBPOSTAL_LABEL_STAIRCASE:
		// case LIBPOSTAL_LABEL_LEVEL:
		// case LIBPOSTAL_LABEL_METRO_STATION:
		// case LIBPOSTAL_LABEL_SUBURB:
		// case LIBPOSTAL_LABEL_CITY_DISTRICT:
		// case LIBPOSTAL_LABEL_STATE_DISTRICT:
		// case LIBPOSTAL_LABEL_ISLAND:
		// case LIBPOSTAL_LABEL_COUNTRY_REGION:
		// case LIBPOSTAL_LABEL_WORLD_REGION:
		default:
		}
	}

	result := a.Format()
	return result, nil
}
