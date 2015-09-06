package util

import (
	"math"
	"sort"
	"strings"
)

// =========================================================================
// 5分钟内topN
type KeyItem struct {
	Key       string
	Timestamp int64
	Item      float64
}

type TopItem []*KeyItem

func (tr TopItem) Len() int {
	return len(tr)
}

func (tr TopItem) Less(i, j int) bool {
	tt := time.Now().Unix()
	if tt-tr[i].Timestamp > 300 && tt-tr[j].Timestamp < 300 {
		return true
	}
	if tt-tr[i].Timestamp < 300 && tt-tr[j].Timestamp > 300 {
		return false
	}
	return tr[i].Item < tr[j].Item
}

func (tr TopItem) Swap(i, j int) {
	tr[i], tr[j] = tr[j], tr[i]
}

type TimeTopN struct {
	num  int
	topN TopItem
}

func NewTimeTopN(n int) *TimeTopN {
	return &TimeTopN{n, make([]*KeyItem, 0, n+1)}
}

func (ttn *TimeTopN) Insert(ar *KeyItem) {
	if ttn.topN.Len() >= ttn.num {
		ttn.topN = append(ttn.topN[:ttn.num], ar)
	} else {
		ttn.topN = append(ttn.topN, ar)
	}
	sort.Sort(sort.Reverse(ttn.topN))
}

func (ttn *TimeTopN) String() string {
	var s string
	n := len(ttn.topN)
	if n > ttn.num {
		n = ttn.num
	}
	for k, v := range ttn.topN[:n] {
		s += fmt.Sprintf("Seq:%d Key:%s Timestamp:%d Item:%f.", k, v.Key, v.Timestamp, v.Item)
	}
	return s
}

func (ttn *TimeTopN) Query() []*KeyItem {
	n := len(ttn.topN)
	if n > ttn.num {
		n = ttn.num
	}
	return ttn.topN[:n]
}
