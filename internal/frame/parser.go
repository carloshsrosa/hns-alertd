package frame

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

const FrameSize = 24

func Parse(b []byte) (Frame, error) {
	if len(b) != FrameSize {
		return Frame{}, fmt.Errorf("frame parse: expected %d bytes, got %d", FrameSize, len(b))
	}
	f := Frame{
		SpacecraftID: binary.LittleEndian.Uint16(b[0:2]),
		Subsystem:    SubsystemID(b[2]),
		Code:         binary.LittleEndian.Uint16(b[4:6]),
		Value:        math.Float64frombits(binary.LittleEndian.Uint64(b[8:16])),
	}
	nanos := int64(binary.LittleEndian.Uint64(b[16:24]))
	f.Timestamp = time.Unix(0, nanos).UTC()
	if f.Subsystem < SubsystemPower || f.Subsystem > SubsystemPayload {
		return Frame{}, fmt.Errorf("frame parse: unknown subsystem %d", f.Subsystem)
	}
	return f, nil
}
