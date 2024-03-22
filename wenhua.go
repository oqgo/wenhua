package wenhua

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	wenhuago "github.com/kkqy/wenhua-go"
	"github.com/oqgo/oqgo"
	chnfutures "github.com/oqgo/util/market/chn-futures"
)

type Wenhua struct {
	ctx    context.Context
	client *wenhuago.Client

	subjectKeyMap sync.Map // 用于将wenhua的symbol转换为oqgo的subjectKey
}

func (w *Wenhua) Init(ctx oqgo.IModuleHandle) error {
	w.ctx = ctx
	c, err := wenhuago.NewClient(context.Background(), "60.190.146.149:8200")
	if err != nil {
		return err
	}
	w.client = c
	return nil
}
func (w *Wenhua) convertSubjectKeyToSymbol(subjectKey oqgo.SubjectKey) (string, error) {
	splited := strings.Split(string(subjectKey), ".")
	if len(splited) != 2 {
		return "", fmt.Errorf("invalid subject key")
	}
	symbol := splited[1]
	w.subjectKeyMap.Store(symbol, subjectKey) // 缓存subjectKey到symbol的映射关系
	reg := regexp.MustCompile(`([A-Za-z]+)(\d+)`)
	res := reg.FindStringSubmatch(symbol)
	productCode := res[1]
	yearMonth := res[2]
	if len(yearMonth) == 3 {
		yearMonth = time.Now().Format("2006")[2:3] + yearMonth
	}
	return productCode + yearMonth, nil
}

func (w *Wenhua) MinuteKlinesByDatetime(subjectKey oqgo.SubjectKey, startTime time.Time, endTime time.Time) ([]oqgo.Kline, error) {
	symbol, err := w.convertSubjectKeyToSymbol(subjectKey)
	if err != nil {
		return nil, err
	}
	oKlines := make([]oqgo.Kline, 0, 200)
	for {
		klines, err := w.client.Klines(symbol, time.Minute, endTime.Add(-time.Minute), 200)
		if err != nil {
			return nil, err
		}
		if len(klines) == 0 {
			return oKlines, nil
		}
		if klines[0].Time.Before(startTime) {
			for i, kline := range klines {
				if !kline.Time.Before(startTime) {
					oKlines = slices.Concat(ConvertKlines(subjectKey, time.Minute, klines[i:]), oKlines)
					return oKlines, nil
				}
			}
		}
		oKlines = slices.Concat(ConvertKlines(subjectKey, time.Minute, klines), oKlines)
		endTime = oKlines[0].Time
	}
}

func (w *Wenhua) MinuteKlinesByTradingDay(subjectKey oqgo.SubjectKey, tradingDay time.Time) ([]oqgo.Kline, error) {
	tradingDay = chnfutures.TradingDayByTime(tradingDay)
	startTime, endTime := chnfutures.TimeRangeByTradingDay(tradingDay) // 获取交易日的开始和结束时间
	return w.MinuteKlinesByDatetime(subjectKey, startTime, endTime)    // 调用新的方法获取分钟线数据
}

func (w *Wenhua) MinuteKlinesUntilAligned(subjectKey oqgo.SubjectKey, count int, endtime time.Time) ([]oqgo.Kline, error) {
	symbol, err := w.convertSubjectKeyToSymbol(subjectKey)
	if err != nil {
		return nil, err
	}
	klines, err := w.client.Klines(symbol, time.Minute, endtime.Add(-time.Minute), count)
	if err != nil {
		return nil, err
	}
	oKlines := ConvertKlines(subjectKey, time.Minute, klines)
	if len(oKlines) > 0 {
		startTime, _ := chnfutures.TimeRangeByTradingDay(chnfutures.TradingDayByTime(oKlines[0].Time))
		endTime := oKlines[0].Time
		oKlines2, err := w.MinuteKlinesByDatetime(subjectKey, startTime, endTime)
		if err != nil {
			return nil, err
		}
		oKlines = slices.Concat(oKlines2, oKlines)
	}
	return oKlines, nil
}

func (w *Wenhua) Name() string {
	return "github.com/oqgo/wenhua"
}

func NewWenhua() *Wenhua {
	return &Wenhua{}
}
