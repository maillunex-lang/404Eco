package models

import "time"

type City struct {
	ID           string             `json:"id"`
	NameRU       string             `json:"nameRu"`
	NameBA       string             `json:"nameBa"`
	X            float64            `json:"x"`
	Y            float64            `json:"y"`
	Population   int                `json:"population"`
	Wealth       float64            `json:"wealth"`
	IndustryBias map[string]float64 `json:"industryBias"`
}

type Product struct {
	ID        string  `json:"id"`
	NameRU    string  `json:"nameRu"`
	Category  string  `json:"category"`
	BasePrice float64 `json:"basePrice"`
	Unit      string  `json:"unit"`
}

type Player struct {
	ID          string    `json:"id"`
	NakamaID    string    `json:"nakamaId,omitempty"`
	DisplayName string    `json:"displayName"`
	Cash        float64   `json:"cash"`
	Reputation  float64   `json:"reputation"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Firm struct {
	ID         string             `json:"id"`
	PlayerID   string             `json:"playerId"`
	Name       string             `json:"name"`
	CityID     string             `json:"cityId"`
	Industry   string             `json:"industry"`
	Level      int                `json:"level"`
	Employees  int                `json:"employees"`
	Capacity   int                `json:"capacity"`
	Inventory  map[string]float64 `json:"inventory"`
	Pricing    map[string]float64 `json:"pricing"`
	Quality    float64            `json:"quality"`
	Brand      float64            `json:"brand"`
	CashflowD  float64            `json:"cashflowDaily"`
	IsPublic   bool               `json:"isPublic"`
	SharePrice float64            `json:"sharePrice"`
	SharesOut  int                `json:"sharesOutstanding"`
	PropertyID string             `json:"propertyId,omitempty"`
}

type Property struct {
	ID        string  `json:"id"`
	CityID    string  `json:"cityId"`
	Kind      string  `json:"kind"`
	Name      string  `json:"name"`
	Size      int     `json:"size"`
	BaseRent  float64 `json:"baseRent"`
	BuyPrice  float64 `json:"buyPrice"`
	OwnerFirm string  `json:"ownerFirm,omitempty"`
	Condition float64 `json:"condition"`
}

type MarketOrder struct {
	ID        string    `json:"id"`
	FirmID    string    `json:"firmId"`
	ProductID string    `json:"productId"`
	Side      string    `json:"side"`
	Qty       int       `json:"qty"`
	Price     float64   `json:"price"`
	Filled    int       `json:"filled"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type WorldState struct {
	Tick       int64     `json:"tick"`
	Day        int       `json:"day"`
	Month      int       `json:"month"`
	Year       int       `json:"year"`
	OilPrice   float64   `json:"oilPrice"`
	GrainIndex float64   `json:"grainIndex"`
	LaborWage  float64   `json:"laborWage"`
	Inflation  float64   `json:"inflation"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type NewsItem struct {
	ID       int64     `json:"id"`
	Tick     int64     `json:"tick"`
	Severity string    `json:"severity"`
	Title    string    `json:"title"`
	Body     string    `json:"body"`
	Created  time.Time `json:"createdAt"`
}

type Snapshot struct {
	World      WorldState            `json:"world"`
	Cities     []City                `json:"cities"`
	Products   []Product             `json:"products"`
	Players    []Player              `json:"players"`
	Firms      []Firm                `json:"firms"`
	Properties []Property            `json:"properties"`
	Orders     []MarketOrder         `json:"orders"`
	News       []NewsItem            `json:"news"`
	Prices     map[string]float64    `json:"marketPrices"`
	Demand     map[string]float64    `json:"marketDemand"`
}
