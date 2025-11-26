package animation

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		frames     map[string][]byte
		wantFrames int32
	}{
		{
			name: "valid 3 frames",
			ctx:  context.Background(),
			frames: map[string][]byte{
				"frame1": []byte("data1"),
				"frame2": []byte("data2"),
				"frame3": []byte("data3"),
			},
			wantFrames: 3,
		},
		{
			name: "valid 1 frame",
			ctx:  context.Background(),
			frames: map[string][]byte{
				"frame1": []byte("data1"),
			},
			wantFrames: 1,
		},
		{
			name:       "empty frames",
			ctx:        context.Background(),
			frames:     map[string][]byte{},
			wantFrames: 0,
		},
		{
			name:       "nil frames",
			ctx:        context.Background(),
			frames:     nil,
			wantFrames: 0,
		},
		{
			name: "nil context",
			ctx:  nil,
			frames: map[string][]byte{
				"frame1": []byte("data1"),
			},
			wantFrames: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := New(tt.ctx, tt.frames)

			require.NotNil(t, am)
			assert.Equal(t, tt.wantFrames, am.numOfFrames)
			assert.Equal(t, tt.ctx, am.ctx)
			assert.Equal(t, tt.frames, am.frames)
			assert.Equal(t, int32(0), am.animationIdx.Load())
			assert.Nil(t, am.state.Load())
		})
	}
}

func TestAnimator_Start_Validation(t *testing.T) {
	tests := []struct {
		name        string
		ctx         context.Context
		frames      map[string][]byte
		numOfFrames int32
		shouldStart bool
	}{
		{
			name:        "nil context",
			ctx:         nil,
			frames:      map[string][]byte{"frame1": []byte("data")},
			numOfFrames: 1,
			shouldStart: false,
		},
		{
			name:        "nil frames",
			ctx:         context.Background(),
			frames:      nil,
			numOfFrames: 1,
			shouldStart: false,
		},
		{
			name:        "empty frames",
			ctx:         context.Background(),
			frames:      map[string][]byte{},
			numOfFrames: 0,
			shouldStart: false,
		},
		{
			name: "valid inputs",
			ctx:  context.Background(),
			frames: map[string][]byte{
				"frame1": []byte("data1"),
				"frame2": []byte("data2"),
			},
			numOfFrames: 2,
			shouldStart: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := New(tt.ctx, tt.frames)

			am.Start()

			time.Sleep(50 * time.Millisecond)

			state := am.state.Load()
			if tt.shouldStart {
				assert.NotNil(t, state, "animation should have started")
				am.Stop()
			} else {
				assert.Nil(t, state, "animation should not have started")
			}
		})
	}
}

func TestAnimator_Start_ConcurrentCalls(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}

	am := New(ctx, frames)

	am.Start()
	time.Sleep(50 * time.Millisecond)

	firstState := am.state.Load()
	require.NotNil(t, firstState, "first animation should start")

	am.Start()
	time.Sleep(50 * time.Millisecond)

	secondState := am.state.Load()
	assert.Equal(t, firstState, secondState, "state should not change on second Start()")

	am.Stop()
}

func TestAnimator_Stop_NonRunning(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{"frame1": []byte("data")}
	am := New(ctx, frames)

	am.Stop()
	assert.Nil(t, am.state.Load())
}

func TestAnimator_Stop_Running(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}
	am := New(ctx, frames)

	am.Start()
	time.Sleep(50 * time.Millisecond)

	require.NotNil(t, am.state.Load(), "animation should be running")

	am.Stop()
	time.Sleep(50 * time.Millisecond)

	assert.Nil(t, am.state.Load(), "animation should be stopped")
}

func TestAnimator_Stop_Multiple(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{"frame1": []byte("data")}
	am := New(ctx, frames)

	am.Start()
	time.Sleep(50 * time.Millisecond)

	am.Stop()
	am.Stop()
	am.Stop()

	assert.Nil(t, am.state.Load())
}

func TestAnimator_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}

	am := New(ctx, frames)
	am.Start()
	time.Sleep(50 * time.Millisecond)

	require.NotNil(t, am.state.Load(), "animation should be running")

	cancel()
	time.Sleep(100 * time.Millisecond)

	am.Stop()
}

func TestAnimator_FrameIndexing(t *testing.T) {
	tests := []struct {
		name        string
		numFrames   int32
		iterations  int
		wantIndices []int32
	}{
		{
			name:        "2 frames, 5 iterations",
			numFrames:   2,
			iterations:  5,
			wantIndices: []int32{1, 2, 1, 2, 1},
		},
		{
			name:        "3 frames, 7 iterations",
			numFrames:   3,
			iterations:  7,
			wantIndices: []int32{1, 2, 3, 1, 2, 3, 1},
		},
		{
			name:        "1 frame, 3 iterations",
			numFrames:   1,
			iterations:  3,
			wantIndices: []int32{1, 1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &Animator{
				numOfFrames: tt.numFrames,
			}

			var gotIndices []int32
			for i := 0; i < tt.iterations; i++ {
				idx := (am.animationIdx.Load() % am.numOfFrames) + 1
				gotIndices = append(gotIndices, idx)
				am.animationIdx.Store(idx)
			}

			assert.Equal(t, tt.wantIndices, gotIndices)
		})
	}
}

func TestAnimator_StartStopCycle(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}

	am := New(ctx, frames)

	for cycle := 1; cycle <= 3; cycle++ {
		am.Start()
		time.Sleep(50 * time.Millisecond)
		assert.NotNil(t, am.state.Load(), "cycle %d: start should work", cycle)

		am.Stop()
		time.Sleep(50 * time.Millisecond)
		assert.Nil(t, am.state.Load(), "cycle %d: stop should work", cycle)
	}
}

func TestAnimator_FrameKeyCalculation(t *testing.T) {
	tests := []struct {
		name         string
		currentIdx   int32
		numFrames    int32
		wantFrameKey string
	}{
		{
			name:         "first frame",
			currentIdx:   0,
			numFrames:    3,
			wantFrameKey: "frame1",
		},
		{
			name:         "second frame",
			currentIdx:   1,
			numFrames:    3,
			wantFrameKey: "frame2",
		},
		{
			name:         "wrap around",
			currentIdx:   2,
			numFrames:    3,
			wantFrameKey: "frame3",
		},
		{
			name:         "single frame",
			currentIdx:   0,
			numFrames:    1,
			wantFrameKey: "frame1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &Animator{
				numOfFrames: tt.numFrames,
			}
			am.animationIdx.Store(tt.currentIdx)

			idx := (am.animationIdx.Load() % am.numOfFrames) + 1
			frameKey := "frame" + string(rune('0'+idx))

			assert.Equal(t, tt.wantFrameKey, frameKey)
		})
	}
}

func TestAnimator_StateAtomicity(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}

	am := New(ctx, frames)

	am.Start()
	time.Sleep(50 * time.Millisecond)

	state1 := am.state.Load()
	require.NotNil(t, state1)

	state2 := am.state.Load()
	assert.Equal(t, state1, state2, "state should remain consistent")

	am.Stop()

	state3 := am.state.Load()
	assert.Nil(t, state3, "state should be nil after stop")
}

func TestAnimator_FramesMapReference(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}

	am := New(ctx, frames)

	assert.Equal(t, frames, am.frames)
	assert.Equal(t, len(frames), int(am.numOfFrames))

	frames["frame3"] = []byte("data3")
	assert.Equal(t, 3, len(am.frames), "manager should reference same map")
}

func TestAnimator_ContextReference(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{"frame1": []byte("data")}

	am := New(ctx, frames)

	assert.Equal(t, ctx, am.ctx)
}

func TestAnimator_NumOfFramesCalculation(t *testing.T) {
	tests := []struct {
		name       string
		frameCount int
	}{
		{"0 frames", 0},
		{"1 frame", 1},
		{"5 frames", 5},
		{"100 frames", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frames := make(map[string][]byte)
			for i := 1; i <= tt.frameCount; i++ {
				frames["frame"+string(rune('0'+i))] = []byte("data")
			}

			am := New(context.Background(), frames)

			assert.Equal(t, int32(tt.frameCount), am.numOfFrames)
		})
	}
}

func TestAnimator_FrameRateZeroDefaultsToDefault(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}

	am := New(ctx, frames)

	assert.Equal(t, DefaultFrameRate, am.frameRate)
}

func TestAnimator_FrameRateCustom(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}
	customRate := 50 * time.Millisecond

	am := New(ctx, frames, customRate)

	assert.Equal(t, customRate, am.frameRate)
}

func TestAnimator_FrameRateMultipleArgsUsesFirst(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}
	firstRate := 50 * time.Millisecond
	secondRate := 100 * time.Millisecond

	am := New(ctx, frames, firstRate, secondRate)

	assert.Equal(t, firstRate, am.frameRate)
}

func TestAnimator_FrameRateFast(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}
	fastRate := 50 * time.Millisecond

	am := New(ctx, frames, fastRate)
	am.Start()
	time.Sleep(150 * time.Millisecond)

	require.NotNil(t, am.state.Load(), "animation should be running")
	assert.Greater(t, am.animationIdx.Load(), int32(0), "should have advanced frames with fast rate")

	am.Stop()
}

func TestAnimator_FrameRateSlow(t *testing.T) {
	ctx := context.Background()
	frames := map[string][]byte{
		"frame1": []byte("data1"),
		"frame2": []byte("data2"),
	}
	slowRate := 1 * time.Second

	am := New(ctx, frames, slowRate)
	am.Start()
	time.Sleep(100 * time.Millisecond)

	require.NotNil(t, am.state.Load(), "animation should be running")
	assert.Equal(t, int32(0), am.animationIdx.Load(), "should not have advanced frames with slow rate")

	am.Stop()
}
