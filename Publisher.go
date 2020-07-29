package eventBus

import (
	"reflect"
)

// 异步发布
func (bus *eventBus) Publish(topic string, events ...interface{}) {
	bus.wg.Add(1)
	go bus.publish(topic, events...)
	return
}

// 同步发布
func (bus *eventBus) PublishSync(topic string, events ...interface{}) {
	bus.publishSync(topic, events...)
	return
}

func (bus *eventBus) publish(topic string, events ...interface{}) {
	bus.publishSync(topic, events...)
	bus.wg.Done()
	return
}

func (bus *eventBus) publishSync(topic string, events ...interface{}) {
	Topic := bus.getTopic(topic)
	bus.callAsync(topic, events, Topic)
	bus.callSync(topic, events, Topic)
	Topic.wg.Wait()
}

// 执行同步订阅回调
func (bus *eventBus) callSync(topic string, events []interface{}, Topic *topic) {
	Topic.RLock()
	syncHandlers := make([]Callback, len(Topic.syncHandlers))
	if len(Topic.syncHandlers) > 0 {
		copy(syncHandlers, Topic.syncHandlers)
	}
	Topic.RUnlock()
	for _, syncHandler := range syncHandlers {
		// fmt.Println(syncHandler)
		err := syncHandler.Callback(topic, events...)
		if err != nil {
			if Topic.transaction {
				return
			}
		}
	}
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
				callbackFunc.Call(params)
				Topic.wg.Done()
			}()
		}
	}
}
