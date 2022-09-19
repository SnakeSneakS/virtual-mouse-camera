package model

import (
	"log"
	"sync"

	"golang.org/x/xerrors"
)

// Mouse Event
type MouseInputEvent int

const (
	MouseLeftDown   MouseInputEvent = iota
	MouseMiddleDown                 //wheel
	MouseRightDown
)

// EventSequence

type EventSequenceState int

// down means "ON"
const (
	HoldEnd  EventSequenceState = iota //end of holding
	HoldDown                           //started holding
	Click                              //clicked
	Stay                               //no change
	Unknown
)

// Available Mouse Events
var MouseEvents []MouseInputEvent = []MouseInputEvent{MouseLeftDown, MouseMiddleDown, MouseRightDown}

// EventStore
type MouseInputEventStore struct {
	m map[MouseInputEvent]bool
	*sync.Mutex
}

func NewMouseInputEventStore() MouseInputEventStore {
	return MouseInputEventStore{
		make(map[MouseInputEvent]bool, 0),
		new(sync.Mutex),
	}
}

func (s MouseInputEventStore) Set(e MouseInputEvent, b bool) {
	s.Lock()
	defer s.Unlock()

	s.m[e] = b
}
func (s MouseInputEventStore) Get(e MouseInputEvent) (bool, error) {
	s.Lock()
	defer s.Unlock()

	if v, ok := s.m[e]; ok {
		delete(s.m, e)
		return v, nil
	}
	return false, xerrors.New("no value is set")
}

// Several Events Store
type MouseInputEventsStore struct {
	m map[MouseInputEvent]([]bool)
	*sync.Mutex
}

func NewMouseInputEventsStore() MouseInputEventsStore {
	return MouseInputEventsStore{
		make(map[MouseInputEvent]([]bool), 0),
		new(sync.Mutex),
	}
}

func (s MouseInputEventsStore) Push(e MouseInputEvent, b bool) {
	s.Lock()
	defer s.Unlock()

	s.m[e] = append(s.m[e][1:], b)
}

func (s MouseInputEventsStore) GetAll(e MouseInputEvent) []bool {
	s.Lock()
	defer s.Unlock()
	return s.m[e]
}

// combineSeq: combine some sequence and do majority decision to prevent chattering
// holdSeq is an sequence required to hold down.
func (s MouseInputEventsStore) GetState(e MouseInputEvent, combineSeq, holdSeq int) EventSequenceState {
	store, ok := s.m[e]
	if !ok {
		log.Println("failed to get state of event: ", e)
		return Unknown
	}

	size := len(store)
	storeCombined := make([]bool, 0)
	i := 0
	for i+combineSeq < size {
		t, f := 0, 0
		for _, b := range store[i : i+combineSeq] {
			if b {
				t++
			} else {
				f++
			}
		}
		storeCombined = append(storeCombined, t >= f)

		i += combineSeq
	}

	type eventSequenceStateWithCount struct {
		isDown bool //isDown
		count  int  //sequence count
	}

	i = len(storeCombined) - 1
	b := storeCombined[i]
	c := 1
	inputs := make([]eventSequenceStateWithCount, 0)
	for i > 0 {
		i--
		if b == storeCombined[i] {
			c++
		} else {
			//押した又は離した
			inputs = append(inputs, eventSequenceStateWithCount{b, c})

			//少なくとも2つわかったら終了
			if len(inputs) == 2 {
				break
			}

			b = storeCombined[i]
			c = 1
		}
	}

	if len(inputs) == 0 {
		return Stay
	} else if len(inputs) == 1 {
		//ちょうどhold時間に達した場合、holdDown又はholdUpに達したというこ
		return Stay
	} else {
		first := inputs[1]
		second := inputs[0]
		if first.isDown == second.isDown {
			return Unknown
		} else if first.isDown && !second.isDown {
			if second.count == 1 {
				if first.count >= holdSeq {
					return HoldEnd
				} else {
					return Click
				}
			} else {
				return Stay
			}

		} else if !first.isDown && second.isDown {
			if second.count == holdSeq {
				return HoldDown
			} else {
				return Stay
			}
		} else {
			return Unknown
		}
	}
}

func (s MouseInputEventsStore) Init(e MouseInputEvent, size int) {
	s.m[e] = make([]bool, size)
}
