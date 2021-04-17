package midi

type WireType uint8

const (
	MessageUnknown WireType = iota
	MessageNoteOff
	MessageNoteOn
	MessageRealtimeReset
)

func MessageType(x uint64) WireType {
	if kind := x >> 56; kind < 4 {
		return WireType(kind)
	}
	return MessageUnknown
}
