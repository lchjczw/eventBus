package core

import (
	mapSet "github.com/deckarep/golang-set"
	"reflect"
	"sync"
)

func (bus *eventBus) SetTransaction(topicName string, tr bool) {
	Topic := bus.getTopic(topicName)
	Topic.transaction = tr
	bus.topicMap.Store(topicName, Topic)
}
func (bus *eventBus) SetHook(topic string, hook Hook) {
	if hook == nil {
		return
	}
	Topic := bus.getTopic(topic)
	Topic.hook = hook
	bus.topicMap.Store(topic, Topic)
}

func (bus *eventBus) WaitAsync() {
	bus.wg.Wait()
}

func (bus *eventBus) getTopic(topicName string) (resTopic *topic) {
	res, ok := bus.topicMap.Load(topicName)
	// res, find := bus.topicMap.LoadOrStore(topicName, resTopic)
	if ok && res != nil {
		resTopic = res.(*topic)
	} else {
		resTopic = &topic{
			syncHandlers:  make([]Handler, 0),
			asyncHandlers: mapSet.NewSet(),
			transaction:   false,
			RWMutex:       sync.RWMutex{},
		}
		bus.topicMap.Store(topicName, resTopic)
	}
	return
}

func isSameFunc(a, b Handler) bool {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	return sf1.Pointer() == sf2.Pointer()
}
