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
	PreEvent        Event   //前置事件
	GetGoodsResults int64   //0:拿货成功 1：下架 2：后天 3：明天
	priority        int64   //订单紧急状况优先级：0:普通 1：紧急
	effective       int     //
	CurState        State   //当前状态
	States          []State //
}

// StallPurchaseOrderEventProcessor is used to handle StallPurchaseOrder actions.
type StallPurchaseOrderEventProcessor struct{}

func (p *StallPurchaseOrderEventProcessor) OnExit(fromState State, args []interface{}) {
	t := args[0].(*StallPurchaseOrder)
	if t.CurState != fromState {
		panic(fmt.Errorf("采购单 %v 的状态与期望的状态 %v 不一致，可能在状态机外被改变了", t, fromState))
	}

	log.Printf("采购单 %d 从状态 %v 改变", t.ID, fromState)
}

func (p *StallPurchaseOrderEventProcessor) PreAction(action Action, fromState State, toState State, args []interface{}) error {
	ts := args[0].(*StallPurchaseOrder)
	ts.EventCount++

	switch action.Name {

	case "获取拿货结果":
		{
			ts.GetGoodsResults = int64(rand.Intn(4))
		}
	case "禁用此货源":
		{
			fmt.Println("禁用此货源")
		}

	case "生成络采购单", "生成次日网络采购单":
		//{
		//	fmt.Println("生成网络采购单")
		//}

	case "此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效":

		{
			fmt.Println("此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效")
		}
	case "检测拿货时效":
		{
			ts.effective = rand.Intn(2)
		}
	default: //其它action

	}

	return nil
}

func (p *StallPurchaseOrderEventProcessor) NextAction(action Action, fromState State, toState State, args []interface{}) error {
	ts := args[0].(*StallPurchaseOrder)
	ts.EventCount++

	switch action.Name {

	case "获取拿货结果":
		{
			ts.GetGoodsResults = int64(rand.Intn(4))
		}
	case "禁用此货源":
		{
			fmt.Println("禁用此货源")
		}

	case "生成络采购单", "生成次日网络采购单":
		//{
		//	fmt.Println("生成网络采购单")
		//}

	case "此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效":

		{
			fmt.Println("此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效")
		}
	case "检测拿货时效":
		{
			ts.effective = rand.Intn(2)
		}
	default: //其它action

	}

	return nil
}

func createNetOrder(ts *StallPurchaseOrder, fsm *StateMachine) {
	event := Urgent
	switch ts.priority {
	case 0:
		event = NotUrgent
		fmt.Println("生成次日网络采购单")
	case 1:
		fmt.Println("生成当日网络采购单")
		event = Urgent
	}

	err := fsm.Trigger(ts.CurState, event, ts, fsm)
	if err != nil {
		fmt.Printf("trigger err: %v", err)
	}

}

func (p *StallPurchaseOrderEventProcessor) OnEnter(toState State, args []interface{}) {
	t := args[0].(*StallPurchaseOrder)
	t.CurState = toState
	t.States = append(t.States, toState)

	log.Printf("采购单 %d 的状态改变为 %v ", t.ID, toState)
}

func (p *StallPurchaseOrderEventProcessor) OnActionFailure(action Action, fromState State, toState State, args []interface{}, err error) {
	t := args[0].(*StallPurchaseOrder)

	log.Printf("采购单 %d 的状态从 %v to %v 改变失败， 原因: %v", t.ID, fromState, toState, err)
}

func TestStallOrder(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	ts := &StallPurchaseOrder{
		ID:       1,
		CurState: PendingPuchase,
		States:   []State{PendingPuchase},
	}


	
	ts.priority = int64(rand.Intn(2))
	fsm := initStallOrderFSM()

	err := fsm.Trigger(ts.CurState, Go2Stall, ts, fsm)
	if err != nil {
		t.Errorf("trigger err: %v", err)
	}

	event := GotTheGoods
	switch ts.GetGoodsResults {
	case 0:
		event = GotTheGoods
	case 1:
		event = HasBeenRemoved
	case 2:
		event = TwoDaysLater
	case 3:
		event = OneDaysLater
	}

	err = fsm.Trigger(ts.CurState, event, ts, fsm)
	ts.PreEvent = event
	if err != nil {
		fmt.Printf("trigger err: %v", err)
	}
	switch ts.GetGoodsResults {
	case 0:

		err = fsm.Trigger(ts.CurState, SendToWarehouse, ts, fsm)
		ts.PreEvent = SendToWarehouse
		if err != nil {
			t.Errorf("trigger err: %v", err)
		}

		err = fsm.Trigger(ts.CurState, Received, ts, fsm)
		if err != nil {
			t.Errorf("trigger err: %v", err)
		}

	case 1:
		createNetOrder(ts, fsm)
	case 2:
		createNetOrder(ts, fsm)
	case 3:
		event := NotExpired
		switch ts.effective {
		case 0:
			event = NotExpired
		case 1:
			event = Overdue
		}

		err := fsm.Trigger(ts.CurState, event, ts, fsm)
		ts.PreEvent = event
		if err != nil {
			fmt.Printf("trigger err: %v", err)
		}

		if ts.effective == 1 {
			createNetOrder(ts, fsm)
		}

	}

	fsm.Export("state.png")

}

func compareStallPurchaseOrder(t1 *StallPurchaseOrder, t2 *StallPurchaseOrder) bool {
	if t1.ID != t2.ID || t1.EventCount != t2.EventCount ||
		t1.CurState != t2.CurState {
		return false
	}

	return fmt.Sprint(t1.States) == fmt.Sprint(t2.States)
}

func initStallOrderFSM() *StateMachine {
	delegate := &DefaultDelegate{P: &StallPurchaseOrderEventProcessor{}}

	transitions := []Transition{
		{From: PendingPuchase, Event: Go2Stall, To: Purchasing, PreAction: GetTheResults},
		{From: Purchasing, Event: GotTheGoods, To: PurchaseSuccessful, PreAction: NoAction},
		{From: PurchaseSuccessful, Event: SendToWarehouse, To: BeingTransported, PreAction: NoAction},
		{From: BeingTransported, Event: Received, To: PurchaseEnd, PreAction: NoAction},
		{From: Purchasing, Event: HasBeenRemoved, To: Purchasefailed, PreAction: DisableSecondarySources},
		{From: Purchasefailed, Event: Urgent, To: PurchaseEnd, PreAction: GenerateTodayOnlineOrder},
		{From: Purchasefailed, Event: NotUrgent, To: PurchaseEnd, PreAction: GenerateTomorrowOnlineOrder},
		{From: Purchasing, Event: TwoDaysLater, To: Purchasefailed, PreAction: SourceFreeze2Days},
		{From: Purchasing, Event: OneDaysLater, To: Purchasing, PreAction: TestingTime},
		{From: Purchasing, Event: Overdue, To: Purchasefailed, PreAction: NoAction},
		{From: Purchasing, Event: NotExpired, To: PendingPuchase, PreAction: NoAction},
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
