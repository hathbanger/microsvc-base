package models

// HealthRequest - response for health
type HealthRequest struct{}

// HealthResponse - response for health
type HealthResponse struct {
	Health bool `json:"health"`
}
