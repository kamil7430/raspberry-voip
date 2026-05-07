package audio

import (
	"context"
	"log"

	"github.com/gen2brain/alsa"
	"github.com/kamil7430/raspberry-voip/internal/config"
	resampling "github.com/tphakala/go-audio-resampling"
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

func NewAudioHandler() *AudioHandler {
	return &AudioHandler{
		In:               make(chan []byte, 10),
		Out:              make(chan []byte, 10),
		CaptureCard:      uint(config.LoadInt("captureCard")),
		CaptureDevice:    uint(config.LoadInt("captureDevice")),
		PlaybackCard:     uint(config.LoadInt("playbackCard")),
		PlaybackDevice:   uint(config.LoadInt("playbackDevice")),
		SampleRate:       uint32(config.LoadInt("sampleRate")),
		CaptureChannels:  1, // capture in mono
		PlaybackChannels: 2, // playback in stereo
		Format:           alsa.SNDRV_PCM_FORMAT_S16_LE,
		PeriodSize:       uint32(config.LoadInt("periodSize")),
		PeriodCount:      uint32(config.LoadInt("periodCount")),
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
	alsaConfig := &alsa.Config{
		Channels:    h.CaptureChannels,
		Rate:        h.SampleRate,
		Format:      h.Format,
		PeriodSize:  h.PeriodSize * h.SampleRate / 48000,
		PeriodCount: h.PeriodCount,
	}

	p, err := alsa.PcmOpen(h.CaptureCard, h.CaptureDevice, alsa.PCM_IN, alsaConfig)
	if err != nil {
		log.Printf("Audio Capture Error: Failed to open PCM device (Card %d, Device %d): %v", h.CaptureCard, h.CaptureDevice, err)
		return
	}
	defer p.Close()

	log.Printf("Audio Capture Started: Card %d, Device %d (Mono, %d Hz)", h.CaptureCard, h.CaptureDevice, h.SampleRate)

	resConfig := &resampling.Config{
		InputRate:  float64(h.SampleRate),
		OutputRate: 48000,
		Channels:   1,
		Quality:    resampling.QualitySpec{Preset: resampling.QualityLow},
	}

	var resampler resampling.Resampler
	if h.SampleRate != 48000 {
		resampler, err = resampling.New(resConfig)
		if err != nil {
			log.Printf("Audio Capture Resampler Error: %v", err)
			return
		}
	}

	bufferSize := alsa.PcmFramesToBytes(p, alsaConfig.PeriodSize)
	buf := make([]byte, bufferSize)

	for {
		select {
		case <-ctx.Done():
			log.Println("Audio Capture Stopped")
			return
		default:
			n, err := p.Read(buf)
			if err != nil {
				log.Printf("Audio Capture Read Error: %v", err)
				continue
			}

			data := buf[:n]
			var outBuf []byte

			if resampler != nil {
				floatBuf := bytesToFloat64(data)
				resampledFloats, err := resampler.Process(floatBuf)
				if err != nil {
					log.Printf("Audio Capture Resampling Error: %v", err)
					continue
				}
				outBuf = float64ToBytes(resampledFloats)
			} else {
				// Copy data to avoid modification if buffer is reused
				outBuf = make([]byte, len(data))
				copy(outBuf, data)
			}

			if len(outBuf) == 0 {
				continue
			}

			select {
			case h.Out <- outBuf:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (h *AudioHandler) playbackRoutine(ctx context.Context) {
	alsaConfig := &alsa.Config{
		Channels:    h.PlaybackChannels,
		Rate:        48000,
		Format:      h.Format,
		PeriodSize:  h.PeriodSize,
		PeriodCount: h.PeriodCount,
	}

	p, err := alsa.PcmOpen(h.PlaybackCard, h.PlaybackDevice, alsa.PCM_OUT, alsaConfig)
	if err != nil {
		log.Printf("Audio Playback Error: Failed to open PCM device (Card %d, Device %d): %v", h.PlaybackCard, h.PlaybackDevice, err)
		return
	}
	defer p.Close()

	log.Printf("Audio Playback Started: Card %d, Device %d (Stereo, 48000 Hz)", h.PlaybackCard, h.PlaybackDevice)

	for {
		select {
		case <-ctx.Done():
			log.Println("Audio Playback Stopped")
			return
		case monoData, ok := <-h.In:
			if !ok {
				return
			}
			stereoData := monoToStereo(monoData)

			_, err := p.Write(stereoData)
			if err != nil {
				log.Printf("Audio Playback Write Error: %v", err)
			}
		}
	}
}

func monoToStereo(mono []byte) []byte {
	// Assumes S16LE (2 bytes per sample)
	stereo := make([]byte, len(mono)*2)
	for i := 0; i+1 < len(mono); i += 2 {
		// Left channel
		stereo[i*2] = mono[i]
		stereo[i*2+1] = mono[i+1]
		// Right channel
		stereo[i*2+2] = mono[i]
		stereo[i*2+3] = mono[i+1]
	}
	return stereo
}

func bytesToFloat64(pcm []byte) []float64 {
	floats := make([]float64, len(pcm)/2)
	for i := 0; i+1 < len(pcm); i += 2 {
		val := int16(pcm[i]) | (int16(pcm[i+1]) << 8)
		floats[i/2] = float64(val) / 32768.0
	}
	return floats
}

func float64ToBytes(floats []float64) []byte {
	pcm := make([]byte, len(floats)*2)
	for i, f := range floats {
		val := int32(f * 32768.0)
		if val > 32767 {
			val = 32767
		} else if val < -32768 {
			val = -32768
		}
		pcm[i*2] = byte(val)
		pcm[i*2+1] = byte(val >> 8)
	}
	return pcm
}
