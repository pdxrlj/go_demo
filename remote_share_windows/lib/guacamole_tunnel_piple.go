package lib

import (
	"sync"
)

type SimpleTunnel struct {
	stream     *Stream
	readerLock sync.Mutex
	writerLock sync.Mutex
}

func NewSimpleTunnel(stream *Stream) *SimpleTunnel {
	return &SimpleTunnel{
		stream:     stream,
		readerLock: sync.Mutex{},
		writerLock: sync.Mutex{},
	}
}

func (t *SimpleTunnel) AcquireReader() *Stream {
	t.readerLock.Lock()
	return t.stream
}

func (t *SimpleTunnel) ReleaseReader() {
	t.readerLock.Unlock()
}

func (t *SimpleTunnel) AcquireWriter() *Stream {
	t.writerLock.Lock()
	return t.stream
}

func (t *SimpleTunnel) ReleaseWriter() {
	t.writerLock.Unlock()
}

func (t *SimpleTunnel) Close() error {
	return t.stream.Close()
}
