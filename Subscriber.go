package eventBus

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// 同步订阅主题
func (bus *eventBus) SubscribeSync(topic string, handler Handler) error {
	return bus.subscribe(topic, handler, false)
}

// 异步订阅主题
func (bus *eventBus) SubscribeAsync(topic string, handler Handler) error {
	return bus.subscribe(topic, handler, true)
}

// 取消订阅
func (bus *eventBus) UnSubscribe(topic string, handler Handler) {
	Topic := bus.getTopic(topic)
	reflectHandler := reflect.ValueOf(handler)
	if Topic.asyncHandlers.Contains(reflectHandler) {
		Topic.asyncHandlers.Remove(reflectHandler)
		return
	}
	Topic.Lock()
	index := findHandler(Topic.syncHandlers, handler)
	if index > -1 {
		Topic.syncHandlers = removeFromSync(Topic.syncHandlers, index)
	}
	// bus.topicMap.Store(topic, Topic)
	Topic.Unlock()
}

// 取消所有订阅
func (bus *eventBus) UnSubscribeAll(handler Handler) {
	bus.topicMap.Range(func(topic, Topic interface{}) bool {
		bus.UnSubscribe(topic.(string), handler)
		return true
	})
}

// 订阅
func (bus *eventBus) subscribe(topic string, handler Handler, async bool) error {
	Topic := bus.getTopic(topic)
	if checkSub(Topic, handler) {
		return errors.New(fmt.Sprintf("topic:%s 重复订阅 ", topic))
	}
	if async {
		reflectHandler := reflect.ValueOf(handler)
		Topic.asyncHandlers.Add(reflectHandler)
	} else {
		Topic.Lock()
		Topic.syncHandlers = append(Topic.syncHandlers, handler)
		Topic.Unlock()
	}
	return nil
}

func checkSub(topic *topic, handler Handler) bool {
	reflectHandler := reflect.ValueOf(handler)
	if topic.asyncHandlers.Contains(reflectHandler) {
		return true
	}
	return checkSyncSub(topic.syncHandlers, handler, topic.RLocker())
}

func checkSyncSub(handlers []Handler, handler Handler, rMutex sync.Locker) bool {
	rMutex.Lock()
	defer rMutex.Unlock()
	if findHandler(handlers, handler) > -1 {
		return true
	}
	return false
}

// 必须在读或写锁保护的情况下调用
func findHandler(handlers []Handler, handler Handler) int {
	for i, subFun := range handlers {
		if isSameFunc(subFun, handler) {
			return i
		}
	}
	return -1
}

// 必须在写锁保护的情况下调用
func removeFromSync(slice []Handler, index int) (result []Handler) {
	result = append(slice[:index], slice[index+1:]...)
	return
}
