package devices

// DatabaseConfig contains configuration settings for database access.
type DatabaseConfig struct {
	Name     string `json:"DbName"`
	User     string `json:"DbUser"`
	Password string `json:"DbPassword"`
	SiteName string
}
