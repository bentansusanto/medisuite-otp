package config

var Config AppConfig

type AppConfig struct {
	RateLimiterMaxRequest float64 `json:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond int     `json:"rateLimiterTimeSecond"`
	Environment           string  `json:"environment"`  // "development", "production", etc.
	CookieDomain          string  `json:"cookieDomain"` // Cookie domain for production
	CookieSecure          bool    `json:"cookieSecure"` // HTTPS requirement for cookies
}
