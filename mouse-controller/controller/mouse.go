package controller

import (
	"github.com/go-vgo/robotgo"
	"github.com/snakesneaks/virtual-mouse-camera/mouse-controller/controller/worker"
	"github.com/snakesneaks/virtual-mouse-camera/mouse-controller/model"
)

type MouseController interface {
	// move mouse position to relative position in screen. X and Y ranges 0 ~ 1.0.
	Move(x, y float64) error

	//when mouse has input
	LeftDown(bool)
	MiddleDown(bool)
	RightDown(bool)
}

type mouseController struct {
	sx, sy int
	worker worker.MouseWorker
}

func NewMouseController() MouseController {
	x, y := robotgo.GetScreenSize()
	c := mouseController{
		sx:     x,
		sy:     y,
		worker: worker.NewMouseWorker(),
	}
	c.worker.Start()
	return &c
}

func (m *mouseController) Move(x, y float64) error {
	//log.Printf("mousePos: (%f,%f)", x, y)

	//mouse position to screen position
	sx := int(float64(m.sx) * x)
	sy := int(float64(m.sy) * y)

	m.worker.Move(sx, sy)
	return nil
}

func (m *mouseController) LeftDown(b bool) {
	m.worker.MouseInput(model.MouseLeftDown, b)
}

func (m *mouseController) MiddleDown(b bool) {
	m.worker.MouseInput(model.MouseMiddleDown, b)
}

func (m *mouseController) RightDown(b bool) {
	m.worker.MouseInput(model.MouseRightDown, b)
}

/*
func main() {

	robotgo.MouseSleep = 10

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	fmt.Println("Ready. Press any key ...")

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Key press => " + string(r))

		switch r {
		case 'a':
			robotgo.MoveRelative(-1, 0)
		case 'w':
			robotgo.MoveRelative(0, -1)
		case 's':
			robotgo.MoveRelative(0, 1)
		case 'd':
			robotgo.MoveRelative(1, 0)
		case ' ':
			robotgo.Click("left", true)
		default:
			log.Println("unhandling keyboard input")
		}
	}

}
*/
