package types

import (
	uuid "github.com/satori/go.uuid"
)

type ItemEntity struct {
	ItemId       uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name         string    `json:"name"`
	Count        int       `json:"count"`
	UnitPrice    float64   `json:"unit_piece"`
	SummaryPrice float64   `json:"summary_price"`
	OrderId      uuid.UUID `json:"order_id"`
	IsDilivered  bool      `json:"delivered"`
	PositionId   int       `json:"position_id"`
}

func (ie *ItemEntity) TableName() string {
	return "items"
}

func (ie *ItemEntity) ToModel() Item {
	return Item{
		Id:          ie.ItemId.String(),
		Name:        ie.Name,
		Count:       ie.Count,
		UnitPrice:   ie.UnitPrice,
		IsDilivered: ie.IsDilivered,
		PositionId:  ie.PositionId,
	}
}

func (i *Item) ToEntity() ItemEntity {
	summaryPrice := float64(i.Count) * i.UnitPrice
	itemId, _ := uuid.FromString(i.Id)
	return ItemEntity{
		ItemId:       itemId,
		Name:         i.Name,
		Count:        i.Count,
		UnitPrice:    i.UnitPrice,
		SummaryPrice: summaryPrice,
		IsDilivered:  i.IsDilivered,
		PositionId:   i.PositionId,
	}
}
