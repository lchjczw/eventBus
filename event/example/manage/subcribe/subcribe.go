package subcribe

import (
	"gitee.com/super_step/eventBus/event/example/manage"
	"gitee.com/super_step/eventBus/event/example/order"
	"gitee.com/super_step/eventBus/event/example/sale"
)

func init() {
	manage.Sale.Register(sale.SaleEvent{})
	manage.Order.Register(order.OrderEvent{})
}