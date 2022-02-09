package order

import (
	"fmt"
	"github.com/lchjczw/eventBus/example/manage"
	"github.com/lchjczw/eventBus/pkg/memstore"
)

type OrderEvent struct{}

func (e *OrderEvent) Before(topic string, ctx *memstore.Store, args ...interface{}) error {
	fmt.Println("before:", topic)
	//return errors.New(fmt.Sprintf("before:检查发生错误"))
	return nil
}

func (e *OrderEvent) AfterSync(topic string, ctx *memstore.Store, args ...interface{}) {
}

func (e *OrderEvent) After(topic string, ctx *memstore.Store, args ...interface{}) {
}

func (e *OrderEvent) Error(topic string, ctx *memstore.Store, err error, args ...interface{}) {
}

func (e *OrderEvent) Handler(topic string, ctx *memstore.Store, args ...interface{}) error {
	fmt.Printf("topic:%s 订单事件:%v\n", topic, args)
	return nil
}

func Order() {
	err := manage.Sale.PublishSyncNoWait(1, "order -> sale")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = manage.Order.PublishSyncNoWait("order -> order")
	if err != nil {
		fmt.Println(err.Error())
	}
}
