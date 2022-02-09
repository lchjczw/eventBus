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

func BenchmarkSale(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sale.Sale()
	}

}

func BenchmarkOrder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		order.Order()
	}

}
