package eventBus

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

func (bus *eventBus) HasSubscribed(topic string) bool {
	return hasSubscribed(topic)
}

func (bus *eventBus) SetTransaction(topicName string, tr bool) bool {
	res, ok := bus.topicMap.Load(topicName)
	if ok {
		Topic := res.(*topic)
		Topic.transaction = tr
	}
	return ok
}

func (bus *eventBus) WaitAsync() {
	bus.wg.Wait()
}

func hasSubscribed(topic string) bool {
	Topic := getTopic(topic)
	if Topic.asyncHandlers.Cardinality() > 0 {
		return true
	}
	Topic.RLock()
	defer Topic.RUnlock()
	if len(Topic.syncHandlers) > 0 {
		return true
	}
	return false
}

func getTopic(topicName string) (resTopic *topic) {
	res, find := bus.topicMap.LoadOrStore(topicName, resTopic)
	if find && res != nil {
		resTopic = res.(*topic)
	}
	return
}

func clone(a, b interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(a); err != nil {
		return err
	}
	if err := dec.Decode(b); err != nil {
		return err
	}
	return nil
}

func isSameFunc(a, b CallbackFunc) bool {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	return sf1.Pointer() == sf2.Pointer()
}
