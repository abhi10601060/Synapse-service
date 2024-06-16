package model

import (
	"log"
	"sync"
	"time"
)

type StreamManager struct {
	Streams map[string]*Stream
	sync.RWMutex
}

func (m *StreamManager) CloseStream(pid string) {
	m.Lock()
	delete(m.Streams, pid)
	m.Unlock()
}

func (m *StreamManager) StartStream(streamer *Streamer) {
	pid := streamer.UserId + ":" + time.Now().Format(time.RFC3339Nano)
	log.Println("Room started with Process Id: ", pid)

	stream := Stream{
		Id:       pid,
		Streamer: streamer,
		Viewers:  &map[string]*Viewer{},
	}
	streamer.Stream = &stream
	m.Lock()
	m.Streams[pid] = &stream
	m.Unlock()

	log.Println("Stream Created and added to Manager: ", m.Streams)
}
