package geoip

import (
	. "geoip/services"
)

type Handler struct {
}

func NewHandler(dbPath string) (h *Handler) {
	InitMaxMindDb(dbPath)
	return &Handler{}
}

// Parameters:
//  - IP
func (h *Handler) IpToGeoData(ip string) (r *GeoData, err error) {
	r, err = IpToGeoData(ip)
	return
}

// Parameters:
//  - IP
func (h *Handler) GetLatlng(ip string) (r *LatLng, err error) {
	var g *GeoData
	g, err = IpToGeoData(ip)

	if err != nil {
		return nil, err
	} else {
		return &LatLng{g.Lat, g.Lng}, nil
	}
}

// Parameters:
//  - IP
func (h *Handler) GetCityName(ip string) (r string, err error) {
	var g *GeoData
	g, err = IpToGeoData(ip)

	if err != nil {
		return "", err
	} else {
		return g.CityName, nil
	}
}

// Parameters:
//  - IP
func (h *Handler) GetCountryName(ip string) (r string, err error) {
	var g *GeoData
	g, err = IpToGeoData(ip)

	if err != nil {
		return "", err
	} else {
		return g.CountryName, nil
	}
}

// Parameters:
//  - IP
func (h *Handler) GetCountryCode(ip string) (r string, err error) {
	var g *GeoData
	g, err = IpToGeoData(ip)

	if err != nil {
		return "", err
	} else {
		return g.CountryIsoCode, nil
	}
}

// Parameters:
//  - IP
func (h *Handler) GetProvince(ip string) (r string, err error) {
	var g *GeoData
	g, err = IpToGeoData(ip)

	if err != nil {
		return "", err
	} else {
		return g.Province, nil
	}
}

func (h *Handler) Ping() (err error) {
	return nil
}
