package subcribe

import (
	"github.com/lchjczw/eventBus/example/manage"
	"github.com/lchjczw/eventBus/example/order"
	"github.com/lchjczw/eventBus/example/sale"
)

func init() {
	manage.Sale.SubscribeSync(&sale.SaleEvent{})
	manage.Order.SubscribeSync(&order.OrderEvent{})
}
