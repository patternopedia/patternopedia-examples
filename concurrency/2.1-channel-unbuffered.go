// =============================================================================
// PATTERN 02: THE CHANNEL — Unbuffered Channel Communication
// =============================================================================
//
// This program demonstrates the fundamental channel operation: two goroutines
// communicating through an unbuffered channel.
//
// WHAT THIS EXAMPLE SHOWS:
// ------------------------
// 1. Creating an unbuffered channel with make(chan T)
// 2. The send operation (ch <- value) blocks until a receiver is ready
// 3. The receive operation (<-ch) unblocks the sender — the "rendezvous"
// 4. Neither side can proceed until both are ready
//
// WHAT THIS EXAMPLE INTENTIONALLY DOES NOT SHOW:
// -----------------------------------------------
// - Buffered channels (covered in this pattern's buffered section)
// - Channel closure and range (see channel-closure-range.go)
// - Select statement (covered in the Select pattern)
// - Context cancellation (covered in the Context pattern)
//
// RUN: go run 2.1-channel-unbuffered.go
//
// =============================================================================

package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println()
	fmt.Println("=== Example 2.1: Unbuffered Channel Communication ===")
	fmt.Println("Two goroutines communicating through an unbuffered channel.")
	fmt.Println()

	// Create an unbuffered channel for string values.
	// Unbuffered means capacity is 0: no storage, just a handoff point.
	ch := make(chan string)

	// Launch a sender goroutine that will try to send a message.
	// This runs concurrently with main().
	go func() {
		fmt.Println("Sender: about to send...")

		// SEND OPERATION: This line BLOCKS until main() executes <-ch.
		// The sender cannot proceed past this point until a receiver arrives.
		ch <- "Hello from sender!"

		// This only prints AFTER the receiver has taken the value.
		fmt.Println("Sender: message sent!")
	}()

	// Sleep to let the sender goroutine reach the blocking send.
	// This demonstrates that the sender is stuck waiting for us.
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Receiver: about to receive...")

	// RECEIVE OPERATION: This unblocks the sender.
	// At this moment, the rendezvous occurs: value transfers, both proceed.
	message := <-ch
	fmt.Println("Receiver: got", message)

	// Brief sleep to let the sender's final print complete before main exits.
	time.Sleep(50 * time.Millisecond)

	fmt.Println()
	fmt.Println("=== KEY INSIGHT ===")
	fmt.Println("With unbuffered channels, neither side can proceed until both are")
	fmt.Println("ready. The sender blocks on ch <- value until a receiver executes")
	fmt.Println("<-ch. This creates a synchronization point: a rendezvous between goroutines.")
	fmt.Println()
}
