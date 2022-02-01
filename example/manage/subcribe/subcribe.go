package subcribe

import (
	"gitee.com/super_step/eventBus/example/manage"
	"gitee.com/super_step/eventBus/example/order"
	"gitee.com/super_step/eventBus/example/sale"
)

func init() {
	manage.Sale.SubscribeSync(&sale.SaleEvent{})
	manage.Order.SubscribeSync(&order.OrderEvent{})
}
