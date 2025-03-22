package models

type OrderType int8

const (
	OrderBill       OrderType = 1
	OrderMembership OrderType = 2
)
