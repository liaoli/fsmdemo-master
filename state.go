package fsm

/*
档口订单的状态
*/

var PendingPuchase = State{
	0,
	"待采购",
	"",
}

var Purchasing = State{
	1,
	"采购中",
	"",
}

var Purchasefailed = State{
	2,
	"采购失败",
	"",
}

var PurchaseSuccessful = State{
	3,
	"采购成功",
	"",
}
var BeingTransported = State{
	4,
	"运输中",
	"",
}
var PurchaseEnd = State{
	5,
	"结束",
	"",
}

var Invalid = State{
	0,
	"无效状态",
	"",
}
