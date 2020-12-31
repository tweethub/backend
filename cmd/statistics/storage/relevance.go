package storage

import (
	"go.uber.org/zap/zapcore"
)

// HML represents the high, median and low amount.
type HML struct {
	High   int64   `json:"high" bson:"high" example:"100"`
	Median float64 `json:"median" bson:"median" example:"50.0"`
	Low    int64   `json:"low" bson:"low" example:"0"`
}

// RelevanceCandle represents the high, median and low amount of retweets and favorites,
// and the amount of tweets that make those amounts in a time span.
type RelevanceCandle struct {
	TimeFrame   string `json:"time_frame" bson:"time_frame" example:"1595700000-1595750000"`
	Favorites   *HML   `json:"favorites" bson:"favorites"`
	Retweets    *HML   `json:"retweets" bson:"retweets"`
	TweetsCount int64  `json:"tweets_count" bson:"tweets_count" example:"100"`
}

type RelevanceCandles []*RelevanceCandle

// Relevance represents relevance candles for a time span.
type Relevance struct {
	TimeSpan string           `json:"time_span" bson:"time_span" example:"10d"`
	Candles  RelevanceCandles `json:"candles" bson:"candles"`
}

// Relevances represents relevances in different time spans.
type Relevances []*Relevance

// NewRelevanceCandle returns relevance candle.
func NewRelevanceCandle() *RelevanceCandle {
	return &RelevanceCandle{
		Favorites:   &HML{},
		Retweets:    &HML{},
		TweetsCount: 0,
	}
}

func (rels Relevances) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, rel := range rels {
		err := enc.AppendObject(rel)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rel *Relevance) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("time_span", rel.TimeSpan)
	return enc.AddArray("candles", rel.Candles)
}

func (cdls RelevanceCandles) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, cdl := range cdls {
		err := enc.AppendObject(cdl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cdl RelevanceCandle) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("time-frame", cdl.TimeFrame)
	enc.AddInt64("tweets-count", cdl.TweetsCount)

	if err := enc.AddObject("favorites", cdl.Favorites); err != nil {
		return err
	}
	if err := enc.AddObject("retweets", cdl.Retweets); err != nil {
		return err
	}
	return nil
}

func (hml *HML) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("high", hml.High)
	enc.AddFloat64("median", hml.Median)
	enc.AddInt64("low", hml.Low)
	return nil
}
