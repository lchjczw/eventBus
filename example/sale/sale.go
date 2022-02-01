package sale

import (
	"fmt"
	"github.com/lchjczw/eventBus/example/manage"
	"github.com/lchjczw/eventBus/pkg/memstore"
)

type SaleEvent struct{}

func (s *SaleEvent) Handler(topic string, ctx *memstore.Store, args ...interface{}) error {
	fmt.Printf("topic:%s 销售单事件:%v\n", topic, args)

	return nil
}

func Sale() {
	manage.Sale.PublishSyncNoWait(1, "sale -> sale")
	manage.Order.PublishSyncNoWait("sale -> order")
}
