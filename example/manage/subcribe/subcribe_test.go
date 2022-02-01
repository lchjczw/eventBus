package subcribe

import (
	"gitee.com/super_step/eventBus/example/order"
	"gitee.com/super_step/eventBus/example/sale"
	"testing"
)

func TestEvent(t *testing.T) {
	order.Order()
	sale.Sale()
}
