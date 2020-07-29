package eventBus

import (
	"errors"
	"reflect"
	"sync"
)

// 同步订阅主题
func (bus *eventBus) Subscribe(topic string, callback CallbackFunc) error {
	return bus.subscribe(topic, callback, false)
}

// 异步订阅主题
func (bus *eventBus) SubscribeAsync(topic string, callback CallbackFunc) error {
	return bus.subscribe(topic, callback, true)
}

// 取消订阅
func (bus *eventBus) UnSubscribe(topic string, callback CallbackFunc) {
	Topic := bus.getTopic(topic)
	reflectCallback := reflect.ValueOf(callback)
	if Topic.asyncHandlers.Contains(reflectCallback) {
		Topic.asyncHandlers.Remove(reflectCallback)
		return
	}
	Topic.Lock()
	index := findCallback(Topic.syncHandlers, callback)
	if index > -1 {
		Topic.syncHandlers = removeFromSync(Topic.syncHandlers, index)
	}
	// bus.topicMap.Store(topic, Topic)
	Topic.Unlock()
}

// todo 取消所有订阅
// func (bus *eventBus) UnSubscribeAll(callback CallbackFunc) {
//
// }

// todo 订阅当前所有频道
// func (bus *eventBus) SubscribeAll(callback CallbackFunc) {
//
// }

// 订阅
func (bus *eventBus) subscribe(topic string, callback CallbackFunc, async bool) error {
	Topic := bus.getTopic(topic)
	if checkSub(Topic, callback) {
		return errors.New("不能重复订阅")
	}
	if async {
		reflectCallback := reflect.ValueOf(callback)
		Topic.asyncHandlers.Add(reflectCallback)
	} else {
		Topic.Lock()
		Topic.syncHandlers = append(Topic.syncHandlers, callback)
		Topic.Unlock()
	}
	return nil
}

func checkSub(topic *topic, callback CallbackFunc) bool {
	reflectCallback := reflect.ValueOf(callback)
	if topic.asyncHandlers.Contains(reflectCallback) {
		return true
	}
	return checkSyncSub(topic.syncHandlers, callback, topic.RLocker())
}

func checkSyncSub(handlers []CallbackFunc, callback CallbackFunc, rMutex sync.Locker) bool {
	rMutex.Lock()
	defer rMutex.Unlock()
	if findCallback(handlers, callback) > -1 {
		return true
	}
	return false
}

// 必须在读或写锁保护的情况下调用
func findCallback(handlers []CallbackFunc, callback CallbackFunc) int {
	for i, subFun := range handlers {
		if isSameFunc(subFun, callback) {
			return i
		}
	}
	return -1
}

// 必须在写锁保护的情况下调用
func removeFromSync(slice []CallbackFunc, index int) (result []CallbackFunc) {
	if index == 0 {
		result = slice[1:]
	} else {
		result = append(slice[:index], slice[index+1:]...)
	}
	return
}
