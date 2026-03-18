package entity

type OrderEvent struct {
	OrderID     string
	OrderNumber string
	Status      OrderStatus
}
