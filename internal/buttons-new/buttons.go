package buttons

import (
	"log"
	"time"

	"github.com/warthog618/go-gpiocdev"
)

type NewButtonHandler struct {
	chip         *gpiocdev.Chip
	answerLine   *gpiocdev.Line
	rejectLine   *gpiocdev.Line
	answerOffset int
	rejectOffset int
	OnAnswer     func()
	OnReject     func()

	debounceDuration time.Duration
	lastAnswerPress  time.Time
	lastRejectPress  time.Time
}

// New initializes the GPIO chip and stores line offsets.
// chipPath may be something like "/dev/gpiochip4".
func New(chipPath string, answerOffset, rejectOffset int) *NewButtonHandler {
	chip, err := gpiocdev.NewChip(chipPath)
	if err != nil {
		log.Printf("GPIOCDEV Error: Failed to open chip %s: %v", chipPath, err)
		return nil
	}

	return &NewButtonHandler{
		chip:             chip,
		answerOffset:     answerOffset,
		rejectOffset:     rejectOffset,
		debounceDuration: 250 * time.Millisecond,
	}
}

func (h *NewButtonHandler) Start() {
	var err error

	h.answerLine, err = h.chip.RequestLine(h.answerOffset,
		gpiocdev.WithEventHandler(h.handleAnswerEvent),
		gpiocdev.WithFallingEdge,
		gpiocdev.WithPullUp,
	)
	if err != nil {
		log.Printf("GPIOCDEV Error: Failed to request answer line (offset %d): %v", h.answerOffset, err)
	}

	h.rejectLine, err = h.chip.RequestLine(h.rejectOffset,
		gpiocdev.WithEventHandler(h.handleRejectEvent),
		gpiocdev.WithFallingEdge,
		gpiocdev.WithPullUp,
	)
	if err != nil {
		log.Printf("GPIOCDEV Error: Failed to request reject line (offset %d): %v", h.rejectOffset, err)
	}
}

func (h *NewButtonHandler) handleAnswerEvent(event gpiocdev.LineEvent) {
	if time.Since(h.lastAnswerPress) < h.debounceDuration {
		return
	}
	h.lastAnswerPress = time.Now()
	if h.OnAnswer != nil {
		log.Println("Button (CDEV): Answer/Call pressed")
		h.OnAnswer()
	}
}

func (h *NewButtonHandler) handleRejectEvent(event gpiocdev.LineEvent) {
	if time.Since(h.lastRejectPress) < h.debounceDuration {
		return
	}
	h.lastRejectPress = time.Now()
	if h.OnReject != nil {
		log.Println("Button (CDEV): Reject/Finish pressed")
		h.OnReject()
	}
}

func (h *NewButtonHandler) Close() {
	if h.answerLine != nil {
		h.answerLine.Close()
	}
	if h.rejectLine != nil {
		h.rejectLine.Close()
	}
	if h.chip != nil {
		h.chip.Close()
	}
}
