package buttons

import (
	"log"
	"time"

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

	concreteHandler := ConcreteButtonHandler{
		handler: handler,
		state:   s,
	}

	handler.OnAnswer = concreteHandler.onAnswer
	handler.OnReject = concreteHandler.onReject

	return &concreteHandler
}

func (h *ConcreteButtonHandler) Start() {
	h.handler.Start()
}

func (h *ConcreteButtonHandler) onAnswer() {
	h.state.AnswerButtonClickChan <- time.Now()
}

func (h *ConcreteButtonHandler) onReject() {
	h.state.RejectButtonClickChan <- time.Now()
}
