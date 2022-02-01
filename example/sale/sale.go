package sale

import (
	"fmt"
	"gitee.com/super_step/eventBus/example/manage"
	"gitee.com/super_step/eventBus/pkg/memstore"
)

type SaleEvent struct{}

func (s *SaleEvent) Handler(topic string, ctx *memstore.Store, events ...interface{}) error {
	fmt.Printf("topic:%s 销售单事件:%v\n", topic, events)

	return nil
}

func Sale() {
	manage.Sale.PublishSyncNoWait(1, "sale -> sale")
	manage.Order.PublishSyncNoWait("sale -> order")
}
