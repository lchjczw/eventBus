package internal

import (
	"github.com/lchjczw/eventBus/pkg/memstore"
	"github.com/lchjczw/eventBus/pkg/myError"
	"reflect"
	"sync"
)

// 异步发布
func (bus *eventBus) PublishAsync(topic string, args ...interface{}) {
	bus.wg.Add(1)
	bus.logger.Debugf("asyncPublish topic:%s args:%v", topic, args)
	go bus.publish(topic, args...)
	return
}

// 同步发布
func (bus *eventBus) PublishSync(topic string, args ...interface{}) error {
	bus.wg.Add(1)
	bus.logger.Debugf("syncPublish topic:%s args:%v", topic, args)
	return bus.publishSync(topic, false, args...)
}

// 同步发布, 不等待异步操作完成
func (bus *eventBus) PublishSyncNoWait(topic string, args ...interface{}) error {
	bus.wg.Add(1)
	bus.logger.Debugf("syncPublishNoWait topic:%s args:%v", topic, args)
	return bus.publishSync(topic, false, args...)
}

func (bus *eventBus) CloseTopic(topic string) {
	bus.topicMap.Delete(topic)
	return
}

func (bus *eventBus) publish(topic string, args ...interface{}) {
	_ = bus.publishSync(topic, true, args...)
	return
}

func (bus *eventBus) publishSync(topic string, wait bool, args ...interface{}) error {
	publishStore := &memstore.Store{}
	wg := &sync.WaitGroup{}
	publishStore.Set("waitGroup", wg)
	Topic := bus.getTopic(topic)
	if Topic.hook != nil {
		err := Topic.hook.Before(topic, publishStore, args...)
		if err != nil {
			err = myError.Warp(err, "前置生命周期出错, 流程停止")
			bus.wg.Done()
			return err
		}
	}
	bus.callAsync(topic, publishStore, args, Topic)
	err := bus.callSync(topic, publishStore, args, Topic)
	if Topic.hook != nil {
		Topic.hook.AfterSync(topic, publishStore, args...)
	}
	if wait {
		bus.waitAsync(topic, publishStore, wg, Topic, args...)
	} else {
		go bus.waitAsync(topic, publishStore, wg, Topic, args...)
	}
	if err != nil {
		return err
	}
	return err
}

func (bus *eventBus) waitAsync(topic string, store *memstore.Store, wg *sync.WaitGroup, Topic *topic, args ...interface{}) {
	if Topic.hook != nil {
		Topic.hook.After(topic, store, args...)
	}
	wg.Wait()
	bus.wg.Done()
	return
}

// 执行同步订阅回调
func (bus *eventBus) callSync(topic string, ctx *memstore.Store, args []interface{}, Topic *topic) error {
	Topic.RLock()
	syncHandlers := make([]Handler, len(Topic.syncHandlers))
	if len(Topic.syncHandlers) > 0 {
		copy(syncHandlers, Topic.syncHandlers)
	}
	Topic.RUnlock()
	var tmpErr error
	for _, syncHandler := range syncHandlers {
		err := syncHandler.Handler(topic, ctx, args...)
		if err != nil {
			bus.logger.Errorf("eventBus(sync): %s%v#%s", topic, args, err.Error())
			if Topic.hook != nil {
				Topic.hook.Error(topic, ctx, err, args...)
			}
			if Topic.transaction {
				return err
			} else {
				tmpErr = err
			}
		}
	}
	return tmpErr
}

// 执行异步订阅回调
func (bus *eventBus) callAsync(topic string, store *memstore.Store, args []interface{}, Topic *topic) {
	for _, asyncHandler := range Topic.asyncHandlers.ToSlice() {
		handler, _ := asyncHandler.(reflect.Value)
		waitGroup, _ := store.Get("waitGroup").(*sync.WaitGroup)
		waitGroup.Add(1)
		go func() {
			// 通过反射调用
			params := []reflect.Value{
				reflect.ValueOf(topic),
				reflect.ValueOf(store),
			}
			for _, arg := range args {
				params = append(params, reflect.ValueOf(arg))
			}
			handlerFunc := handler.MethodByName("Handler")
			result := handlerFunc.Call(params)
			if len(result) > 0 && !result[0].IsNil() {
				err, ok := result[0].Interface().(error)
				if ok && err != nil {
					bus.logger.Errorf("eventBus(async): %s%v#%s", topic, args, err.Error())
					if Topic.hook != nil {
						Topic.hook.Error(topic, store, err, args...)
					}
				}
			}
			waitGroup.Done()
		}()
	}
}
