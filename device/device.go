package device

import (
	"errors"
	"os/exec"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/rtmididrv"

	driver "gitlab.com/gomidi/rtmididrv"

	"fmt"
	"log"
)

type MidiDevicer interface {
	closeSession() error
	Start() error
	Stop()
}

type noteEventFunc func(p *reader.Position, channel, key, vel uint8)
type sysexEventFunc func(p *reader.Position, data []byte)

type notesState map[uint8]uint8

var DevicePort = "Ableton Push 2 Live Port"

type MidiDevice struct {
	In         *midi.In
	Out        *midi.Out
	Driver     *rtmididrv.Driver
	NotesState notesState
	// Writer writes to midi input
	Writer *writer.Writer

	// Writer reads from midi input
	Reader *reader.Reader
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (p *MidiDevice) closeSession() error {
	return p.Driver.Close()
}

// NewMidiDevice - returns a new instance of MidiDevice
func NewMidiDevice() MidiDevicer {
	return &MidiDevice{}
}

func (p *MidiDevice) Start() error {

	if p.In != nil || p.Out != nil {
		panic(errors.New("Device driver already started"))
	}

	drv, err := driver.New()

	must(err)
	p.Driver = drv

	in, err := midi.OpenIn(drv, -1, DevicePort)
	out, err := midi.OpenOut(drv, -1, DevicePort)

	if err != nil {
		panic(err)
	}

	must(in.Open())
	must(out.Open())

	wr := writer.New(out)

	p.Writer = wr

	p.In = &in
	p.Out = &out

	rd := reader.New(
		reader.Device(func(_ reader.Position, name string) {
			log.Println(name)
		}),
		reader.NoteOn(func(_ *reader.Position, channel uint8, key uint8, velocity uint8) {
			// use this logo to know which key
			println(key)
			// change the midi key to match your device expected midi key
			if key == 36 {
				// call a function that uses the cmd package to execute a command or whatever you need
				executeCommand()
			}

		}),
		reader.RTStart(func() {
			log.Println("Start msg")

		}),
		reader.RTStop(func() {
			log.Println("Stop msg")

		}),
		reader.SysEx(func(r *reader.Position, data []byte) {
			fmt.Sprint("Sysex got %s", string(data))
		}),
	)

	p.Reader = rd

	rd.ListenTo(*p.In)
	println("Device listener started")
	return nil
}

func (p *MidiDevice) Stop() {
	p.closeSession()
}

func executeCommand() {
	println("Executing command")
	cmd := exec.Command("open", "-a", "Ableton Live 10 Suite")

	if out, err := cmd.Output(); err != nil {
		log.Fatal(err)
	} else if out != nil {
		log.Println(string(out))
	}
}
