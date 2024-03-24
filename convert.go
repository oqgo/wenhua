package wenhua

import (
	"time"

	wenhuago "github.com/kkqy/wenhua-go"
	"github.com/oqgo/oqgo"
	chnfutures "github.com/oqgo/util/market/chn-futures"
)

func ConvertKlines(subjectKey oqgo.SubjectKey, duration time.Duration, klines []wenhuago.Kline) []oqgo.Kline {
	oKlines := make([]oqgo.Kline, len(klines))
	for i, kline := range klines {
		oKlines[i] = oqgo.Kline{
			SubjectKey:  subjectKey,
			Duration:    time.Minute,
			TradingDate: chnfutures.TradingDayByTime(kline.Time),
			Time:        kline.Time,
			OpenPrice:   float64(kline.Open),
			HighPrice:   float64(kline.High),
			LowPrice:    float64(kline.Low),
			ClosePrice:  float64(kline.Close),
		}
	}
	return oKlines
}
