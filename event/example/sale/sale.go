package sale

import (
	"fmt"
	"gitee.com/super_step/eventBus/event/example/manage"
	"github.com/kataras/iris/v12/core/memstore"
)

type SaleEvent struct{}

func (s *SaleEvent) Callback(topic string, ctx *memstore.Store, events ...interface{}) error {
	fmt.Printf("topic:%s 销售单事件:%v\n", topic, events)

	return nil
}

func Sale() {
	manage.Sale.PublishSyncNoWait(1, "sale -> sale")
	manage.Order.PublishSyncNoWait("sale -> order")
}
