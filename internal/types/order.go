package types

import (
	uuid "github.com/satori/go.uuid"
)

type OrderEntity struct {
	Id           uuid.UUID    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedDate  int64        `json:"created_date"`
	FinishDate   int64        `json:"finish_date"`
	Table        int          `json:"table"`
	Paid         bool         `json:"paid"`
	Status       string       `json:"status"`
	SummaryPrice float64      `json:"summary_price"`
	Items        []ItemEntity `json:"item" gorm:"foreignKey:OrderId"`
}

func (oe *OrderEntity) TableName() string {
	return "orders"
}

func (oe *OrderEntity) ToModel() Order {
	items := make([]Item, 0)
	for _, ie := range oe.Items {
		items = append(items, ie.ToModel())
	}
	return Order{
		Id:           oe.Id.String(),
		CreatedDate:  oe.CreatedDate,
		FinishDate:   oe.FinishDate,
		Table:        oe.Table,
		Items:        items,
		Paid:         oe.Paid,
		Status:       oe.Status,
		SummaryPrice: oe.SummaryPrice,
	}
}
