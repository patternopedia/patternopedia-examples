// =============================================================================
// PATTERN 03: THE BUFFERED CHANNEL — Buffer Capacity and Blocking
// =============================================================================
//
// This program demonstrates the core property of buffered channels: sends
// proceed without blocking until the buffer is full, and receives proceed
// without blocking until the buffer is empty.
//
// WHAT THIS EXAMPLE SHOWS:
// ------------------------
// 1. Creating a buffered channel with make(chan T, capacity)
// 2. Sends do not block while the buffer has room
// 3. The buffer fills to capacity, at which point the next send would block
// 4. Draining the buffer with sequential receives
// 5. FIFO ordering guarantee (values come out in the order they went in)
//
// WHAT THIS EXAMPLE INTENTIONALLY DOES NOT SHOW:
// -----------------------------------------------
// - Channel closure and draining (see 3.2-buffered-channel-close-drain.go)
// - Producer-consumer with goroutines (covered in Pipeline pattern)
// - Select with buffered channels (covered in the Select pattern)
// - Buffered channels as semaphores (covered in the Semaphore pattern)
//
// RUN: go run 3.1-buffered-channel-basics.go
//
// =============================================================================

package main

import "fmt"

func main() {
	fmt.Println()
	fmt.Println("=== Example 3.1: Buffer Capacity and Blocking ===")
	fmt.Println("Sends proceed without blocking until the buffer is full.")
	fmt.Println()

	// Create a buffered channel with capacity 5.
	// This channel can hold up to 5 values before any send blocks.
	ch := make(chan int, 5)

	// Phase 1: Fill the buffer without blocking.
	// Each send succeeds immediately because the buffer has room.
	for i := 1; i <= 5; i++ {
		ch <- i
		fmt.Printf("Sent %d: len=%d cap=%d\n", i, len(ch), cap(ch))
	}

	fmt.Println()
	fmt.Println("Buffer is now FULL (len == cap).")
	fmt.Println("The next send would BLOCK until a receiver drains a value.")
	fmt.Println()

	// Phase 2: Drain the buffer.
	// Each receive succeeds immediately because the buffer has values.
	for i := 1; i <= 5; i++ {
		val := <-ch
		fmt.Printf("Received %d: len=%d cap=%d\n", val, len(ch), cap(ch))
	}

	fmt.Println()
	fmt.Println("Buffer is now EMPTY (len == 0).")
	fmt.Println("The next receive would BLOCK until a sender provides a value.")

	// Verify FIFO ordering: values came out in the order they went in.
	ch2 := make(chan string, 3)
	ch2 <- "first"
	ch2 <- "second"
	ch2 <- "third"

	fmt.Println()
	fmt.Println("FIFO verification:")
	fmt.Println(<-ch2) // first
	fmt.Println(<-ch2) // second
	fmt.Println(<-ch2) // third

	fmt.Println()
	fmt.Println("=== KEY INSIGHT ===")
	fmt.Println("Buffered channels decouple send and receive timing. Sends do not")
	fmt.Println("block until the buffer is full, and receives do not block until")
	fmt.Println("the buffer is empty. The buffer absorbs temporary rate mismatches.")
	fmt.Println()
}
