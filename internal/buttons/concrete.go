package buttons

import (
	"log"
	"time"

	"github.com/kamil7430/raspberry-voip/internal/audio"
	"github.com/kamil7430/raspberry-voip/internal/config"
	"github.com/kamil7430/raspberry-voip/internal/display"
	"github.com/kamil7430/raspberry-voip/internal/state"
	"github.com/kamil7430/raspberry-voip/internal/tcp"
)

type ConcreteButtonHandler struct {
	handler *NewButtonHandler
	state   *state.State
	display *display.DisplayController
	audio   *audio.AudioHandler
}

func NewConcreteButtonHandler(s *state.State, d *display.DisplayController, a *audio.AudioHandler) *ConcreteButtonHandler {
	handler := New(
		config.LoadString("chipPath"),
		config.LoadInt("answerGpioPin"),
		config.LoadInt("rejectGpioPin"),
	)
	if handler == nil {
		log.Fatal("could not create chip button handler")
	}

	concreteHandler := ConcreteButtonHandler{
		handler: handler,
		state:   s,
		display: d,
		audio:   a,
	}

	handler.OnAnswer = concreteHandler.onAnswer
	handler.OnReject = concreteHandler.onReject

	return &concreteHandler
}

func (h *ConcreteButtonHandler) Start() {
	h.handler.Start()
}

func (h *ConcreteButtonHandler) onAnswer() {
	if h.state.GetConnectionContext() != nil {
		select {
		case h.state.AnswerButtonClickChan <- time.Now():
		default:
		}
	} else {
		err := tcp.Dial(
			h.state.GetDialingAddress(),
			h.state,
			h.display,
			h.audio,
		)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *ConcreteButtonHandler) onReject() {
	select {
	case h.state.RejectButtonClickChan <- time.Now():
	default:
	}
}
