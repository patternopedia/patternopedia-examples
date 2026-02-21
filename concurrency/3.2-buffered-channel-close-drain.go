// =============================================================================
// PATTERN 03: THE BUFFERED CHANNEL — Close and Drain
// =============================================================================
//
// This program demonstrates a behavior unique to buffered channels: closing a
// channel does not discard buffered values. All values in the buffer at the
// time of closure remain accessible and can be drained by receivers.
//
// WHAT THIS EXAMPLE SHOWS:
// ------------------------
// 1. Closing a buffered channel while values remain in the buffer
// 2. Buffered values survive closure and can be drained with range
// 3. After draining, receives return zero values with ok=false
// 4. The comma-ok idiom to distinguish real zero values from closure
//
// WHAT THIS EXAMPLE INTENTIONALLY DOES NOT SHOW:
// -----------------------------------------------
// - Sending to a closed channel (panics; covered in the Channel pattern)
// - Goroutine-based producer-consumer (covered in Pipeline pattern)
// - Graceful shutdown with context (covered in the Context pattern)
// - Select with closed channels (covered in the Select pattern)
//
// RUN: go run 3.2-buffered-channel-close-drain.go
//
// =============================================================================

package main

import "fmt"

func main() {
	fmt.Println()
	fmt.Println("=== Example 3.2: Close and Drain ===")
	fmt.Println("Buffered values survive channel closure and can be drained.")
	fmt.Println()

	// Create a buffered channel and load it with values.
	ch := make(chan int, 5)
	ch <- 10
	ch <- 20
	ch <- 30

	fmt.Printf("Before close: len=%d cap=%d\n", len(ch), cap(ch))

	// Close the channel while values are still in the buffer.
	// This signals "no more sends" but does NOT discard buffered values.
	close(ch)
	fmt.Println("Channel closed.")
	fmt.Printf("After close:  len=%d cap=%d\n", len(ch), cap(ch))
	fmt.Println()

	// Drain all buffered values using range.
	// The range loop receives each buffered value in FIFO order,
	// then exits when the channel is both closed and empty.
	fmt.Println("Draining with range:")
	for val := range ch {
		fmt.Printf("  Received: %d\n", val)
	}
	fmt.Println("Range loop exited (channel closed and drained).")
	fmt.Println()

	// After draining, receives return zero values with ok=false.
	val, ok := <-ch
	fmt.Printf("After drain: val=%d ok=%v\n", val, ok)
	fmt.Println()

	// Demonstrate the comma-ok idiom to distinguish
	// a real zero value from a closed-channel zero value.
	ch2 := make(chan int, 3)
	ch2 <- 0  // A real zero!
	ch2 <- 42
	close(ch2)

	fmt.Println("Distinguishing real values from closure:")
	for {
		v, open := <-ch2
		if !open {
			fmt.Printf("  val=%d ok=%v (channel closed)\n", v, open)
			break
		}
		fmt.Printf("  val=%d ok=%v (real value)\n", v, open)
	}

	fmt.Println()
	fmt.Println("=== KEY INSIGHT ===")
	fmt.Println("Closing a buffered channel does not discard its contents. All")
	fmt.Println("buffered values survive closure and can be drained with range")
	fmt.Println("or comma-ok. This enables graceful shutdown: close to signal")
	fmt.Println("completion, then let consumers process remaining work.")
	fmt.Println()
}
