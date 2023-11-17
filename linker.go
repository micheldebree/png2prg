package png2prg

import (
	"fmt"
	"io"
)

type Word uint16

func NewWord(bHi, bLo byte) Word {
	return Word(uint16(bHi)<<8 + uint16(bLo))
}
func (w Word) String() string {
	return fmt.Sprintf("0x%04x", uint16(w))
}
func (w Word) Low() byte {
	return byte(w & 0xff)
}
func (w Word) High() byte {
	return byte(w >> 8)
}
func (w Word) Bytes() []byte {
	return []byte{w.Low(), w.High()}
}

func BytesToWord(bLo, bHi byte) Word {
	return Word(uint16(bHi)<<8 + uint16(bLo))
}

const MaxMemory = 0xffff

type Linker struct {
	cursor  Word
	payload [MaxMemory + 1]byte
	block   [MaxMemory + 1]bool
	used    [MaxMemory + 1]bool
}

func NewLinker(start Word) *Linker {
	return &Linker{cursor: start}
}

func (l *Linker) Used() bool {
	return l.used[l.cursor] || l.block[l.cursor]
}

func (l *Linker) Block(start, end Word) {
	for i := start; i < end; i++ {
		l.block[i] = true
	}
}

func (l *Linker) Cursor() Word {
	return l.cursor
}

func (l *Linker) SetCursor(v Word) {
	l.cursor = v
}

func (l *Linker) SetByte(addr Word, v byte) {
	l.payload[addr] = v
	l.used[addr] = true
}

// CursorWrite sets the cursor and writes b to payload.
func (l *Linker) CursorWrite(cursor Word, b []byte) (n int, err error) {
	l.cursor = cursor
	return l.Write(b)
}

// Write writes b to payload at cursor address and increases the cursor with amount of bytes written.
func (l *Linker) Write(b []byte) (n int, err error) {
	if int(l.cursor)+len(b) > MaxMemory {
		return 0, fmt.Errorf("linker: out of memory error, cursor %s, length %d", l.cursor, len(b))
	}
	for i := 0; i < len(b); i++ {
		if l.Used() {
			return n, fmt.Errorf("linker: memory overlap error, cursor %s, length %d", l.cursor, len(b))
		}
		l.payload[l.cursor] = b[i]
		l.used[l.cursor] = true
		l.cursor++
		n++
	}
	return n, nil
}

func (l *Linker) WritePrg(prg []byte) (n int, err error) {
	if len(prg) < 3 {
		return 0, fmt.Errorf("prg too short to write. length: %d", len(prg))
	}
	l.cursor = BytesToWord(prg[0], prg[1])
	return l.Write(prg[2:])
}

func (l *Linker) StartAddress() Word {
	for i := Word(0); i <= MaxMemory; i++ {
		if l.used[i] {
			return i
		}
	}
	return MaxMemory
}

func (l *Linker) EndAddress() Word {
	for i := Word(MaxMemory); i >= 0; i-- {
		if l.used[i] {
			return i + 1
		}
	}
	return 0
}

// Bytes returns a slice of the used payload.
func (l *Linker) Bytes() []byte {
	start := l.StartAddress()
	end := l.EndAddress()
	if end < start {
		return []byte{}
	}
	return l.payload[start:end]
}

// WriteTo writes the 2 byte startaddress and all used memory linked in one .prg to w.
func (l *Linker) WriteTo(w io.Writer) (n int64, err error) {
	start := l.StartAddress()
	end := l.EndAddress()
	if start >= end {
		return n, fmt.Errorf("linker: Write failed %s >= %s: %w", start, end, err)
	}
	m, err := w.Write(start.Bytes())
	n += int64(m)
	if err != nil {
		return n, fmt.Errorf("linker: Write failed start address %s: %w", start, err)
	}
	m, err = w.Write(l.payload[start:end])
	if err != nil {
		return n, fmt.Errorf("linker: Write failed %s - %s: %w", start, end, err)
	}
	return n, nil
}
