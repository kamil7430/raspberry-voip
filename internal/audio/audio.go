package audio

import (
	"context"
	"log"

	"github.com/gen2brain/alsa"
)

type AudioHandler struct {
	In  chan []byte // Incoming audio chunks to be played
	Out chan []byte // Outgoing audio chunks captured from device

	CaptureCard    uint
	CaptureDevice  uint
	PlaybackCard   uint
	PlaybackDevice uint

	SampleRate uint32
	Channels   uint32
	Format     alsa.PcmFormat

	PeriodSize  uint32
	PeriodCount uint32

	cancel context.CancelFunc
}

func NewAudioHandler(capCard, capDev, playCard, playDev uint) *AudioHandler {
	return &AudioHandler{
		In:             make(chan []byte, 10),
		Out:            make(chan []byte, 10),
		CaptureCard:    capCard,
		CaptureDevice:  capDev,
		PlaybackCard:   playCard,
		PlaybackDevice: playDev,
		SampleRate:     44100,
		Channels:       1,
		Format:         alsa.SNDRV_PCM_FORMAT_S16_LE,
		PeriodSize:     1024,
		PeriodCount:    4,
	}
}

func (h *AudioHandler) Start(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	h.cancel = cancel

	go h.captureRoutine(ctx)
	go h.playbackRoutine(ctx)
}

func (h *AudioHandler) Stop() {
	if h.cancel != nil {
		h.cancel()
		h.cancel = nil
	}
}

func (h *AudioHandler) captureRoutine(ctx context.Context) {
	config := &alsa.Config{
		Channels:    h.Channels,
		Rate:        h.SampleRate,
		Format:      h.Format,
		PeriodSize:  h.PeriodSize,
		PeriodCount: h.PeriodCount,
	}

	p, err := alsa.PcmOpen(h.CaptureCard, h.CaptureDevice, alsa.PCM_IN, config)
	if err != nil {
		log.Printf("Audio Capture Error: Failed to open PCM device (Card %d, Device %d): %v", h.CaptureCard, h.CaptureDevice, err)
		return
	}
	defer p.Close()

	log.Printf("Audio Capture Started: Card %d, Device %d", h.CaptureCard, h.CaptureDevice)

	bufferSize := alsa.PcmFramesToBytes(p, config.PeriodSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("Audio Capture Stopped")
			return
		default:
			buf := make([]byte, bufferSize)
			_, err := p.Read(buf)
			if err != nil {
				log.Printf("Audio Capture Read Error: %v", err)
				continue
			}

			select {
			case h.Out <- buf:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (h *AudioHandler) playbackRoutine(ctx context.Context) {
	config := &alsa.Config{
		Channels:    h.Channels,
		Rate:        h.SampleRate,
		Format:      h.Format,
		PeriodSize:  h.PeriodSize,
		PeriodCount: h.PeriodCount,
	}

	p, err := alsa.PcmOpen(h.PlaybackCard, h.PlaybackDevice, alsa.PCM_OUT, config)
	if err != nil {
		log.Printf("Audio Playback Error: Failed to open PCM device (Card %d, Device %d): %v", h.PlaybackCard, h.PlaybackDevice, err)
		return
	}
	defer p.Close()

	log.Printf("Audio Playback Started: Card %d, Device %d", h.PlaybackCard, h.PlaybackDevice)

	for {
		select {
		case <-ctx.Done():
			log.Println("Audio Playback Stopped")
			return
		case data := <-h.In:
			_, err := p.Write(data)
			if err != nil {
				log.Printf("Audio Playback Write Error: %v", err)
			}
		}
	}
}
