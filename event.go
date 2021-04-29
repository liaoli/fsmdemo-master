package fsm

/*
档口订单的事件
*/

var Go2Stall = Event{
	0,
	"去档口拿货",
}

var GotTheGoods = Event{
	1,
	"拿货成功",
}

var SendToWarehouse = Event{
	2,
	"运送中",
}

var Received = Event{
	3,
	"已收货",
}

var HasBeenRemoved = Event{
	4,
	"该商品已经下架",
}
var TwoDaysLater = Event{
	5,
	"后天拿到货",
}

var OneDaysLater = Event{
	6,
	"明天拿到货",
}

var Overdue = Event{
	7,
	"拿货超过有效期",
}

var NotExpired = Event{
	8,
	"拿货未超过有效期",
}

var Urgent = Event{
	9,
	"优先级(紧急)",
}

var NotUrgent = Event{
	10,
	"优先级(非紧急)",
}
