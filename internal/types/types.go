// Code generated by goctl. DO NOT EDIT.
package types

type AddToOrder struct {
	Id string `path:"id"`
	Items []ItemCreate `json:"items"`
}

type OrderCreate struct {
	Table int    `json:"table"`
	Items []ItemCreate `json:"items"`
}

type Order struct {
	Id           string  `json:"id"`
	CreatedDate  int64   `json:"created_date"`
	FinishDate   int64   `json:"finish_date"`
	Table        int     `json:"table"`
	Items        []Item  `json:"item"`
	Paid         bool    `json:"paid"`
	Status       string  `json:"status"`
	SummaryPrice float64 `json:"summary_price"`
}

type Item struct {
	Id           string  `json:"id"`
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	UnitPrice float64 `json:"unit_piece"`
	IsDilivered bool `json:"delivered"`
	PositionId int `json:"position_id"`
}

type ItemCreate struct {
	Name       string  `json:"name"`
	Count      int     `json:"count"`
	UnitPrice float64 `json:"unit_piece"`
	IsDilivered bool `json:"delivered"`
	PositionId int `json:"position_id"`
}

type Table struct {
	Number int
	Token string
}

type TableNumber struct {
	Number int `path:"table"`
}