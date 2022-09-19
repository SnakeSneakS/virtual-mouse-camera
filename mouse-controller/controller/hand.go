package controller

import (
	"encoding/json"
	"log"
	"math"

	"github.com/snakesneaks/virtual-mouse-camera/mouse-controller/model"
)

const sensitivity = 1.5

type handler struct {
	mouseController MouseController
}

func NewHandler() handler {
	return handler{
		mouseController: NewMouseController(),
	}
}

func (h *handler) HandleHandLandmarksBytes(b []byte) error {

	var landmarks model.HandLandmark
	if err := json.Unmarshal(b, &landmarks); err != nil {
		log.Fatal(err)
	}

	//log.Println(landmarks)
	if err := h.handleHandLandmarks(landmarks); err != nil {
		return err
	}

	return nil
}

func (h *handler) handleHandLandmarks(landmarks model.HandLandmark) error {
	//log.Println(landmarks)
	index_finger_mcp := landmarks[5]
	index_finger_tip := landmarks[8]
	middle_finger_mcp := landmarks[9]
	middle_finger_tip := landmarks[12]
	//ring_finger_mcp := landmarks[13]
	//ring_finger_tip := landmarks[16]

	mid := index_finger_tip.Add(middle_finger_tip).Multiply(model.NewPosition(0.5, 0.5, 0.5)) //midth of index finger & middle finger
	mousePosRaw := mid.Multiply(model.NewPosition(-1, 1, 1)).Add(model.NewPosition(1, 0, 0))  //arrangement
	mousePosRefined := h.coodinatePositionWithSensitivity(mousePosRaw, sensitivity)

	if h.isTouched(index_finger_mcp, index_finger_tip, middle_finger_mcp, middle_finger_tip) {
		h.mouseController.LeftDown(true)
	}
	/*else if h.isTouched(middle_finger_mcp, middle_finger_tip, ring_finger_mcp, ring_finger_tip) {
		h.mouseController.RightDown()
	}*/

	if err := h.mouseController.Move(mousePosRefined.X, mousePosRefined.Y); err != nil {
		return err
	}
	return nil
}

// 座標に少しの値をかけることによって、隅々までマウスを動かせるようにする
func (h *handler) coodinatePositionWithSensitivity(p model.Position, sensitivity float64) model.Position {
	mp := p.Sub(
		model.NewPosition(0.5, 0.5, 0.5),
	).Multiply(
		model.NewPosition(sensitivity, sensitivity, sensitivity),
	).Add(
		model.NewPosition(0.5, 0.5, 0.5),
	)
	if mp.X > 1 {
		mp.X = 1
	}
	if mp.Y > 1 {
		mp.Y = 1
	}
	return mp
}

// p1_1からp1_2の指がp2_1からp2_2にくっついているかどうか
// テキトーに計算してる
func (h *handler) isTouched(p1_1, p1_2, p2_1, p2_2 model.Position) bool {
	vec1 := p1_2.Sub(p1_1)
	vec2 := p2_2.Sub(p2_1)
	sum := vec1.Add(vec2)
	between := vec1.Sub(vec2)

	c := (sum.X + sum.Y + sum.Z) / (between.X + between.Y + between.Z)

	return math.Abs(c) > 10
}
