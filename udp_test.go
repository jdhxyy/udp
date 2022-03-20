package udp

import (
	"github.com/jdhxyy/lagan"
	"testing"
	"time"
)

func TestCase1(t *testing.T) {
	lagan.SetFilterLevel(lagan.LevelDebug)
	_ = Load(0, 2002, 4096)
	Send([]uint8("jdh99"), 0x7f000001, 2003)

	select {
	case <-time.After(time.Second * 5):
	}
}
