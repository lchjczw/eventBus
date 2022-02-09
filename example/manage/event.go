package manage

import (
	"github.com/lchjczw/eventBus/event"
)

var (
	r1 = event.Root()

	Sale  = r1.Event(`sale`, "销售相关")
	Status = Sale.Event(`status`,`状态相关`)
	Order = r1.Event(`order`, "订单相关")
)
