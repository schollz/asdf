package emitter

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/schollz/logger"
	"go.bug.st/serial"
)

var crowsSetup bool
var crowsChecked time.Time
var crowMutex sync.Mutex
var murder []CrowConnection

type CrowConnection struct {
	Conn  serial.Port
	Batch string
}

type Crow struct {
	Pitch   int
	Env     int
	Attack  float64
	Decay   float64
	Sustain float64
	Release float64
}

func (c Crow) NoteOn(note int, velocity int) {
	if c.Pitch > 0 {
		log.Debugf("[crow%d] output[%d] note_on %d %d", (c.Pitch-1)/4, (c.Pitch-1)%4+1, note, velocity)
		crowCommand(c.Pitch,
			fmt.Sprintf(".volts=%2.3f\n", (float64(note)-48.0)/12.0),
		)
	}
	if c.Env > 0 {
		crowCommand(c.Env, "(true)")
	}
}

func (c Crow) NoteOff(note int) {
	if c.Pitch > 0 {
		log.Debugf("crow%d[%d]: note_ff %d", (c.Pitch-1)/4, c.Pitch, note)
	}
	if c.Env > 0 {
		crowCommand(c.Env, "(false)")
	}
}

func (c Crow) Set(param string, value int) {
	if (param == "attack" || param == "decay" || param == "sustain" || param == "release") && c.Env > 0 {
		log.Debugf("crow%d[%d]: set %s=%d", (c.Env-1)/4, c.Env, param, value)
		switch param {
		case "attack":
			(&c).Attack = float64(value) / 1000.0
		case "decay":
			(&c).Decay = float64(value) / 1000.0
		case "sustain":
			(&c).Sustain = float64(value) / 100.0
		case "release":
			(&c).Release = float64(value) / 1000.0
		}
		cmd := fmt.Sprintf(".action=adsr(%3.3f,%3.3f,%3.3f,%3.3f)", c.Attack, c.Decay, c.Sustain, c.Release)
		crowCommand(c.Env, cmd)
	}
}

func NewCrow(pitch int, envelope int) (c Crow, err error) {
	setupCrows()
	if !crowsSetup {
		err = fmt.Errorf("no crows connected")
	}

	c = Crow{
		Pitch:   (pitch-1)%4 + 1,
		Env:     (envelope-1)%4 + 1,
		Attack:  5,
		Decay:   0.1,
		Sustain: 10.0,
		Release: 5,
	}
	if pitch > 0 && (pitch-1)/4 >= len(murder) {
		err = fmt.Errorf("crow%d not connected", (c.Pitch-1)/4)
	}
	if envelope > 0 && (envelope-1)/4 >= len(murder) {
		err = fmt.Errorf("crow%d not connected", (c.Env-1)/4)
	}
	if envelope > 0 {
		// setup envelope
		cmd := fmt.Sprintf(".action=adsr(%3.3f,%3.3f,%3.3f,%3.3f)", c.Attack, c.Decay, c.Sustain, c.Release)
		crowCommand(c.Env, cmd)
	}
	return
}

func setupCrows() {
	if crowsSetup {
		return
	}
	if time.Since(crowsChecked) < 10*time.Second {
		return
	}
	crowsChecked = time.Now()
	log.Trace("setting up crows")
	crowMutex.Lock()
	defer crowMutex.Unlock()
	ports, err := serial.GetPortsList()
	if err != nil {
		log.Error(err)
		return
	}
	log.Tracef("ports: %+v", ports)
	mode := &serial.Mode{
		BaudRate: 115200,
	}
	for _, port := range ports {
		if strings.Contains(port, "ttyS0") {
			log.Tracef("skipping %+v", port)
			continue
		}
		log.Tracef("connecting to %+v", port)
		var conn serial.Port
		conn, err = serial.Open(port, mode)
		if err != nil {
			log.Tracef("could not open: %+v", err)
			continue
		}
		conn.SetReadTimeout(1 * time.Second)
		_, err = conn.Write([]byte("^^version"))
		if err != nil {
			conn.Close()
			log.Error(err)
			continue
		}
		buf, err := read(&conn)
		if err != nil {
			conn.Close()
			log.Error(err)
			continue
		}
		if bytes.Contains(buf, []byte("v2")) || bytes.Contains(buf, []byte("v3")) || bytes.Contains(buf, []byte("v4")) {
			log.Tracef("[%s] found crow version info", port)
			// setup default
			_, err = conn.Write([]byte("^^clearscript"))
			if err != nil {
				conn.Close()
				log.Error(err)
				continue
			}
			buf, err = read(&conn)
			if err != nil {
				conn.Close()
				log.Error(err)
				continue
			}
			log.Debugf("[%s] crow connected", port)
			murder = append(murder, CrowConnection{
				Conn: conn,
			})
		} else {
			log.Tracef("[%s] closing", port)
			conn.Close()
			log.Tracef("[%s] not a crow", port)
		}

	}
	log.Debugf("found %d crows", len(murder))
	crowsSetup = len(murder) > 0
	return
}

func read(conn *serial.Port) (r []byte, err error) {
	// read while there is something to read
	buf := make([]byte, 100)
	for {
		n, err := (*conn).Read(buf)
		if err != nil || n == 0 {
			break
		}
		r = append(r, buf[:n]...)
	}
	if len(r) == 0 {
		err = fmt.Errorf("unable to read")
	}
	log.Tracef("read %d bytes: %s", len(r), r)
	return
}

func CrowFlush() (err error) {
	if !crowsSetup {
		return
	}
	crowMutex.Lock()
	defer crowMutex.Unlock()
	for crowIndex := range murder {
		if murder[crowIndex].Batch == "" {
			continue
		}
		cmd := murder[crowIndex].Batch + "\n"
		murder[crowIndex].Batch = ""
		log.Tracef("[crow%d] flush: '%s'", crowIndex, strings.TrimSpace(cmd))
		_, err = murder[crowIndex].Conn.Write([]byte(cmd))
		if err != nil {
			log.Error(err)
			return
		}
		murder[crowIndex].Conn.SetReadTimeout(10 * time.Millisecond)
		buf := make([]byte, 512)
		n, err := murder[crowIndex].Conn.Read(buf)
		if err != nil {
			log.Error(err)
		} else if n > 0 {
			log.Tracef("read %d bytes: %s", n, buf[:n])
		}
	}
	return
}

func crowCommand(index int, cmd string) (err error) {
	if !crowsSetup {
		return
	}
	indexCrow := (index - 1) / 4
	luaIndex := (index-1)%4 + 1
	if indexCrow >= len(murder) {
		err = fmt.Errorf("index %d exceeds number of crows (%d)", indexCrow, len(murder))
	}
	cmd = fmt.Sprintf("output[%d]%s", luaIndex, cmd)
	log.Tracef("[crow%d]: %s", indexCrow, cmd)
	cmd = strings.TrimSpace(cmd)
	if len(murder[indexCrow].Batch) > 0 {
		murder[indexCrow].Batch += ";"
	}
	murder[indexCrow].Batch += cmd
	return
}
