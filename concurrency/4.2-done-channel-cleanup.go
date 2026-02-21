// =============================================================================
// PATTERN 04: THE DONE CHANNEL -- Graceful Shutdown with Cleanup
// =============================================================================
//
// This program demonstrates the done channel pattern with graceful shutdown:
// when cancellation is signaled, the worker finishes its current unit of
// work and performs cleanup before exiting.
//
// WHAT THIS EXAMPLE SHOWS:
// ------------------------
// 1. A worker that processes items from a work channel
// 2. Using select to multiplex between work and cancellation
// 3. Graceful shutdown: finishing current work before exiting
// 4. Cleanup logic (resource release) after the cancellation loop exits
// 5. Coordinating shutdown across multiple goroutines with a single close
//
// WHAT THIS EXAMPLE INTENTIONALLY DOES NOT SHOW:
// -----------------------------------------------
// - context.Context for cancellation (covered in the Context pattern)
// - Timeout-based cancellation (covered in the Context pattern)
// - Error propagation during shutdown (covered in the Errgroup pattern)
// - Draining the work channel after cancellation (covered in Pipeline)
//
// RUN: go run 4.2-done-channel-cleanup.go
//
// =============================================================================

package main

import (
	"fmt"
	"time"
)

// =============================================================================
// WORKER WITH CLEANUP
// =============================================================================

// gracefulWorker processes items from a work channel until it receives a
// done signal. On cancellation, it finishes processing the current item
// (if any), performs cleanup, and signals completion via the finished channel.
//
// This pattern is common in production systems where workers hold resources
// (database connections, file handles, network sockets) that must be released
// cleanly on shutdown.
func gracefulWorker(id int, work <-chan string, done <-chan struct{}, finished chan<- struct{}) {
	defer close(finished) // Signal that this goroutine has fully exited

	processed := 0

	for {
		select {
		case <-done:
			// =========================================================
			// CLEANUP PHASE
			// =========================================================
			//
			// The done channel is closed. We exit the work loop and
			// perform any necessary cleanup. In a real system, this
			// would close database connections, flush buffers, release
			// file handles, or deregister from a service registry.
			fmt.Printf("  Worker %d: done signal received, cleaning up (%d items processed)\n", id, processed)
			fmt.Printf("  Worker %d: releasing resources...\n", id)
			time.Sleep(50 * time.Millisecond) // Simulate cleanup work
			fmt.Printf("  Worker %d: cleanup complete\n", id)
			return

		case item, ok := <-work:
			if !ok {
				// Work channel was closed (no more items).
				fmt.Printf("  Worker %d: work channel closed, exiting (%d items processed)\n", id, processed)
				return
			}
			// Process the work item.
			fmt.Printf("  Worker %d: processing %q\n", id, item)
			time.Sleep(80 * time.Millisecond) // Simulate processing
			processed++
		}
	}
}

func main() {
	fmt.Println()
	fmt.Println("=== Example 4.2: Done Channel with Cleanup ===")
	fmt.Println("Graceful shutdown - workers finish current work and clean up resources.")
	fmt.Println()

	// =========================================================================
	// SETUP: CHANNELS FOR WORK, CANCELLATION, AND COMPLETION
	// =========================================================================
	//
	// work:      buffered channel carrying work items to the workers
	// done:      signal channel; close to cancel all workers
	// finished1: worker 1 closes this when fully exited
	// finished2: worker 2 closes this when fully exited

	work := make(chan string, 5)
	done := make(chan struct{})
	finished1 := make(chan struct{})
	finished2 := make(chan struct{})

	// =========================================================================
	// LAUNCH TWO WORKERS
	// =========================================================================
	//
	// Both workers share the same done channel. When we close it, BOTH
	// workers receive the signal simultaneously. This is the broadcast
	// property of channel closure: one close, N receivers.

	go gracefulWorker(1, work, done, finished1)
	go gracefulWorker(2, work, done, finished2)

	// =========================================================================
	// SEND WORK ITEMS
	// =========================================================================
	//
	// We send several items into the work channel. Workers pick them
	// up concurrently. After a brief delay, we signal cancellation
	// before all items are necessarily processed.

	items := []string{"task-A", "task-B", "task-C", "task-D", "task-E", "task-F"}
	for _, item := range items {
		work <- item
	}
	fmt.Println("Main: all work items sent")

	// =========================================================================
	// CANCEL AFTER A DELAY
	// =========================================================================
	//
	// Timeline:
	//   0ms        workers start processing (80ms per item)
	//   ~300ms     main closes done channel
	//   ~300ms+    workers see <-done in their select, enter cleanup
	//   ~350ms+    workers finish cleanup, close finished channels
	//   ~350ms+    main unblocks, program exits

	time.Sleep(300 * time.Millisecond)
	fmt.Println()
	fmt.Println("Main: signaling shutdown (closing done channel)...")
	close(done)

	// =========================================================================
	// WAIT FOR BOTH WORKERS TO FINISH CLEANUP
	// =========================================================================
	//
	// We wait on each worker's finished channel. This guarantees both
	// workers have completed their cleanup before main exits.

	<-finished1
	<-finished2
	fmt.Println()
	fmt.Println("Main: all workers have shut down cleanly")

	fmt.Println()
	fmt.Println("=== KEY INSIGHT ===")
	fmt.Println("The done channel pattern enables graceful shutdown. Workers check")
	fmt.Println("for cancellation alongside their normal work using select, finish")
	fmt.Println("their current task, release resources, and signal completion. One")
	fmt.Println("close(done) cancels all workers simultaneously.")
	fmt.Println()
}
