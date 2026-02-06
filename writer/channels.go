package writer

// Channels for data from listeners to the summarizer
// Each channel contains an individual data packet, already decoded from its original format

type ChannelsSet struct {
	NavtexCh      chan string
	PositionCh    chan string // TODO this will be a special struct, not string
	TextMessageCh chan string
}

func NewChannelsSet() *ChannelsSet {
	return &ChannelsSet{
		NavtexCh:      make(chan string, 10),
		PositionCh:    make(chan string, 10),
		TextMessageCh: make(chan string, 10),
	}
}
