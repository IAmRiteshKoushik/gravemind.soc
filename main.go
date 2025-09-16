package main

import (
	"fmt"
	"time"

	"github.com/IAmRiteshKoushik/gravemind/cmd"
	wf "github.com/IAmRiteshKoushik/gravemind/workflows"
)

// For goroutines which might panic, there is a requirement to restart them
// as they are critical to the application
func RecoverRoutine(name string, routine func()) {
	go func() {
		for {
			// Panic recovery block is wrapping the goroutine
			// In-case there is an exit, this function is triggered
			// and the panic is recoverred and logged.
			defer func() {
				if err := recover(); err != nil {
					cmd.Log.Fatal(fmt.Sprintf(name+" panicked. %w", err))
					// Adding a delay to prevent a tight loop
					time.Sleep(5 * time.Second)
				}
			}()
			routine()
			// If the goroutine exists normally, then break out of the loop
			break
		}
	}()
}

func main() {
	RecoverRoutine("bounty-stream", wf.ReadBountyStream)
	RecoverRoutine("achivement-stream", wf.ReadAchivementStream)
	RecoverRoutine("solution-stream", wf.ReadSolutionStream)
}
