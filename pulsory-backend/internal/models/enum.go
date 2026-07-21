package models

type WebsiteStatus string

const (
	StatusUp      WebsiteStatus = "Up"
	StatusDown    WebsiteStatus = "Down"
	StatusUnknown WebsiteStatus = "Unknown"
)