package storage

import (
	"os"
	"sync"
)

type Physical struct {
	file *os.File
	mu   sync.Mutex
}

func OpenPhysical(path string) (*Physical, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Physical{file: f}, nil
}

func (p *Physical) Append(data []byte) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	stat, err := p.file.Stat()
	if err != nil {
		return 0, err
	}
	offset := stat.Size()

	_, err = p.file.Write(data)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

func (p *Physical) Read(offset int64, length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := p.file.ReadAt(buf, offset)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
