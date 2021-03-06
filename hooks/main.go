package hooks

import (
	"fmt"
	"lts/logging"
	"lts/watcher"
	"time"

	"lts/screen"
)

// Returns a blocking channel that, when the hook is dead, sends.
func RegisterHook(name string, callback func()) chan int {
	done := make(chan int)

	go func() {

		callbackWrapped := func(reason string) {
			screen.Clear()
			screen.MoveTopLeft()
			logging.Successf("%v [%v]\n\n", reason, name)
			callback()
		}

		switch name {
		case "change":
			// File system change.
			w, err := watcher.New("./")

			if err != nil {
				logging.Errorf("Could not register change hook: %v", err)
				break
			}

			done := make(chan int)
			w.OnChange(func(file string) {
				if GetFileEligibility(file) {
					logging.Infof("Change: %v\n", file)
					DebounceCallback(func() {
						callbackWrapped("Changes were detected in " + file)
					})
				}
			})
			<-done
		case "periodic":
			for {
				time.Sleep(30 * time.Second)
				callbackWrapped("Periodic rerun")
			}
		case "change-all":
			// File system change.
			w, err := watcher.NewRecursive("./")

			if err != nil {
				logging.Errorf("Could not register change hook: %v", err)
				break
			}

			done := make(chan int)
			w.OnChange(func(file string) {
				if GetFileEligibility(file) {
					logging.Infof("Change: %v\n", file)
					DebounceCallback(func() {
						callbackWrapped("Changes were detected in " + file)
					})
				}
			})
			<-done
		}
		fmt.Printf("Close %v\n", name)
		done <- 1
		close(done)
	}()

	return done
}

func WaitForAllHooks(hooks []chan int) {
	for _, v := range hooks {
		<-v
	}
}
