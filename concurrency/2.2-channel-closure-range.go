// =============================================================================
// PATTERN 02: THE CHANNEL — Channel Closure and Range
// =============================================================================
//
// This program demonstrates closing a channel to signal completion and using
// range to consume all values until the channel closes.
//
// WHAT THIS EXAMPLE SHOWS:
// ------------------------
// 1. The producer-consumer pattern with channels
// 2. Closing a channel with close(ch) to signal "no more values"
// 3. Using range to receive all values until channel closes
// 4. The producer owns the channel and is responsible for closing it
//
// WHAT THIS EXAMPLE INTENTIONALLY DOES NOT SHOW:
// -----------------------------------------------
// - Unbuffered channel mechanics (see channel-unbuffered.go)
// - The comma-ok idiom for manual close detection (covered in the pattern page)
// - Buffered channels (covered in this pattern's buffered section)
// - Select statement (covered in the Select pattern)
//
// RUN: go run 2.2-channel-closure-range.go
//
// =============================================================================

package main

import "fmt"

func main() {
	fmt.Println()
	fmt.Println("=== Example 2.2: Channel Closure and Range ===")
	fmt.Println("Closing a channel to signal completion; range to consume all values.")
	fmt.Println()

	// Create an unbuffered channel for the producer-consumer pattern.
	ch := make(chan string)

	// PRODUCER: Runs in a separate goroutine, sends values, then closes.
	// The producer "owns" the channel and is responsible for closing it.
	go func() {
		words := []string{"Go", "channels", "are", "first-class", "values!"}

		// Send each word through the channel.
		// Each send blocks until the consumer receives.
		for _, word := range words {
			ch <- word
		}

		// CLOSE: Signal that no more values will be sent.
		// This is how the producer tells the consumer "I'm done."
		// After this, the consumer's range loop will exit.
		close(ch)
	}()

	// CONSUMER: Uses range to receive all values until channel closes.
	// The range loop automatically handles the comma-ok check internally.
	// When the channel closes and drains, the loop exits cleanly.
	for word := range ch {
		fmt.Println("Received:", word)
	}

	// We only reach this line after the channel is closed and drained.
	fmt.Println("Channel closed, loop exited")

	fmt.Println()
	fmt.Println("=== KEY INSIGHT ===")
	fmt.Println("The range construct handles the comma-ok check internally, exiting")
	fmt.Println("when the channel is closed and drained. This is the idiomatic way")
	fmt.Println("to consume all values from a channel.")
	fmt.Println()
}
