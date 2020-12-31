package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap/zapcore"
)

type Aggregation mongo.Pipeline

func (agn Aggregation) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, step := range agn {
		qry := step.Map()
		return enc.AppendObject(Query(qry))
	}
	return nil
}

func (agn Aggregation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return enc.AddArray("pipeline", agn)
}
