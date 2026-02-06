package writer

// Channels for data from listeners to the summarizer
// Each channel contains an individual data packet, already decoded from its original format

type ChannelsSet struct {
	NavtexCh      chan string
	PositionCh    chan PositionMessage
	TextMessageCh chan string
}

func NewChannelsSet() *ChannelsSet {
	return &ChannelsSet{
		NavtexCh:      make(chan string, 10),
		PositionCh:    make(chan PositionMessage, 10),
		TextMessageCh: make(chan string, 10),
	}
}

// Format for individual data packets

type PositionMessage struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}
