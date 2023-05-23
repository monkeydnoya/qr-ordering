package types

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type OrderEntity struct {
	Id           uuid.UUID    `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedDate  time.Time    `json:"created_date"`
	Table        int          `json:"table"`
	Paid         bool         `json:"paid"`
	Status       string       `json:"status"`
	SummaryPrice float64      `json:"summary_price"`
	Items        []ItemEntity `json:"item" gorm:"foreignKey:OrderId"`
}

func (oe *OrderEntity) TableName() string {
	return "orders"
}

func (oe *OrderEntity) ToModel() CreatedOrder {
	items := make([]Item, 0)
	for _, ie := range oe.Items {
		items = append(items, ie.ToModel())
	}
	return CreatedOrder{
		Id:           oe.Id.String(),
		CreatedDate:  oe.CreatedDate.Format(time.RFC3339),
		Table:        oe.Table,
		Items:        items,
		Paid:         oe.Paid,
		Status:       oe.Status,
		SummaryPrice: oe.SummaryPrice,
	}
}
