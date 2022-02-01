package event

import (
	"errors"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12/core/memstore"
)


type cycle struct {
	BeforeFlag      bool
	AfterFlag       bool
	AfterSyncFlag   bool
	OnErrorFlag     bool
	StoreCount      int
	BeforeErrorFlag bool
	AfterCount      int
	AfterSyncCount  int
}

func (cycle *cycle) OnBefore(topic string, ctx *memstore.Store, events ...interface{}) error {
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnBefore %s: %+v", topic, events)
	cycle.BeforeFlag = true
	return nil
}

func (cycle *cycle) OnBeforeError(topic string, ctx *memstore.Store, events ...interface{}) error {
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnBeforeError %s: %+v", topic, events)
	cycle.BeforeErrorFlag = true
	return errors.New("OnBeforeError")
}

func (cycle *cycle) OnAfter(topic string, ctx *memstore.Store, events ...interface{}) {
	count, _ := ctx.GetInt("count")
	cycle.StoreCount = count + 1
	cycle.AfterCount++
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnAfter %s: %+v", topic, events)
	cycle.AfterFlag = true
}

func (cycle *cycle) OnAfterSync(topic string, ctx *memstore.Store, events ...interface{}) {
	count, _ := ctx.GetInt("count")
	cycle.StoreCount = count + 1
	cycle.AfterSyncCount++
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnAfterSync %s: %+v", topic, events)
	cycle.AfterSyncFlag = true
}

func (cycle *cycle) OnError(topic string, ctx *memstore.Store, err error, events ...interface{}) {
	golog.Default.Infof("OnError %s: %+v %s", topic, events, err.Error())
	cycle.OnErrorFlag = true
}
