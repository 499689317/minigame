/***********************************************************************
*【单线程】时间轮
***********************************************************************/
package timer

import (
	"time"
)

const (
	WHEEL_NUM     = 5
	TIME_TICK_LEN = 25
)

var (
	WHEEL_BIT  = [...]uint{8, 6, 6, 6, 5}
	WHEEL_SIZE [WHEEL_NUM]uint
	WHEEL_CAP  [WHEEL_NUM]uint
)

// ------------------------------------------------------------
// node
type TimerNode struct {
	prev     *TimerNode
	next     *TimerNode
	timeDead int
	interval int
	total    int
	callback func()
}

func NewTimerNode() *TimerNode {
	ret := new(TimerNode)
	ret._init()
	return ret
}
func (self *TimerNode) _init() {
	self.prev = self
	self.next = self //circle
}
func (self *TimerNode) _Callback() {
	self.total -= self.interval
	if self.total > 0 {
		self.timeDead += self.interval
		G_TimerMgr._AddTimerNode(self.interval, self)
	}
	self.callback()
}

// ------------------------------------------------------------/
// wheel
type Wheel struct {
	//每个slot维护的node链表为一个环，如此可以简化插入删除的操作
	//slot.next为node链表中第一个节点，prev为node的最后一个节点
	slots   []TimerNode
	slotIdx uint
}

func NewWheel(size int) *Wheel {
	ret := new(Wheel)
	ret.slots = make([]TimerNode, size)
	for i := 0; i < size; i++ {
		ret.slots[i]._init()
	}
	return ret
}
func (self *Wheel) GetCurSlot() *TimerNode {
	return &self.slots[self.slotIdx]
}
func (self *Wheel) size() uint { return uint(len(self.slots)) }

// ------------------------------------------------------------
// manager
var G_TimerMgr = _NewTimerMgr()

type TimerMgr struct {
	wheels    [WHEEL_NUM]Wheel
	readyNode TimerNode
}

func _NewTimerMgr() *TimerMgr {
	if len(WHEEL_BIT) != WHEEL_NUM {
		panic("WHEEL_NUM isn't matching of WHEEL_BIT")
	}
	ret := new(TimerMgr)
	ret.readyNode._init()
	for i := 0; i < WHEEL_NUM; i++ {
		var wheelCap uint
		for j := 0; j <= i; j++ {
			wheelCap += WHEEL_BIT[j]
		}
		WHEEL_CAP[i] = 1 << wheelCap
		WHEEL_SIZE[i] = 1 << WHEEL_BIT[i]
		ret.wheels[i].slots = make([]TimerNode, WHEEL_SIZE[i])
	}
	return ret
}
func (self *TimerMgr) Refresh(time_elapse, timenow int) {
	tickCnt := time_elapse / TIME_TICK_LEN
	for i := 0; i < tickCnt; i++ { //扫过的slot均超时
		isCascade := false
		wheel := &self.wheels[0]
		slot := wheel.GetCurSlot()

		wheel.slotIdx++
		if wheel.slotIdx >= wheel.size() {
			wheel.slotIdx = 0
			isCascade = true
		}
		node := slot.next
		slot.next = slot //清空当前格子
		slot.prev = slot
		for node != slot { //环形链表遍历
			tmp := node
			node = node.next //得放在前面，后续函数调用，可能会更改node的链接关系
			self._AddToReadyNode(tmp)
		}
		if isCascade {
			self._Cascade(1, timenow) //跳级
		}
	}
	self._DoTimeOutCallBack()
}
func (self *TimerMgr) AddTimerSec(callback func(), delay, cd, total int) {
	node := NewTimerNode()
	node.callback = callback
	node.interval = cd * 1000
	node.total = total * 1000
	node.timeDead = time.Now().Nanosecond()/int(time.Millisecond) + delay*1000
	self._AddTimerNode(delay*1000, node)
}
func (self *TimerMgr) _AddTimerNode(msec int, node *TimerNode) {
	var slot *TimerNode
	tickCnt := uint(msec / TIME_TICK_LEN)
	if tickCnt < WHEEL_CAP[0] {
		idx := (self.wheels[0].slotIdx + tickCnt) & (WHEEL_SIZE[0] - 1) //2的N次幂位操作取余
		slot = &self.wheels[0].slots[idx]
	} else {
		for i := 1; i < WHEEL_NUM; i++ {
			if tickCnt < WHEEL_CAP[i] {
				preCap := WHEEL_CAP[i-1] //上一级总容量即为本级的一格容量
				idx := (self.wheels[i].slotIdx + tickCnt/preCap - 1) & (WHEEL_SIZE[i] - 1)
				slot = &self.wheels[i].slots[idx]
				break
			}
		}
	}
	node.prev = slot.prev //插入格子的prev位置(尾节点)
	node.prev.next = node
	node.next = slot
	slot.prev = node
}
func (self *TimerMgr) RemoveTimer(node *TimerNode) {
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	node.prev = nil
	node.next = nil
}
func (self *TimerMgr) _Cascade(wheelIdx int, timenow int) {
	if wheelIdx < 1 || wheelIdx >= WHEEL_NUM {
		return
	}
	isCascade := false
	wheel := &self.wheels[wheelIdx]
	slot := wheel.GetCurSlot()
	//【Bug】须先更新槽位————扫格子时加新Node，不能再放入当前槽位了
	wheel.slotIdx++
	if wheel.slotIdx >= wheel.size() {
		wheel.slotIdx = 0
		isCascade = true
	}
	link := slot.next
	slot.next = slot //清空当前格子
	slot.prev = slot
	for link != slot {
		node := link
		link = link.next
		if node.timeDead <= timenow {
			self._AddToReadyNode(node)
		} else {
			//【Bug】加新Node，须加到其它槽位，本槽位已扫过(失效，等一整轮才会再扫到)
			self._AddTimerNode(node.timeDead-timenow, node)
		}
	}
	if isCascade {
		self._Cascade(wheelIdx+1, timenow)
	}
}
func (self *TimerMgr) _AddToReadyNode(node *TimerNode) {
	node.prev = self.readyNode.prev
	node.prev.next = node
	node.next = &self.readyNode
	self.readyNode.prev = node
}
func (self *TimerMgr) _DoTimeOutCallBack() {
	node := self.readyNode.next
	for node != &self.readyNode {
		tmp := node
		node = node.next
		tmp._Callback()
	}
	self.readyNode.next = &self.readyNode
	self.readyNode.prev = &self.readyNode
}
