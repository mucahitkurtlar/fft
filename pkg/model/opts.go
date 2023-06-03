package model

type CrawlerOpts struct {
	MaxPageCount   uint16
	GoRoutineCount uint8
	GoToTimeout    float64
	NetIdleTimeout float64
}
