// =============================================================================
// PATTERN 04: THE DONE CHANNEL -- Basic Goroutine Cancellation
// =============================================================================
//
// This program demonstrates the fundamental done channel pattern: using a
// closed channel to signal one or more goroutines to stop their work.
//
// WHAT THIS EXAMPLE SHOWS:
// ------------------------
// 1. Creating a done channel (chan struct{}) for signaling
// 2. Launching a worker goroutine that checks for cancellation via select
// 3. Closing the done channel to broadcast a stop signal
// 4. The zero-value trick: close(done) unblocks all receivers simultaneously
// 5. Waiting for the worker to acknowledge cancellation
//
// WHAT THIS EXAMPLE INTENTIONALLY DOES NOT SHOW:
// -----------------------------------------------
// - Cleanup or graceful shutdown (see 4.2-done-channel-cleanup.go)
// - context.Context cancellation (covered in the Context pattern)
// - Multiple done channels with select (covered in the Select pattern)
// - Error propagation during cancellation (covered in the Errgroup pattern)
//
// RUN: go run 4.1-done-channel-basic.go
//
// =============================================================================

package main

import (
	"fmt"
	"time"
)

// =============================================================================
// WORKER FUNCTION
// =============================================================================

// worker simulates a long-running goroutine that performs periodic work.
//
// It checks for cancellation on every iteration using select. When the
// done channel is closed, the <-done case becomes immediately receivable
// (returning the zero value of struct{}), and the goroutine exits.
//
// The finished channel signals back to the caller that the worker has
// fully stopped and returned.
func worker(done <-chan struct{}, finished chan<- struct{}) {
	defer close(finished) // Signal that this goroutine has returned

	iteration := 0
	for {
		select {
		case <-done:
			// The done channel was closed. Time to stop.
			fmt.Println("  Worker: received done signal, exiting")
			return
		default:
			// No signal yet, keep working.
		}

		// Simulate a unit of work.
		iteration++
		fmt.Printf("  Worker: iteration %d\n", iteration)
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	fmt.Println()
	fmt.Println("=== Example 4.1: Basic Done Channel ===")
	fmt.Println("Using close(done) to signal a goroutine to stop.")
	fmt.Println()

	// =========================================================================
	// SETUP: CREATE THE SIGNALING CHANNELS
	// =========================================================================
	//
	// done: used to broadcast cancellation to the worker.
	//       We never send values through it; we only close it.
	//       chan struct{} costs zero bytes per element, making
	//       it ideal for pure signaling.
	//
	// finished: used by the worker to signal it has fully stopped.
	//           This lets main() wait for the worker to exit cleanly.

	done := make(chan struct{})     // Signal channel: close to cancel
	finished := make(chan struct{}) // Acknowledgment: worker closes when done

	// =========================================================================
	// LAUNCH THE WORKER
	// =========================================================================
	//
	// The worker receives done as <-chan struct{} (receive-only) so it
	// cannot accidentally close the done channel itself. Only the
	// owner (main) controls the done channel's lifecycle.

	go worker(done, finished)

	// =========================================================================
	// LET THE WORKER RUN, THEN CANCEL
	// =========================================================================
	//
	// Timeline:
	//   0ms        worker starts, iterates every 100ms
	//   350ms      main closes done channel
	//   350ms+     worker sees <-done, prints exit message, closes finished
	//   350ms+     main receives from finished, program ends
	//
	// The close(done) call unblocks ALL goroutines waiting on <-done.
	// Even though we have only one worker here, the same close would
	// cancel 1,000 workers simultaneously. This is the broadcast
	// property of channel closure.

	time.Sleep(350 * time.Millisecond)
	fmt.Println("Main: sending done signal (closing channel)...")
	close(done)

	// =========================================================================
	// WAIT FOR WORKER TO FINISH
	// =========================================================================
	//
	// We block until the worker closes the finished channel.
	// This ensures the worker has fully exited before main returns.

	<-finished
	fmt.Println("Main: worker has exited")

	fmt.Println()
	fmt.Println("=== KEY INSIGHT ===")
	fmt.Println("Closing a channel unblocks all goroutines waiting on it. This")
	fmt.Println("broadcast property makes close(done) the idiomatic way to signal")
	fmt.Println("cancellation to any number of goroutines simultaneously.")
	fmt.Println()
}
