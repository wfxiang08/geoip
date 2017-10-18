package geoip

import (
	"context"
	. "geoip/services"
)

type Handler struct {
}

func NewHandler(dbPath string) (h *Handler) {
	InitMaxMindDb(dbPath, true)
	return &Handler{}
}

// Parameters:
//  - IP
func (h *Handler) IpToGeoData(ctx context.Context, ip string) (r *GeoData, err error) {
	r, err = IpToGeoData(ip)
	return
}

// Parameters:
//  - IP
func (h *Handler) GetLatlng(ctx context.Context, ip string) (r *LatLng, err error) {
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
func (h *Handler) GetCityName(ctx context.Context, ip string) (r string, err error) {
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
func (h *Handler) GetCountryName(ctx context.Context, ip string) (r string, err error) {
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
func (h *Handler) GetCountryCode(ctx context.Context, ip string) (r string, err error) {
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
func (h *Handler) GetProvince(ctx context.Context, ip string) (r string, err error) {
	var g *GeoData
	g, err = IpToGeoData(ip)

	if err != nil {
		return "", err
	} else {
		return g.Province, nil
	}
}

func (h *Handler) Ping(ctx context.Context) (err error) {
	return nil
}
