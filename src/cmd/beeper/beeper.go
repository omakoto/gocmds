// Just play an 880 hz short beep sound.
package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// The //go:embed directive MUST be placed immediately above the variable declaration.
// We embed the file as a raw byte slice ([]byte).
//
//go:embed beep.wav
var embeddedWavFile []byte

func main() {
	// 1. Create a reader for the embedded byte slice
	// We use bytes.NewReader to treat the []byte as an io.ReadSeeker,
	// which is what the wav.Decode function expects.
	wavReader := bytes.NewReader(embeddedWavFile)

	// 2. Decode the WAV file
	streamer, format, err := wav.Decode(wavReader)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// 3. Initialize the speaker
	// The second argument is a buffer size.
	// This example uses the format's sample rate divided by 10 as a buffer size.
	// You may need to adjust this depending on your system for better performance.
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}

	// 4. Play the sound
	fmt.Println("Playing 'test.wav'...")

	speaker.Play(streamer)

	// We'll just wait for the streamer to finish manually.
	// In a real application, you might use a WaitGroup or a separate goroutine
	// to manage playback state.
	for streamer.Position() < streamer.Len() {
		time.Sleep(50)
	}

	fmt.Println("Playback finished.")

	// Terminate the speaker after playback is finished
	// The speaker is not thread-safe, so you must use speaker.Lock/Unlock
	// to ensure it doesn't try to access the audio device after termination.
	speaker.Close()
}
