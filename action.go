package fsm


var NoAction = Action{

}
var GetTheResults = Action{
	1,
	"获取拿货结果",
}

var DisableSecondarySources = Action{
	2,
	"禁用次货源",
}
var GenerateTodayOnlineOrder = Action{
	3,
	"生成当日网络采购单",
}
var GenerateTomorrowOnlineOrder = Action{
	4,
	"生成次日网络采购单",
}
var SourceFreeze2Days = Action{
	5,
	"此货源冻结2天2天不使用此货源采购，关闭拿货采购单状态无效",
}
var TestingTime = Action{
	6,
	"检测拿货时效",
}
