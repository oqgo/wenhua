package wenhua

import (
	"time"

	"github.com/oqgo/oqgo"
)

type Tick struct {
	subjectKey oqgo.SubjectKey
	lastPrice  float64
	time       time.Time
}

func (t *Tick) LastPrice() oqgo.Price {
	return oqgo.Price(t.lastPrice)
}

func (t *Tick) SubjectKey() oqgo.SubjectKey {
	return t.subjectKey
}

func (t *Tick) Time() time.Time {
	return t.time
}

func (t *Tick) TradingDay() time.Time {
	panic("unimplemented")
}
