// Copyright 2019 The Oto Authors / 2025 Makoto Onuki
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/ebitengine/oto/v3"
)

var (
	sampleRate   = flag.Int("samplerate", 48000, "sample rate")
	channelCount = flag.Int("channelcount", 1, "number of channel")
	format       = flag.String("format", "s16le", "source format (u8, s16le, or f32le)")
	freq         = flag.Float64("freq", 880, "frequency in Hz")
	duration     = flag.Int("duration", 180, "duration in ms")
	amp          = flag.Float64("amp", 0.3, "amplitude (0.0-1.0)")
)

type SineWave struct {
	freq   float64
	amp    float64
	length int64
	pos    int64

	channelCount int
	format       oto.Format

	remaining []byte
}

func formatByteLength(format oto.Format) int {
	switch format {
	case oto.FormatFloat32LE:
		return 4
	case oto.FormatUnsignedInt8:
		return 1
	case oto.FormatSignedInt16LE:
		return 2
	default:
		panic(fmt.Sprintf("unexpected format: %d", format))
	}
}

func NewSineWave(freq, amp float64, duration time.Duration, channelCount int, format oto.Format) *SineWave {
	l := int64(channelCount) * int64(formatByteLength(format)) * int64(*sampleRate) * int64(duration) / int64(time.Second)
	l = l / 4 * 4
	return &SineWave{
		freq:         freq,
		amp:          amp,
		length:       l,
		channelCount: channelCount,
		format:       format,
	}
}

func (s *SineWave) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		copy(s.remaining, s.remaining[n:])
		s.remaining = s.remaining[:len(s.remaining)-n]
		return n, nil
	}

	if s.pos == s.length {
		return 0, io.EOF
	}

	eof := false
	if s.pos+int64(len(buf)) > s.length {
		buf = buf[:s.length-s.pos]
		eof = true
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := float64(*sampleRate) / float64(s.freq)

	num := formatByteLength(s.format) * s.channelCount
	p := s.pos / int64(num)
	switch s.format {
	case oto.FormatFloat32LE:
		for i := 0; i < len(buf)/num; i++ {
			bs := math.Float32bits(float32(math.Sin(2*math.Pi*float64(p)/length) * s.amp))
			for ch := 0; ch < *channelCount; ch++ {
				buf[num*i+4*ch] = byte(bs)
				buf[num*i+1+4*ch] = byte(bs >> 8)
				buf[num*i+2+4*ch] = byte(bs >> 16)
				buf[num*i+3+4*ch] = byte(bs >> 24)
			}
			p++
		}
	case oto.FormatUnsignedInt8:
		for i := 0; i < len(buf)/num; i++ {
			const max = 127
			b := int(math.Sin(2*math.Pi*float64(p)/length) * s.amp * max)
			for ch := 0; ch < *channelCount; ch++ {
				buf[num*i+ch] = byte(b + 128)
			}
			p++
		}
	case oto.FormatSignedInt16LE:
		for i := 0; i < len(buf)/num; i++ {
			const max = 32767
			b := int16(math.Sin(2*math.Pi*float64(p)/length) * s.amp * max)
			for ch := 0; ch < *channelCount; ch++ {
				buf[num*i+2*ch] = byte(b)
				buf[num*i+1+2*ch] = byte(b >> 8)
			}
			p++
		}
	}

	s.pos += int64(len(buf))

	n := len(buf)
	if origBuf != nil {
		n = copy(origBuf, buf)
		s.remaining = buf[n:]
	}

	if eof {
		return n, io.EOF
	}
	return n, nil
}

func play(context *oto.Context, freq, amp float64, duration time.Duration, channelCount int, format oto.Format) *oto.Player {
	p := context.NewPlayer(NewSineWave(freq, amp, duration, channelCount, format))
	p.Play()
	return p
}

func run() error {
	op := &oto.NewContextOptions{}
	op.SampleRate = *sampleRate
	op.ChannelCount = *channelCount

	switch *format {
	case "f32le":
		op.Format = oto.FormatFloat32LE
	case "u8":
		op.Format = oto.FormatUnsignedInt8
	case "s16le":
		op.Format = oto.FormatSignedInt16LE
	default:
		return fmt.Errorf("format must be u8, s16le, or f32le but: %s", *format)
	}
	c, ready, err := oto.NewContext(op)
	if err != nil {
		return err
	}
	<-ready

	duration := time.Duration(*duration) * time.Millisecond
	play(c, *freq, *amp, duration, op.ChannelCount, op.Format)
	time.Sleep(duration)

	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}
