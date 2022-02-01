package unitTest

import (
	"errors"
	"fmt"
	"gitee.com/super_step/eventBus"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12/core/memstore"
	"testing"
)

const (
	syncTopic  = "syncTopic"
	asyncTopic = "asyncTopic"
)

var (
	testBus   eventBus.EventBus
	callbacks []eventBus.Callback
	testHook  MyHook
)

type callback struct {
	Name      string
	Err       bool
	Recursion bool
}

type MyHook struct {
	BeforeFlag      bool
	AfterFlag       bool
	AfterSyncFlag   bool
	OnErrorFlag     bool
	StoreCount      int
	AfterCount      int
	AfterSyncCount  int
}

func (cycle *MyHook) Before(topic string, ctx *memstore.Store, args ...interface{}) error {
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("Before %s: %+v", topic, args)
	cycle.BeforeFlag = true
	return nil
}

func (cycle *MyHook) After(topic string, ctx *memstore.Store, args ...interface{}) {
	count, _ := ctx.GetInt("count")
	cycle.StoreCount = count + 1
	cycle.AfterCount++
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("After %s: %+v", topic, args)
	cycle.AfterFlag = true
}

func (cycle *MyHook) AfterSync(topic string, ctx *memstore.Store, args ...interface{}) {
	count, _ := ctx.GetInt("count")
	cycle.StoreCount = count + 1
	cycle.AfterSyncCount++
	ctx.Set("count", cycle.StoreCount)
	golog.Default.Infof("AfterSync %s: %+v", topic, args)
	cycle.AfterSyncFlag = true
}

func (cycle *MyHook) Error(topic string, ctx *memstore.Store, err error, args ...interface{}) {
	golog.Default.Infof("Error %s: %+v %s", topic, args, err.Error())
	cycle.OnErrorFlag = true
}

func (callback *callback) Callback(topic string, ctx *memstore.Store, args ...interface{}) error {
	if len(args) == 0 {
		golog.Default.Infof("%s# %s: %v", callback.Name, topic, "Recursioned")
	} else {
		golog.Default.Infof("%s# %s: %v", callback.Name, topic, args)
		if callback.Recursion {
			testBus.PublishAsync(topic)
		}
	}
	if callback.Err {
		return errors.New(callback.Name)
	}
	return nil
}

func TestSub(t *testing.T) {
	for index, callback := range callbacks {
		err := testBus.SubscribeSync(syncTopic, callback)
		if err != nil {
			t.Error(index, err.Error())
			return
		}
		err = testBus.SubscribeSync(syncTopic, callback)
		if err == nil {
			t.Error("重复订阅未报告异常")
			return
		}
	}
}

func TestSubAsync(t *testing.T) {
	for index, callback := range callbacks {
		err := testBus.SubscribeAsync(asyncTopic, callback)
		if err != nil {
			t.Error(index, err.Error())
			return
		}
		err = testBus.SubscribeAsync(asyncTopic, callback)
		if err == nil {
			t.Error("重复订阅未报告异常")
			return
		}
	}
}

func TestPublish(t *testing.T) {
	testBus.SetTransaction(syncTopic, true)
	testBus.PublishAsync(syncTopic, "publish test", "transaction")
	testBus.WaitAsync()
	testBus.SetTransaction(syncTopic, false)
	testBus.PublishAsync(syncTopic, "publish test")
	testBus.PublishAsync(asyncTopic, "publish test")
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
	testBus.SetHook(asyncTopic, &testHook)
	testBus.SetHook(syncTopic, &testHook)

	TestPublish(t)
	TestPublishSync(t)
}

func TestCycleError(t *testing.T) {
	if testHook.AfterCount > 0 {
		t.Errorf("CycleBeforeError逻辑出错: %d", testHook.StoreCount)
	}
	if testHook.AfterSyncCount > 0 {
		t.Errorf("CycleBeforeError逻辑出错: %d", testHook.StoreCount)
	}
}

func TestSetCycle(t *testing.T) {
	testBus.SetHook(asyncTopic, &testHook)
	testBus.SetHook(syncTopic, &testHook)
	TestPublish(t)
	TestPublishSync(t)
}

func TestCycle(t *testing.T) {
	testBus.WaitAsync()
	if !testHook.BeforeFlag {
		t.Error("CycleBefore回调失败")
	}
	if !testHook.AfterFlag {
		t.Error("CycleAfter回调失败")
	}
	if !testHook.AfterSyncFlag {
		t.Error("CycleAfterSync回调失败")
	}
	if !testHook.OnErrorFlag {
		t.Error("CycleOnError回调失败")
	}
	if testHook.StoreCount == 0 {
		t.Errorf("Store测试失败: %d", testHook.StoreCount)
	}
}

func TestUnSubAsync(t *testing.T) {
	for index, callback := range callbacks {
		testBus.PublishAsync(syncTopic, "Sync", "UnSub", index)
		testBus.WaitAsync()
		testBus.UnSubscribe(syncTopic, callback)
	}
}

func TestUnSubSync(t *testing.T) {
	for index, callback := range callbacks {
		testBus.PublishAsync(asyncTopic, "Async", "UnSub", index)
		testBus.WaitAsync()
		testBus.UnSubscribe(asyncTopic, callback)
	}
}

func TestUnSubAll(t *testing.T) {
	TestSub(t)
	TestSubAsync(t)
	length := len(callbacks)
	for i := 1; i <= length; i++ {
		testBus.PublishAsync(asyncTopic, "Async", "UnSubAll", i)
		testBus.PublishAsync(syncTopic, "Sync", "UnSubAll", i)
		testBus.WaitAsync()
		testBus.UnSubscribeAll(callbacks[length-i])
	}

	for index, callback := range callbacks {
		testBus.PublishAsync(asyncTopic, "Async", "UnSubAll", index)
		testBus.PublishAsync(syncTopic, "Sync", "UnSubAll", index)
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
	err := testBus.SubscribeSync(syncTopic, &callbackRecursion)
	if err != nil {
		t.Error("Recursion", err.Error())
		return
	}
	testBus.PublishAsync(syncTopic, "Recursion")
	testBus.WaitAsync()
	testBus.UnSubscribe(syncTopic, &callbackRecursion)
}

func TestCloseTopic(t *testing.T) {
	TestSub(t)
	TestSubAsync(t)
	testBus.PublishAsync(syncTopic, "BeforeClose")
	testBus.WaitAsync()
	testBus.CloseTopic(syncTopic)
	testBus.PublishAsync(syncTopic, "Closed")
	testBus.WaitAsync()
}

func TestMain(m *testing.M) {
	//golog.Default.Printer = pio.NewPrinter("", os.Stdout).
	//	EnableDirectOutput().Hijack(logger.Hijacker).SetSync(false)
	golog.Default.SetTimeFormat("2006-01-02 15:04:05")
	golog.Default.SetLevel("debug")
	testBus = eventBus.NewBus(golog.Default)
	testHook = MyHook{}
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
