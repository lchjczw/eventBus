package eventBus

import (
	"errors"
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
	Topic := getTopic(topic)
	if Topic.asyncHandlers.Contains(callback) {
		Topic.asyncHandlers.Remove(callback)
		return
	}
	Topic.Lock()
	index := findCallback(Topic.syncHandlers, callback)
	if index > -1 {
		Topic.syncHandlers = removeFromSync(Topic.syncHandlers, index)
	}
	Topic.Unlock()
}

// 订阅
func (bus *eventBus) subscribe(topic string, callback CallbackFunc, async bool) error {
	Topic := getTopic(topic)
	if checkSub(Topic, callback) {
		return errors.New("不能重复订阅")
	}
	if async {
		Topic.asyncHandlers.Add(callback)
	} else {
		Topic.Lock()
		for _, subFun := range Topic.syncHandlers {
			if isSameFunc(subFun, callback) {
				break
			}
		}
		Topic.syncHandlers = append(Topic.syncHandlers, callback)
		Topic.Unlock()
	}
	return nil
}

func checkSub(topic *topic, callback CallbackFunc) bool {
	if topic.asyncHandlers.Contains(callback) {
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
	sliceLength := len(slice)
	if index == 0 {
		result = slice[index+1:]
	} else if index == sliceLength {
		result = slice[:index]
	} else {
		result = append(slice[:index], slice[index+1:]...)
	}
	return
}
