package eventBus

import (
	"gitee.com/super_step/go_utils/myError"
	"github.com/kataras/iris/v12/core/memstore"
	"reflect"
	"sync"
)

// 异步发布
func (bus *eventBus) Publish(topic string, events ...interface{}) {
	bus.wg.Add(1)
	bus.logger.Debugf("asyncPublish topic:%s events:%v", topic, events)
	go bus.publish(topic, events...)
	return
}

// 同步发布
func (bus *eventBus) PublishSync(topic string, events ...interface{}) error {
	bus.wg.Add(1)
	bus.logger.Debugf("syncPublish topic:%s events:%v", topic, events)
	return bus.publishSync(topic, false, events...)
}

// 同步发布, 不等待异步操作完成
func (bus *eventBus) PublishSyncNoWait(topic string, events ...interface{}) error {
	bus.wg.Add(1)
	bus.logger.Debugf("syncPublishNoWait topic:%s events:%v", topic, events)
	return bus.publishSync(topic, false, events...)
}

func (bus *eventBus) CloseTopic(topic string) {
	bus.topicMap.Delete(topic)
	return
}

func (bus *eventBus) publish(topic string, events ...interface{}) {
	_ = bus.publishSync(topic, true, events...)
	return
}

func (bus *eventBus) publishSync(topic string, wait bool, events ...interface{}) error {
	publishStore := &memstore.Store{}
	wg := &sync.WaitGroup{}
	publishStore.Set("waitGroup", wg)
	Topic := bus.getTopic(topic)
	if Topic.beforeCallback != nil {
		err := Topic.beforeCallback(topic, publishStore, events...)
		if err != nil {
			err = myError.Warp(err, "前置生命周期出错, 流程停止")
			bus.wg.Done()
			return err
		}
	}
	bus.callAsync(topic, publishStore, events, Topic)
	err := bus.callSync(topic, publishStore, events, Topic)
	if Topic.afterSyncCallback != nil {
		Topic.afterSyncCallback(topic, publishStore, events...)
	}
	if wait {
		bus.waitAsync(topic, publishStore, wg, Topic, events...)
	} else {
		go bus.waitAsync(topic, publishStore, wg, Topic, events...)
	}
	if err != nil {
		return err
	}
	return err
}

func (bus *eventBus) waitAsync(topic string, store *memstore.Store, wg *sync.WaitGroup, Topic *topic, events ...interface{}) {
	if Topic.afterCallback != nil {
		Topic.afterCallback(topic, store, events...)
	}
	wg.Wait()
	bus.wg.Done()
	return
}

// 执行同步订阅回调
func (bus *eventBus) callSync(topic string, ctx *memstore.Store, events []interface{}, Topic *topic) error {
	Topic.RLock()
	syncHandlers := make([]Callback, len(Topic.syncHandlers))
	if len(Topic.syncHandlers) > 0 {
		copy(syncHandlers, Topic.syncHandlers)
	}
	Topic.RUnlock()
	var tmpErr error
	for _, syncHandler := range syncHandlers {
		err := syncHandler.Callback(topic, ctx, events...)
		if err != nil {
			bus.logger.Errorf("eventBus(sync): %s%v#%s", topic, events, err.Error())
			if Topic.onErrorCallback != nil {
				Topic.onErrorCallback(topic, ctx, err, events...)
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
func (bus *eventBus) callAsync(topic string, store *memstore.Store, events []interface{}, Topic *topic) {
	for _, asyncHandler := range Topic.asyncHandlers.ToSlice() {
		callback, _ := asyncHandler.(reflect.Value)
		waitGroup, _ := store.Get("waitGroup").(*sync.WaitGroup)
		waitGroup.Add(1)
		go func() {
			// 通过反射调用
			params := []reflect.Value{
				reflect.ValueOf(topic),
				reflect.ValueOf(store),
			}
			for _, event := range events {
				params = append(params, reflect.ValueOf(event))
			}
			callbackFunc := callback.MethodByName("Callback")
			result := callbackFunc.Call(params)
			if len(result) > 0 && !result[0].IsNil() {
				err, ok := result[0].Interface().(error)
				if ok && err != nil {
					bus.logger.Errorf("eventBus(async): %s%v#%s", topic, events, err.Error())
					if Topic.onErrorCallback != nil {
						Topic.onErrorCallback(topic, store, err, events...)
					}
				}
			}
			waitGroup.Done()
		}()
	}
}
