package buttons

/*
import (
	"log"
	"time"

	"github.com/warthog618/gpio"
)

type ButtonHandler struct {
	answerPin *gpio.Pin
	rejectPin *gpio.Pin

	OnAnswer func()
	OnReject func()

	debounceDuration time.Duration
	lastAnswerPress  time.Time
	lastRejectPress  time.Time
}

func NewButtonHandler(answerBCM, rejectBCM int) *ButtonHandler {
	err := gpio.Open()
	if err != nil {
		log.Printf("GPIO Error: Failed to open: %v", err)
		return nil
	}

	h := &ButtonHandler{
		answerPin:        gpio.NewPin(answerBCM),
		rejectPin:        gpio.NewPin(rejectBCM),
		debounceDuration: 250 * time.Millisecond,
	}

	h.answerPin.Input()
	h.answerPin.PullUp()

	h.rejectPin.Input()
	h.rejectPin.PullUp()

	return h
}

func (h *ButtonHandler) Start() {

	h.answerPin.Watch(gpio.EdgeFalling, func(p *gpio.Pin) {
		if p.Read() != gpio.Low {
			return
		}
		if time.Since(h.lastAnswerPress) < h.debounceDuration {
			return
		}
		h.lastAnswerPress = time.Now()
		if h.OnAnswer != nil {
			log.Println("Button: Answer/Call pressed")
			h.OnAnswer()
		}
	})

	h.rejectPin.Watch(gpio.EdgeFalling, func(p *gpio.Pin) {
		if p.Read() != gpio.Low {
			return
		}
		if time.Since(h.lastRejectPress) < h.debounceDuration {
			return
		}
		h.lastRejectPress = time.Now()
		if h.OnReject != nil {
			log.Println("Button: Reject/Finish pressed")
			h.OnReject()
		}
	})
}

func (h *ButtonHandler) Close() {
	if h.answerPin != nil {
		h.answerPin.Unwatch()
	}
	if h.rejectPin != nil {
		h.rejectPin.Unwatch()
	}
	gpio.Close()
}
*/
