package eventBus

import (
	"reflect"
)

// 异步发布
func (bus *eventBus) Publish(topic string, events ...interface{}) {
	bus.logger.Debugf("asyncPublish topic:%s events:%v", topic, events)
	bus.wg.Add(1)
	go bus.publish(topic, events...)
	return
}

// 同步发布
func (bus *eventBus) PublishSync(topic string, events ...interface{}) error {
	bus.logger.Debugf("syncPublish topic:%s events:%v", topic, events)
	return bus.publishSync(topic, true, events...)
}

// 同步发布, 不等待异步操作完成
func (bus *eventBus) PublishSyncNoWait(topic string, events ...interface{}) error {
	bus.logger.Debugf("syncPublishNoWait topic:%s events:%v", topic, events)
	return bus.publishSync(topic, false, events...)
}

func (bus *eventBus) publish(topic string, events ...interface{}) {
	_ = bus.publishSync(topic, true, events...)
	bus.wg.Done()
	return
}

func (bus *eventBus) publishSync(topic string, wait bool, events ...interface{}) error {
	Topic := bus.getTopic(topic)
	bus.callAsync(topic, events, Topic)
	err := bus.callSync(topic, events, Topic)
	if err != nil {
		return err
	}
	if wait {
		Topic.wg.Wait()
	}
	return err
}

// 执行同步订阅回调
func (bus *eventBus) callSync(topic string, events []interface{}, Topic *topic) error {
	Topic.RLock()
	syncHandlers := make([]Callback, len(Topic.syncHandlers))
	if len(Topic.syncHandlers) > 0 {
		copy(syncHandlers, Topic.syncHandlers)
	}
	Topic.RUnlock()
	var tmpErr error
	for _, syncHandler := range syncHandlers {
		// fmt.Println(syncHandler)
		err := syncHandler.Callback(topic, events...)
		if err != nil {
			bus.logger.Errorf("eventBus(sync): %s%v#%s", topic, events, err.Error())
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
func (bus *eventBus) callAsync(topic string, events []interface{}, Topic *topic) {
	for _, asyncHandler := range Topic.asyncHandlers.ToSlice() {
		callback, ok := asyncHandler.(reflect.Value)
		if ok {
			Topic.wg.Add(1)
			go func() {
				// 通过反射调用
				params := []reflect.Value{
					reflect.ValueOf(topic),
				}
				for _, event := range events {
					params = append(params, reflect.ValueOf(event))
				}
				callbackFunc := callback.MethodByName("Callback")
				result := callbackFunc.Call(params)
				if len(result) > 0 && !result[0].IsNil() {
					err, ok := result[0].Interface().(error)
					if ok {
						bus.logger.Errorf("eventBus(async): %s%v#%s", topic, events, err.Error())
					}
				}
				Topic.wg.Done()
			}()
		}
	}
}
