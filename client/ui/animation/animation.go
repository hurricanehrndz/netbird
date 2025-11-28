// Package animation provides a thread-safe animator for systray icons.
// It supports configurable frame rates and graceful cancellation via context.
package animation

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"fyne.io/systray"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultFrameRate = 300 * time.Millisecond
)

type state struct {
	cancel context.CancelFunc
}

type Animator struct {
	ctx          context.Context
	state        atomic.Pointer[state]
	frames       map[string][]byte
	framesMu     sync.RWMutex
	animationIdx atomic.Int32
	numOfFrames  int32
	frameRate    time.Duration
	wg           sync.WaitGroup
}

// New creates a new Animator with the provided context and frames.
// The numOfFrames is automatically calculated from the length of the frames map.
// An optional frameRate can be provided; if omitted or zero, DefaultFrameRate (300ms) is used.
func New(ctx context.Context, frames map[string][]byte, frameRate ...time.Duration) *Animator {
	rate := DefaultFrameRate
	if len(frameRate) > 0 && frameRate[0] > 0 {
		rate = frameRate[0]
	}
	return &Animator{
		ctx:         ctx,
		frames:      frames,
		numOfFrames: int32(len(frames)),
		frameRate:   rate,
	}
}

// Start starts the animation.
func (a *Animator) Start() {
	if a.frames == nil || a.numOfFrames == 0 || a.ctx == nil {
		log.Error("animation manager not properly initialized")
		return
	}

	ctx, cancel := context.WithCancel(a.ctx)
	newState := &state{cancel: cancel}

	if !a.state.CompareAndSwap(nil, newState) {
		cancel()
		log.Debug("animation in progress")
		return
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		ticker := time.NewTicker(a.frameRate)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("stopped animation")
				return
			case <-ticker.C:
				idx := (a.animationIdx.Load() % a.numOfFrames) + 1
				a.animationIdx.Store(idx)
				frameKey := fmt.Sprintf("frame%d", idx)

				a.framesMu.RLock()
				frame := a.frames[frameKey]
				a.framesMu.RUnlock()

				systray.SetTemplateIcon(frame, frame)
			}
		}
	}()

	log.Info("animation started")
}

// Stop stops the animation.
func (a *Animator) Stop() {
	oldState := a.state.Swap(nil)

	if oldState == nil {
		log.Debug("animation is not running")
		return
	}

	oldState.cancel()
	a.wg.Wait()
}

// UpdateFrames updates the animation frames and resets the animation index.
// If animation is running, it will continue with the new frames.
func (a *Animator) UpdateFrames(frames map[string][]byte) {
	if frames == nil {
		log.Warn("attempted to update with empty frames")
		return
	}

	a.framesMu.Lock()
	a.frames = frames
	a.numOfFrames = int32(len(frames))
	a.framesMu.Unlock()

	a.animationIdx.Store(0)
}
