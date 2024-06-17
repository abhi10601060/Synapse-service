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
		Viewers:  map[string]*Viewer{},
	}
	streamer.Stream = &stream
	m.Lock()
	m.Streams[pid] = &stream
	m.Unlock()

	log.Println("Stream Created and added to Manager: ", m.Streams)
}

func (m *StreamManager) AddViewerToStream(pid  string, viewer *Viewer) bool {

	stream, isExist := m.Streams[pid]
	if  !isExist {
		log.Println("stream does not exist : ", pid)
		return false
	}
	log.Println("stream exists : ", pid)
	viewer.Stream = stream
	stream.Lock()
	stream.Viewers[viewer.UserId] = viewer
	stream.Unlock()

	log.Println("Viewer added to stream...")
	return true
}

func (m *StreamManager) IsStreamExist(pid string) bool{
	_, isExist := m.Streams[pid]
	return isExist
}