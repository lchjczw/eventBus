package eventBus

import (
	"bytes"
	"encoding/gob"
)

type controller struct{}

func (ctr *controller) HasSubscribed(topic string) bool {
	return hasSubscribed(topic)
}

func (ctr *controller) WaitAsync() {
	return
}

func hasSubscribed(topic string) bool {
	return false
}

func getTopic(topicName string) (resTopic topic) {
	res, find := bus.topicMap.LoadOrStore(topicName, resTopic)
	if find && res != nil {
		resTopic = res.(topic)
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
