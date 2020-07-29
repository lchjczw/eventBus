package eventBus

import (
	mapSet "github.com/deckarep/golang-set"
	"reflect"
	"sync"
)

// func (bus *eventBus) HasSubscribed(topic string) bool {
// 	return bus.hasSubscribed(topic)
// }

func (bus *eventBus) SetTransaction(topicName string, tr bool) {
	Topic := bus.getTopic(topicName)
	Topic.transaction = tr
	bus.topicMap.Store(topicName, Topic)
}

func (bus *eventBus) WaitAsync() {
	bus.wg.Wait()
}

// func (bus *eventBus) hasSubscribed(topic string) bool {
// 	Topic := bus.getTopic(topic)
// 	if Topic.asyncHandlers.Cardinality() > 0 {
// 		return true
// 	}
// 	Topic.RLock()
// 	defer Topic.RUnlock()
// 	if len(Topic.syncHandlers) > 0 {
// 		return true
// 	}
// 	return false
// }

func (bus *eventBus) getTopic(topicName string) (resTopic *topic) {
	res, ok := bus.topicMap.Load(topicName)
	// res, find := bus.topicMap.LoadOrStore(topicName, resTopic)
	if ok && res != nil {
		resTopic = res.(*topic)
	} else {
		resTopic = &topic{
			syncHandlers:  make([]CallbackFunc, 0),
			asyncHandlers: mapSet.NewSet(),
			transaction:   false,
			RWMutex:       sync.RWMutex{},
		}
		bus.topicMap.Store(topicName, resTopic)
	}
	return
}

func isSameFunc(a, b CallbackFunc) bool {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	return sf1.Pointer() == sf2.Pointer()
}
