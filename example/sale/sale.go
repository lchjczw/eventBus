package sale

import (
	"fmt"
	"github.com/lchjczw/eventBus/example/manage"
	"github.com/lchjczw/eventBus/pkg/memstore"
)

type SaleEvent struct{}

func (s *SaleEvent) Handler(topic string, ctx *memstore.Store, args ...interface{}) error {
	switch topic {
	case manage.Sale.CompleteTopic():
		fmt.Printf("topic:%s 销售单事件:%v\n", topic, args)
		break
	case manage.Status.CompleteTopic():
		break
	}
	return nil
}

func Sale() {
	manage.Sale.PublishSyncNoWait(1, "sale -> sale")
	manage.Order.PublishSyncNoWait("sale -> order")
}
