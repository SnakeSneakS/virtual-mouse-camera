package worker

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/snakesneaks/virtual-mouse-camera/mouse-controller/model"
)

const FPS = 30
const waitKeyTime = 1000 / FPS
const historySize = 100
const combineSequence = 1                        //チャタリング予防. ただしこれをしたらクリックイベントが複数起きるなど問題ある.
const holdSequence = FPS / (combineSequence * 3) //何度eventが続いたらholdとみなすか

type MouseWorker interface {
	Start()
	Stop()

	Move(sx, sy int)
	MouseInput(model.MouseInputEvent, bool)
}

type mouseWorker struct {
	isWorking bool

	mouseEventStore  *model.MouseInputEventStore  //現在のマウスイベントを保持する
	mouseEventsStore *model.MouseInputEventsStore //過去のマウスイベントを保持する
}

func NewMouseWorker() MouseWorker {
	w := mouseWorker{
		isWorking:        false,
		mouseEventStore:  model.NewMouseInputEventStore(),
		mouseEventsStore: model.NewMouseInputEventsStore(),
	}
	w.mouseEventsStore.Init(model.MouseLeftDown, historySize)
	w.mouseEventsStore.Init(model.MouseMiddleDown, historySize)
	w.mouseEventsStore.Init(model.MouseRightDown, historySize)
	return &w
}

func (w *mouseWorker) Start() {
	if w.isWorking {
		log.Println("mouse controller is already working")
		return
	}
	w.isWorking = true
	log.Println("mouse controller started working")
	go func() {
		for w.isWorking {
			for _, event := range model.MouseEvents {
				if b, err := w.mouseEventStore.Get(event); err == nil {
					w.mouseEventsStore.Push(event, b)
				} else {
					w.mouseEventsStore.Push(event, false)
				}

				eventState := w.mouseEventsStore.GetState(event, combineSequence, holdSequence)
				switch eventState {
				case model.Click:
					log.Println("state: click")
					robotgo.Click(w.mouseInputEvent2mouseInputKey(event))
				case model.HoldEnd:
					log.Println("state: holdEnd")
					robotgo.Toggle(w.mouseInputEvent2mouseInputKey(event), "up")
				case model.HoldDown:
					log.Println("state: holdDown")
					robotgo.Toggle(w.mouseInputEvent2mouseInputKey(event), "down")
				case model.Stay:
					//do nothing
				case model.Unknown:
					log.Println("unknown state: ", eventState)
				}
			}
			time.Sleep(time.Millisecond * waitKeyTime)
		}
		log.Println("mouse controller stopped")
	}()
}

func (w *mouseWorker) Stop() {
	if !w.isWorking {
		log.Println("mouse controller is not working")
	}
	w.isWorking = false
}

func (w *mouseWorker) Move(sx, sy int) {
	robotgo.MoveSmooth(sx, sy)
}

// when down
func (w *mouseWorker) MouseInput(e model.MouseInputEvent, b bool) {
	w.mouseEventStore.Set(e, b)
}

func (w *mouseWorker) mouseInputEvent2mouseInputKey(e model.MouseInputEvent) string {
	switch e {
	case model.MouseLeftDown:
		return "left"
	case model.MouseMiddleDown:
		return "center"
	case model.MouseRightDown:
		return "right"
	default:
		log.Fatalln("unhandling input event!")
		return ""
	}
}
