package subcribe

import (
	"github.com/lchjczw/eventBus/example/order"
	"github.com/lchjczw/eventBus/example/sale"
	"testing"
)

func TestEvent(t *testing.T) {
	order.Order()
	sale.Sale()
}
