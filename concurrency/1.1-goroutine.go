// =============================================================================
// PATTERN 01: THE GOROUTINE
// =============================================================================
//
// This program demonstrates the fundamental goroutine pattern: spawning
// concurrent execution with the `go` keyword.
//
// WHAT THIS EXAMPLE SHOWS:
// ------------------------
// 1. Sequential execution: calling functions directly (blocks until complete)
// 2. Concurrent execution: using `go` to spawn goroutines (returns immediately)
// 3. The fundamental problem: main() exits before goroutines complete
//
// WHAT THIS EXAMPLE INTENTIONALLY DOES NOT SHOW:
// ----------------------------------------------
// - Channels (covered in the Channel pattern)
// - WaitGroup (covered in the Barrier/WaitGroup pattern)
// - Mutexes (covered in synchronization patterns)
//
// The `time.Sleep()` at the end is INTENTIONALLY fragile. It demonstrates
// that spawning goroutines is easy, but coordinating them requires the
// patterns that follow.
//
// RUN: go run 1-goroutine.go
//
// =============================================================================

package main

import (
	"fmt"
	"time"
)

// =============================================================================
// SIMULATING WORK
// =============================================================================

// simulateWork represents an I/O-bound operation like an API call or DB query.
//
// In real applications, this function body would contain actual I/O:
//
//	resp, err := http.Get(url)           // Network request
//	rows, err := db.Query(sql)           // Database query
//	data, err := os.ReadFile(path)       // File read
//
// We use time.Sleep() to simulate the waiting that happens during I/O.
// The key insight: during this "sleep", a goroutine consumes almost no
// resources - it's parked by the Go runtime, and the OS thread is free
// to run other goroutines.
func simulateWork(name string, duration time.Duration) {
	fmt.Printf("  Starting: %s\n", name)
	time.Sleep(duration) // Simulates waiting for I/O
	fmt.Printf("  Finished: %s (took %v)\n", name, duration)
}

// =============================================================================
// TASK DEFINITION
// =============================================================================

// task holds information about a unit of work.
// Using a simple struct keeps the example focused on goroutines, not data structures.
type task struct {
	name     string
	duration time.Duration
}

func main() {
	fmt.Println()
	fmt.Println("=== Example 1: Sequential vs Concurrent Execution ===")
	fmt.Println("Comparing four simulated I/O-bound tasks run sequentially vs concurrently.")
	fmt.Println()

	// Define our tasks - these simulate real-world I/O operations
	// with realistic latencies (API calls, DB queries, etc.)
	tasks := []task{
		{"API Call", 400 * time.Millisecond},
		{"DB Query", 350 * time.Millisecond},
		{"File I/O", 300 * time.Millisecond},
		{"Cache Check", 200 * time.Millisecond},
	}

	// =========================================================================
	// SEQUENTIAL EXECUTION
	// =========================================================================
	//
	// Each function call BLOCKS until it completes. The next task cannot
	// start until the previous one finishes.
	//
	// Timeline:
	//   |---API Call---|---DB Query---|---File I/O---|---Cache---|
	//   0             400            750           1050        1250ms
	//
	// Total time = sum of all durations

	fmt.Println("=== SEQUENTIAL EXECUTION ===")
	fmt.Println("(Each task waits for the previous to complete)")
	fmt.Println()

	seqStart := time.Now()

	// Execute tasks one at a time - each call blocks until complete
	for _, t := range tasks {
		simulateWork(t.name, t.duration) // This BLOCKS until the task finishes
	}

	seqDuration := time.Since(seqStart)
	fmt.Printf("\nSequential total: %v\n", seqDuration.Round(time.Millisecond))

	// =========================================================================
	// CONCURRENT EXECUTION (with goroutines)
	// =========================================================================
	//
	// The `go` keyword spawns a new goroutine - a lightweight thread managed
	// by the Go runtime. The key difference:
	//
	//   simulateWork(...)      // BLOCKS - waits for completion
	//   go simulateWork(...)   // RETURNS IMMEDIATELY - work happens concurrently
	//
	// Timeline:
	//   |---API Call (400ms)---|
	//   |--DB Query (350ms)--|
	//   |-File I/O (300ms)-|
	//   |-Cache (200ms)-|
	//   0                    400ms
	//
	// Total time ~= longest task (400ms), not sum (1250ms)

	fmt.Println("\n=== CONCURRENT EXECUTION ===")
	fmt.Println("(All tasks start simultaneously)")
	fmt.Println()

	concStart := time.Now()

	// THE GOROUTINE PATTERN:
	// ----------------------
	// Spawn a goroutine for each task. Each `go` statement returns
	// immediately - the loop completes in microseconds, not seconds.
	for _, t := range tasks {
		// Why the anonymous function with parameters?
		// -------------------------------------------
		// We pass t.name and t.duration as parameters to capture their
		// VALUES at this iteration. Without this, all goroutines might
		// see the final loop values (though Go 1.22+ fixes this for
		// loop variables, the parameter style is still clearer).
		go func(name string, duration time.Duration) {
			simulateWork(name, duration)
		}(t.name, t.duration) // not strictly necessary in Go 1.22+, but clearer
	}

	// THE FUNDAMENTAL PROBLEM:
	// ------------------------
	// The loop above completes almost instantly because `go` returns
	// immediately. But our goroutines are still running! If main()
	// exits now, all goroutines die immediately - they won't finish
	// their work.
	//
	// We MUST wait somehow. Here we use time.Sleep(), which is fragile:
	// - If tasks take longer than expected, they get killed
	// - If tasks finish early, we waste time waiting
	// - We have no way to know when goroutines actually finish
	//
	// This fragility is WHY we need the coordination patterns that follow:
	// - Channels: communicate completion
	// - WaitGroup: wait for multiple goroutines
	// - Context: propagate cancellation
	//
	// For this demo, we know the longest task is 400ms, so we wait 500ms.
	time.Sleep(500 * time.Millisecond)

	concDuration := time.Since(concStart)
	fmt.Printf("\nConcurrent total: %v\n", concDuration.Round(time.Millisecond))

	// =========================================================================
	// RESULTS
	// =========================================================================

	speedup := float64(seqDuration) / float64(concDuration)

	fmt.Println("\n=== RESULTS ===")
	fmt.Printf("Sequential: %v (tasks ran one at a time)\n", seqDuration.Round(time.Millisecond))
	fmt.Printf("Concurrent: %v (tasks ran simultaneously)\n", concDuration.Round(time.Millisecond))
	fmt.Printf("Speedup: %.1fx faster with goroutines!\n", speedup)
	fmt.Printf("Time saved: %v\n", (seqDuration - concDuration).Round(time.Millisecond))

	fmt.Println("\n=== KEY INSIGHT ===")
	fmt.Println("Spawning goroutines is trivial: just add `go` before a function call.")
	fmt.Println("But notice we used time.Sleep() to wait - that's fragile!")
	fmt.Println("How do we KNOW when goroutines finish? That's what Channels and")
	fmt.Println("WaitGroup solve - we will look at those patterns soon.")
	fmt.Println()
}
