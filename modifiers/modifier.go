package modifiers

import (
	"financialExchange/enum"
)

// возвращает функцию-модификатор для обновления на покупку.
func ModifierForOrderBuy(sellerQuantity float64) func(enum.Order) enum.Order {
	return func(order enum.Order) enum.Order {
		// Увеличиваем уже исполненный объем ордера на величину текущей сделки
		order.Filled += sellerQuantity

		// Проверяем полностью ли исполнен ордер
		if order.Filled >= order.Quantity {
			// Если объем достиг или превысил целевой, помечаем как полностью исполненный
			order.Status = enum.OrderStatusFilled
		} else {
			order.Status = enum.OrderStatusPartiallyFilled
		}

		// Возвращаем обновленную копию структуры ордера
		return order
	}
}
