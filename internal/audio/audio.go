package audio

import (
	"context"
	"log"

	"github.com/gen2brain/alsa"
)

type AudioHandler struct {
	In  chan []byte // Incoming audio chunks (expected to be Mono from capture)
	Out chan []byte // Outgoing audio chunks captured from device (Mono)

	CaptureCard    uint
	CaptureDevice  uint
	PlaybackCard   uint
	PlaybackDevice uint

	SampleRate       uint32
	CaptureChannels  uint32 // 1 (Mono)
	PlaybackChannels uint32 // 2 (Stereo)
	Format           alsa.PcmFormat

	PeriodSize  uint32
	PeriodCount uint32

	cancel context.CancelFunc
}

func NewAudioHandler(capCard, capDev, playCard, playDev uint) *AudioHandler {
	return &AudioHandler{
		In:               make(chan []byte, 10),
		Out:              make(chan []byte, 10),
		CaptureCard:      capCard,
		CaptureDevice:    capDev,
		PlaybackCard:     playCard,
		PlaybackDevice:   playDev,
		SampleRate:       44100,
		CaptureChannels:  1, // Capture in Mono
		PlaybackChannels: 2, // Playback in Stereo
		Format:           alsa.SNDRV_PCM_FORMAT_S16_LE,
		PeriodSize:       1024,
		PeriodCount:      4,
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
		Channels:    h.CaptureChannels,
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

	log.Printf("Audio Capture Started: Card %d, Device %d (Mono)", h.CaptureCard, h.CaptureDevice)

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
				log.Printf("Audio recording chunk")
			case <-ctx.Done():
				return
			}
		}
	}
}

func (h *AudioHandler) playbackRoutine(ctx context.Context) {
	config := &alsa.Config{
		Channels:    h.PlaybackChannels,
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

	log.Printf("Audio Playback Started: Card %d, Device %d (Stereo)", h.PlaybackCard, h.PlaybackDevice)

	for {
		select {
		case <-ctx.Done():
			log.Println("Audio Playback Stopped")
			return
		case monoData := <-h.In:
			stereoData := monoToStereo(monoData)
			
			_, err := p.Write(stereoData)
			log.Printf("Audio playback chunk")
			if err != nil {
				log.Printf("Audio Playback Write Error: %v", err)
			}
		}
	}
}

func monoToStereo(mono []byte) []byte {
	stereo := make([]byte, len(mono)*2)

	for i := 0; i < len(mono); i += 2 {
		if i+1 >= len(mono) {
			break 
		}

		// Grab the 2 bytes making up the 16-bit mono sample
		b1 := mono[i]
		b2 := mono[i+1]

		// Write to Left channel
		stereo[i*2] = b1
		stereo[i*2+1] = b2

		// Write to Right channel (duplicate of Left)
		stereo[i*2+2] = b1
		stereo[i*2+3] = b2
	}

	return stereo
}