package types

import (
	uuid "github.com/satori/go.uuid"
)

type ItemEntity struct {
	ItemId       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name         string    `json:"name"`
	Count        int       `json:"count"`
	PiecePrice   float64   `json:"price_per_piece"`
	SummaryPrice float64   `json:"summary_price"`
	OrderId      uuid.UUID `json:"order_id"`
}

func (ie *ItemEntity) TableName() string {
	return "items"
}

func (ie *ItemEntity) ToModel() Item {
	return Item{
		Name:       ie.Name,
		Count:      ie.Count,
		PiecePrice: ie.PiecePrice,
	}
}

func (i *Item) ToEntity() ItemEntity {
	summaryPrice := float64(i.Count) * i.PiecePrice
	return ItemEntity{
		Name:         i.Name,
		Count:        i.Count,
		PiecePrice:   i.PiecePrice,
		SummaryPrice: summaryPrice,
	}
}
