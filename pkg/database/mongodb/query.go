package mongodb

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap/zapcore"
)

// Query represents mongodb query.
type Query map[string]interface{}

// queries represents mongodb queries.
type queries []interface{}

func (qry Query) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for key, value := range qry {
		switch v := value.(type) {
		case string:
			enc.AddString(key, v)
		case int:
			enc.AddInt(key, v)
		case float64:
			enc.AddFloat64(key, v)
		case []interface{}:
			if err := enc.AddArray(key, queries(v)); err != nil {
				return err
			}
		case bson.M:
			if err := enc.AddObject(key, Query(v)); err != nil {
				return err
			}
		case bson.D:
			if err := enc.AddObject(key, Query(v.Map())); err != nil {
				return err
			}
		default:
			return errors.New("unknown type")
		}
	}
	return nil
}

func (qrs queries) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, value := range qrs {
		switch val := value.(type) {
		case bson.M:
			if err := enc.AppendObject(Query(val)); err != nil {
				return err
			}
		case bson.D:
			if err := enc.AppendObject(Query(val.Map())); err != nil {
				return err
			}
		case string:
			enc.AppendString(val)
		case int:
			enc.AppendInt(val)
		case float64:
			enc.AppendFloat64(val)
		default:
			return errors.New("unknown type")
		}
	}
	return nil
}
