package order

import (
	"fmt"
	"gitee.com/super_step/eventBus/event/example/manage"
	"github.com/kataras/iris/v12/core/memstore"
)

type OrderEvent struct{}

func (OrderEvent) Callback(topic string, ctx *memstore.Store, events ...interface{}) error {
	fmt.Printf("topic:%s 订单事件:%v", topic, events)
	return nil
}

func Order() {
	manage.Sale.PublishSyncNoWait(1, "2314")
	manage.Order.PublishSyncNoWait("lchjczw")
}
