package graph

import time "time"

type Order struct {
	ID         string           `json:"id"`
	CreatedAt  time.Time        `json:"createdAt"`
	TotalPrice float64          `json:"totalPrice"`
	Products   []OrderedProduct `json:"products"`
}
