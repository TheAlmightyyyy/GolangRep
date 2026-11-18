package enum

import (
	"math"
	"time"

	"github.com/google/uuid"
)

type Side int8

const (
	SideBuy  = iota + 1 // 1 - статус "покупка"
	SideSell            // 2 - статус "продажа"
)

// Stringer
// type Stringer interface
//
//	{
//		String() string
//	}
//
// slog.Info("Заказ создан", "type", order.Side)
func (side Side) String() string {
	switch side {
	case SideBuy:
		return "покупка"
	case SideSell:
		return "продажа"
	default:
		return "операция не определена"
	}
}

type OrderType int8

const (
	OrderTypeLimit  OrderType = iota + 1 //заявка по фиксированной цене
	OrderTypeMarket                      //заявка по рыночной цене
)

func (t OrderType) String() string {
	switch t {
	case OrderTypeLimit:
		return "заявка на фиксированную цену"
	case OrderTypeMarket:
		return "заявка по рыночной стоимости"
	default:
		return "тип заявки не определён"
	}
}

type OrderStatus int8

const (
	OrderStatusNew             OrderStatus = iota + 1 //зарегестрирована на бирже
	OrderStatusPartiallyFilled                        //частично выполнена
	OrderStatusFilled                                 //выполнена
	OrderStatusCancelled                              //отменена трейдом
	OrderStatusRejected                               //отклонена биржей(несоблюдение правил или недостаточно средств для осуществления сделки)
)

func (status OrderStatus) String() string {
	switch status {
	case OrderStatusNew:
		return "заявка харегестрирована"
	case OrderStatusPartiallyFilled:
		return "заявка частично выполнена"
	case OrderStatusFilled:
		return "заявка выполнена"
	case OrderStatusCancelled:
		return "Трейдер отменил заявку"
	case OrderStatusRejected:
		return "заявка заморожена биржей"
	default:
		return "статус не определяется"
	}
}

type Trader struct {
	ID        uuid.UUID
	Name      string
	Cash      float64
	CreatedAt time.Time
}

type Order struct {
	ID        uuid.UUID
	TraderID  uuid.UUID
	Side      Side
	Type      OrderType
	Symbo     string  //AAPL
	Price     float64 //за единицу товара
	Quantity  float64 //желаемый обьём при совершении сделки
	Filled    float64 //частично исполненный обьём по заказу
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (order *Order) Remaining() float64 {
	return order.Quantity - order.Filled
}

func (order *Order) IsActive() bool {
	return order.Status == OrderStatusNew || order.Status == OrderStatusPartiallyFilled
}

type Trade struct {
	ID          uuid.UUID
	Symbol      string
	BuyOrderID  uuid.UUID
	SellOrderID uuid.UUID
	BuyerID     uuid.UUID
	SellerID    uuid.UUID
	Price       float64
	Quantity    float64
	ExecutedAt  time.Time
}

// position1.price = 0.1
// position2.price = 0.2
// position1.price = 0.1 + position2.price = 0.2 = 0.3
// false
type Position struct {
	TraderID uuid.UUID
	Symbol   string
	Quantity float64
	AvgPrice float64
	Realized float64
}

func RoundPrice(price float64) float64 {
	return math.Round(price*100) / 100
}
