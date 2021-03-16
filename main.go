package main

import (
	"os"
	"os/signal"
	"sync"

	"github.com/txbrown/midi-os-control/device"
)

func waitForCtrlC() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
	}()
	end_waiter.Wait()
}

func main() {
	d := device.NewMidiDevice()

	d.Start()

	waitForCtrlC()

	d.Stop()
}
