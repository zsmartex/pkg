package datatypes

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// NullString wrapper around sql.NullString
type NullInt64 struct {
	sql.NullInt64
}

// IsZero method is called by bson.IsZero in Mongo for type = NullTime
func (x NullInt64) IsZero() bool {
	return !x.Valid
}

// MarshalBSONValue method is called by bson.Marshal in Mongo for type = NullString
func (x *NullInt64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !x.Valid {
		return bsontype.Null, nil, nil
	}

	valueType, b, err := bson.MarshalValue(x.Int64)
	return valueType, b, err
}

// UnmarshalBSONValue method is called by bson.Unmarshal in Mongo for type = NullString
func (x *NullInt64) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}

	var uInt64 int64
	if err := raw.Unmarshal(&uInt64); err != nil {
		return err
	}

	if raw.Value == nil {
		x.Valid = false
		return nil
	}

	x.Valid = true
	x.Int64 = uInt64
	return nil
}
