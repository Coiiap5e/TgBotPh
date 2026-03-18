package model

import "time"

type Shoot struct {
	Id            int
	ShootDate     time.Time
	StartTime     time.Time
	EndTime       time.Time
	ShootPrice    int
	ShootLocation string
	ShootType     string
	Notes         string
	PriceUSD      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Clients []ShootClientInfo
}

type ShootClientInfo struct {
	ClientID         int
	FirstName        string
	LastName         string
	Phone            string
	IsMainClient     bool
	RelationshipType string
}

type Client struct {
	Id               int
	FirstName        string
	LastName         string
	Phone            string
	SocialNetworkUrl string
	CreatedAt        time.Time
}