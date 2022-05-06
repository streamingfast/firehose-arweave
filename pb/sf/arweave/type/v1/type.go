package pbcodec

import "time"

func (b *Block) Time() time.Time {
	return time.Unix(int64(b.Timestamp), 0)
}
