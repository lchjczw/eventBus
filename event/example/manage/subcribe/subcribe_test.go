package subcribe

import (
	"gitee.com/super_step/eventBus/event/example/order"
	"gitee.com/super_step/eventBus/event/example/sale"
	"testing"
)

func TestEvent(t *testing.T) {
	order.Order()
	sale.Sale()
}
