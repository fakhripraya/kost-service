package entities

// Geolocation is an entity to communicate with the matched geolocation response data from any request
type Geolocation struct {
	GeoData []GeolocationDetail `json:"data"`
}

// GeolocationDetail is an entity to communicate with the matched geolocation detail response data from any request
type GeolocationDetail struct {
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Type               string  `json:"type"`
	Distance           float64 `json:"distance"`
	Name               string  `json:"name"`
	Number             string  `json:"number"`
	PostalCode         string  `json:"postal_code"`
	Street             string  `json:"street"`
	Confidence         float64 `json:"confidence"`
	Region             string  `json:"region"`
	RegionCode         string  `json:"region_code"`
	County             string  `json:"county"`
	Locality           string  `json:"locality"`
	AdministrativeArea string  `json:"administrative_area"`
	Neighbourhood      string  `json:"neighbourhood"`
	Country            string  `json:"country"`
	CountryCode        string  `json:"country_code"`
	Continent          string  `json:"continent"`
	Label              string  `json:"label"`
}
