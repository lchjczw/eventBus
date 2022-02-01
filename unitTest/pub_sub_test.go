package unitTest

import (
	"errors"
	"fmt"
	"gitee.com/super_step/eventBus"
	"gitee.com/super_step/go_utils/logger"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12/core/memstore"
	"github.com/kataras/pio"
	"os"
	"testing"
)

const (
	syncTopic  = "syncTopic"
	asyncTopic = "asyncTopic"
)

var (
	testBus   eventBus.EventBus
	callbacks []eventBus.Callback
	testCycle cycle
)

type callback struct {
	Name      string
	Err       bool
	Recursion bool
}

type cycle struct {
	BeforeFlag      bool
	AfterFlag       bool
	AfterSyncFlag   bool
	OnErrorFlag     bool
	StoreCount      int
	BeforeErrorFlag bool
	AfterCount      int
	AfterSyncCount  int
}

func (cycle *cycle) OnBefore(topic string, ctx *memstore.Store, events ...interface{}) error {
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnBefore %s: %+v", topic, events)
	cycle.BeforeFlag = true
	return nil
}

func (cycle *cycle) OnBeforeError(topic string, ctx *memstore.Store, events ...interface{}) error {
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnBeforeError %s: %+v", topic, events)
	cycle.BeforeErrorFlag = true
	return errors.New("OnBeforeError")
}

func (cycle *cycle) OnAfter(topic string, ctx *memstore.Store, events ...interface{}) {
	count, _ := ctx.GetInt("count")
	cycle.StoreCount = count + 1
	cycle.AfterCount++
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnAfter %s: %+v", topic, events)
	cycle.AfterFlag = true
}

func (cycle *cycle) OnAfterSync(topic string, ctx *memstore.Store, events ...interface{}) {
	count, _ := ctx.GetInt("count")
	cycle.StoreCount = count + 1
	cycle.AfterSyncCount++
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("OnAfterSync %s: %+v", topic, events)
	cycle.AfterSyncFlag = true
}

func (cycle *cycle) OnError(topic string, ctx *memstore.Store, err error, events ...interface{}) {
	golog.Default.Infof("OnError %s: %+v %s", topic, events, err.Error())
	cycle.OnErrorFlag = true
}

func (callback *callback) Callback(topic string, ctx *memstore.Store, events ...interface{}) error {
	if len(events) == 0 {
		golog.Default.Infof("%s# %s: %v", callback.Name, topic, "Recursioned")
	} else {
		golog.Default.Infof("%s# %s: %v", callback.Name, topic, events)
		if callback.Recursion {
			testBus.Publish(topic)
		}
	}
	if callback.Err {
		return errors.New(callback.Name)
	}
	return nil
}

func TestSub(t *testing.T) {
	for index, callback := range callbacks {
		err := testBus.Register(syncTopic, callback)
		if err != nil {
			t.Error(index, err.Error())
			return
		}
		err = testBus.Register(syncTopic, callback)
		if err == nil {
			t.Error("重复订阅未报告异常")
			return
		}
	}
}

func TestSubAsync(t *testing.T) {
	for index, callback := range callbacks {
		err := testBus.RegisterAsync(asyncTopic, callback)
		if err != nil {
			t.Error(index, err.Error())
			return
		}
		err = testBus.RegisterAsync(asyncTopic, callback)
		if err == nil {
			t.Error("重复订阅未报告异常")
			return
		}
	}
}

func TestPublish(t *testing.T) {
	testBus.SetTransaction(syncTopic, true)
	testBus.Publish(syncTopic, "publish test", "transaction")
	testBus.WaitAsync()
	testBus.SetTransaction(syncTopic, false)
	testBus.Publish(syncTopic, "publish test")
	testBus.Publish(asyncTopic, "publish test")
	testBus.WaitAsync()
}

func TestPublishSync(t *testing.T) {
	err := testBus.PublishSync(syncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
	err = testBus.PublishSyncNoWait(syncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
	err = testBus.PublishSync(asyncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
	err = testBus.PublishSyncNoWait(asyncTopic, "PublishSync")
	if err != nil {
		golog.Default.Error(err.Error())
	}
}

func TestSetCycleError(t *testing.T) {
	testBus.SetCycleBefore(asyncTopic, testCycle.OnBeforeError)
	testBus.SetCycleBefore(syncTopic, testCycle.OnBeforeError)
	TestPublish(t)
	TestPublishSync(t)
}

func TestCycleError(t *testing.T) {
	if !testCycle.BeforeErrorFlag {
		t.Error("CycleBeforeErrorFlag回调失败")
	}
	if testCycle.AfterCount > 0 {
		t.Errorf("CycleBeforeError逻辑出错: %d", testCycle.StoreCount)
	}
	if testCycle.AfterSyncCount > 0 {
		t.Errorf("CycleBeforeError逻辑出错: %d", testCycle.StoreCount)
	}
}

func TestSetCycle(t *testing.T) {
	testBus.SetCycleBefore(asyncTopic, testCycle.OnBefore)
	testBus.SetCycleAfterAll(asyncTopic, testCycle.OnAfter)
	testBus.SetCycleAfterSync(asyncTopic, testCycle.OnAfterSync)
	testBus.SetCycleError(asyncTopic, testCycle.OnError)
	testBus.SetCycleBefore(syncTopic, testCycle.OnBefore)
	testBus.SetCycleAfterAll(syncTopic, testCycle.OnAfter)
	testBus.SetCycleAfterSync(syncTopic, testCycle.OnAfterSync)
	testBus.SetCycleError(syncTopic, testCycle.OnError)
	TestPublish(t)
	TestPublishSync(t)
}

func TestCycle(t *testing.T) {
	testBus.WaitAsync()
	if !testCycle.BeforeFlag {
		t.Error("CycleBefore回调失败")
	}
	if !testCycle.AfterFlag {
		t.Error("CycleAfter回调失败")
	}
	if !testCycle.AfterSyncFlag {
		t.Error("CycleAfterSync回调失败")
	}
	if !testCycle.OnErrorFlag {
		t.Error("CycleOnError回调失败")
	}
	if testCycle.StoreCount == 0 {
		t.Errorf("Store测试失败: %d", testCycle.StoreCount)
	}
}

func TestUnSubAsync(t *testing.T) {
	for index, callback := range callbacks {
		testBus.Publish(syncTopic, "Sync", "UnSub", index)
		testBus.WaitAsync()
		testBus.UnSubscribe(syncTopic, callback)
	}
}

func TestUnSubSync(t *testing.T) {
	for index, callback := range callbacks {
		testBus.Publish(asyncTopic, "Async", "UnSub", index)
		testBus.WaitAsync()
		testBus.UnSubscribe(asyncTopic, callback)
	}
}

func TestUnSubAll(t *testing.T) {
	TestSub(t)
	TestSubAsync(t)
	length := len(callbacks)
	for i := 1; i <= length; i++ {
		testBus.Publish(asyncTopic, "Async", "UnSubAll", i)
		testBus.Publish(syncTopic, "Sync", "UnSubAll", i)
		testBus.WaitAsync()
		testBus.UnSubscribeAll(callbacks[length-i])
	}

	for index, callback := range callbacks {
		testBus.Publish(asyncTopic, "Async", "UnSubAll", index)
		testBus.Publish(syncTopic, "Sync", "UnSubAll", index)
		testBus.WaitAsync()
		testBus.UnSubscribeAll(callback)
	}
}

func TestRecursionSub(t *testing.T) {
	callbackRecursion := callback{
		Name:      "Recursion",
		Err:       false,
		Recursion: true,
	}
	err := testBus.Register(syncTopic, &callbackRecursion)
	if err != nil {
		t.Error("Recursion", err.Error())
		return
	}
	testBus.Publish(syncTopic, "Recursion")
	testBus.WaitAsync()
	testBus.UnSubscribe(syncTopic, &callbackRecursion)
}

func TestCloseTopic(t *testing.T) {
	TestSub(t)
	TestSubAsync(t)
	testBus.Publish(syncTopic, "BeforeClose")
	testBus.WaitAsync()
	testBus.CloseTopic(syncTopic)
	testBus.Publish(syncTopic, "Closed")
	testBus.WaitAsync()
}

func TestMain(m *testing.M) {
	golog.Default.Printer = pio.NewPrinter("", os.Stdout).
		EnableDirectOutput().Hijack(logger.Hijacker).SetSync(false)
	golog.Default.SetTimeFormat("2006-01-02 15:04:05")
	golog.Default.SetLevel("debug")
	testBus = eventBus.NewBus(golog.Default)
	testCycle = cycle{}
	for i := 0; i < 4; i++ {
		callbackStruck := newCallback(fmt.Sprintf("%d", i), i == 2)
		callbacks = append(callbacks, callbackStruck)
	}
	m.Run()
	testBus.WaitAsync()
}

func newCallback(name string, err bool) eventBus.Callback {
	callbackStruck := callback{
		Name: name,
		Err:  err,
	}
	return &callbackStruck
}
