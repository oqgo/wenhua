package wenhua_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oqgo/oqgo"
	"github.com/oqgo/wenhua"
)

func TestWenhua(t *testing.T) {
	engine, _ := oqgo.NewEngine(context.Background(), ".", nil)
	wenhua := wenhua.NewWenhua()
	engine.RegisterModule(wenhua)
	engine.Init(context.Background())
	klines, _ := wenhua.MinuteKlinesByDatetime(oqgo.SubjectKey("SHFE.ni2405"), time.Date(2024, 3, 21, 9, 0, 0, 0, time.Local), time.Date(2024, 3, 22, 9, 0, 0, 0, time.Local))
	for _, kline := range klines {
		fmt.Println(kline)
	}
	fmt.Println(len(klines))
	klines, _ = wenhua.MinuteKlinesByTradingDay(oqgo.SubjectKey("SHFE.ni2405"), time.Date(2024, 3, 21, 9, 0, 0, 0, time.Local))
	for _, kline := range klines {
		fmt.Println(kline)
	}
	fmt.Println(len(klines))
	klines, _ = wenhua.MinuteKlinesUntilAligned(oqgo.SubjectKey("SHFE.ni2405"), 20, time.Date(2024, 3, 22, 14, 28, 0, 0, time.Local))
	for _, kline := range klines {
		fmt.Println(kline)
	}
	fmt.Println(len(klines))
}
