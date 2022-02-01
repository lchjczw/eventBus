package eventBus

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// 同步订阅主题
func (bus *eventBus) SubscribeSync(topic string, callback Callback) error {
	return bus.subscribe(topic, callback, false)
}

// 异步订阅主题
func (bus *eventBus) SubscribeAsync(topic string, callback Callback) error {
	return bus.subscribe(topic, callback, true)
}

// 取消订阅
func (bus *eventBus) UnSubscribe(topic string, callback Callback) {
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

// 取消所有订阅
func (bus *eventBus) UnSubscribeAll(callback Callback) {
	bus.topicMap.Range(func(topic, Topic interface{}) bool {
		bus.UnSubscribe(topic.(string), callback)
		return true
	})
}

// 订阅
func (bus *eventBus) subscribe(topic string, callback Callback, async bool) error {
	Topic := bus.getTopic(topic)
	if checkSub(Topic, callback) {
		return errors.New(fmt.Sprintf("topic:%s 重复订阅 ", topic))
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

func checkSub(topic *topic, callback Callback) bool {
	reflectCallback := reflect.ValueOf(callback)
	if topic.asyncHandlers.Contains(reflectCallback) {
		return true
	}
	return checkSyncSub(topic.syncHandlers, callback, topic.RLocker())
}

func checkSyncSub(handlers []Callback, callback Callback, rMutex sync.Locker) bool {
	rMutex.Lock()
	defer rMutex.Unlock()
	if findCallback(handlers, callback) > -1 {
		return true
	}
	return false
}

// 必须在读或写锁保护的情况下调用
func findCallback(handlers []Callback, callback Callback) int {
	for i, subFun := range handlers {
		if isSameFunc(subFun, callback) {
			return i
		}
	}
	return -1
}

// 必须在写锁保护的情况下调用
func removeFromSync(slice []Callback, index int) (result []Callback) {
	result = append(slice[:index], slice[index+1:]...)
	return
}
