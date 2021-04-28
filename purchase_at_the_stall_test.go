package fsm

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type StallPurchaseOrder struct {
	ID              uint64
	EventCount      uint64
	PreEvent        string   //前置时间
	GetGoodsResults int64    //0:拿货成功 1：下架 2：后天 3：明天
	priority        int64    //订单紧急状况优先级：0:普通 1：紧急
	effective       int      //
	State           string   //当前状态
	States          []string //
}

// StallPurchaseOrderEventProcessor is used to handle StallPurchaseOrder actions.
type StallPurchaseOrderEventProcessor struct{}

func (p *StallPurchaseOrderEventProcessor) OnExit(fromState string, args []interface{}) {
	t := args[0].(*StallPurchaseOrder)
	if t.State != fromState {
		panic(fmt.Errorf("采购单 %v 的状态与期望的状态 %s 不一致，可能在状态机外被改变了", t, fromState))
	}

	log.Printf("采购单 %d 从状态 %s 改变", t.ID, fromState)
}

func (p *StallPurchaseOrderEventProcessor) Action(action string, fromState string, toState string, args []interface{}) error {
	ts := args[0].(*StallPurchaseOrder)
	ts.EventCount++

	switch action {
	case "获取拿货结果":
		{
			ts.GetGoodsResults = int64(rand.Intn(4))
		}
	case "禁用此货源":
		{
			fmt.Println("禁用此货源")
		}

	case "生成络采购单":
		{
			//fmt.Println("生成当日网络采购单")
		}

	case "生成次日网络采购单":
		{
			//fmt.Println("生成次日网络采购单")
		}

	case "此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效":

		{
			fmt.Println("此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效")
		}
	case "检测拿货时效":
		{
			ts.effective = rand.Intn(2)
		}

	case "送至仓库":

	case "收货":

	default: //其它action

	}

	return nil
}

func createNetOrder(ts *StallPurchaseOrder, fsm *StateMachine) {
	event := "紧急需要"
	switch ts.priority {
	case 0:
		event = "不紧急需要"
		fmt.Println("生成次日网络采购单")
	case 1:
		fmt.Println("生成当日网络采购单")
		event = "紧急需要"
	}

	err := fsm.Trigger(ts.State, event, ts, fsm)
	if err != nil {
		fmt.Printf("trigger err: %v", err)
	}

}

func (p *StallPurchaseOrderEventProcessor) OnEnter(toState string, args []interface{}) {
	t := args[0].(*StallPurchaseOrder)
	t.State = toState
	t.States = append(t.States, toState)

	log.Printf("采购单 %d 的状态改变为 %s ", t.ID, toState)
}

func (p *StallPurchaseOrderEventProcessor) OnActionFailure(action string, fromState string, toState string, args []interface{}, err error) {
	t := args[0].(*StallPurchaseOrder)

	log.Printf("采购单 %d 的状态从 %s to %s 改变失败， 原因: %v", t.ID, fromState, toState, err)
}

func TestStallOrder(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	ts := &StallPurchaseOrder{
		ID:     1,
		State:  "等待采购",
		States: []string{"等待采购"},
	}

	ts.priority = int64(rand.Intn(2))
	fsm := initStallOrderFSM()

	err := fsm.Trigger(ts.State, "去档口拿货", ts, fsm)
	if err != nil {
		t.Errorf("trigger err: %v", err)
	}

	event := "拿货结果(拿货成功)"
	switch ts.GetGoodsResults {
	case 0:
		event = "拿货结果(拿货成功)"
	case 1:
		event = "拿货结果(该货物已经下架)"
	case 2:
		event = "拿货结果(后天拿货)"
	case 3:
		event = "拿货结果(明天拿货)"
	}

	err = fsm.Trigger(ts.State, event, ts, fsm)
	ts.PreEvent = event
	if err != nil {
		fmt.Printf("trigger err: %v", err)
	}

	switch ts.GetGoodsResults {
	case 0:

		err = fsm.Trigger(ts.State, "送至仓库", ts, fsm)
		ts.PreEvent = "送至仓库"
		if err != nil {
			t.Errorf("trigger err: %v", err)
		}

		err = fsm.Trigger(ts.State, "收货", ts, fsm)
		if err != nil {
			t.Errorf("trigger err: %v", err)
		}

	case 1:
		createNetOrder(ts, fsm)
	case 2:
		createNetOrder(ts, fsm)
	case 3:
		event := "未超过有效期"
		switch ts.effective {
		case 0:
			event = "未超过有效期"

		case 1:
			event = "超过有效期"
		}

		err := fsm.Trigger(ts.State, event, ts, fsm)
		ts.PreEvent = event
		if err != nil {
			fmt.Printf("trigger err: %v", err)
		}

		if ts.effective == 1 {
			createNetOrder(ts, fsm)
		}

	}

	//fsm.Export("state.png")

}

func compareStallPurchaseOrder(t1 *StallPurchaseOrder, t2 *StallPurchaseOrder) bool {
	if t1.ID != t2.ID || t1.EventCount != t2.EventCount ||
		t1.State != t2.State {
		return false
	}

	return fmt.Sprint(t1.States) == fmt.Sprint(t2.States)
}

func initStallOrderFSM() *StateMachine {
	delegate := &DefaultDelegate{P: &StallPurchaseOrderEventProcessor{}}

	transitions := []Transition{
		{From: "等待采购", Event: "去档口拿货", To: "采购中", Action: "获取拿货结果"},
		{From: "采购中", Event: "拿货结果(拿货成功)", To: "采购成功", Action: "发货"},
		{From: "采购成功", Event: "送至仓库", To: "采购成功", Action: ""},
		{From: "采购成功", Event: "收货", To: "结束", Action: "订单完结"},
		{From: "采购中", Event: "拿货结果(该货物已经下架)", To: "采购失败", Action: "禁用此货源"},
		{From: "采购失败", Event: "紧急需要", To: "结束", Action: "生成当日网络采购单"},
		{From: "采购失败", Event: "不紧急需要", To: "结束", Action: "生成次日网络采购单"},
		{From: "采购中", Event: "拿货结果(后天拿货)", To: "采购失败", Action: "此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效"},
		{From: "采购中", Event: "拿货结果(明天拿货)", To: "采购中", Action: "检测拿货时效"},
		{From: "采购中", Event: "超过有效期", To: "采购失败", Action: "创建网采单子"},
		{From: "采购中", Event: "未超过有效期", To: "等待采购", Action: "去档口"},
	}

	return NewStateMachine(delegate, transitions...)
}
func IsExistItem(value interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}
