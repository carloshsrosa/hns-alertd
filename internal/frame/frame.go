package frame

import "time"

type SubsystemID uint8

const (
	SubsystemPower    SubsystemID = 1
	SubsystemThermal  SubsystemID = 2
	SubsystemComm     SubsystemID = 3
	SubsystemAttitude SubsystemID = 4
	SubsystemPayload  SubsystemID = 5
)

type Frame struct {
	SpacecraftID uint16
	Subsystem    SubsystemID
	Code         uint16
	Value        float64
	Timestamp    time.Time
}
