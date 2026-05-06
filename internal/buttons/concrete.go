package buttons

import (
	"log"

	"github.com/kamil7430/raspberry-voip/internal/state"
)

const (
	chipPath      = "/dev/gpiochip0"
	answerGpioPin = 18
	rejectGpioPin = 25
)

type ConcreteButtonHandler struct {
	handler *NewButtonHandler
	state   *state.State
}

func NewConcreteButtonHandler(s *state.State) *ConcreteButtonHandler {
	handler := New(
		chipPath,
		answerGpioPin,
		rejectGpioPin,
	)
	if handler == nil {
		log.Fatal("could not create chip button handler")
	}

	handler.OnAnswer = onAnswer
	handler.OnReject = onReject

	return &ConcreteButtonHandler{
		handler: handler,
		state:   s,
	}
}

func (h *ConcreteButtonHandler) Start() {
	h.handler.Start()
}

func onAnswer() {

}

func onReject() {

}
