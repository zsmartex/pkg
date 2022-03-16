package datatypes

import (
	"database/sql"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// NullString wrapper around sql.NullString
type NullBool struct {
	sql.NullBool
}

// IsZero method is called by bson.IsZero in Mongo for type = NullTime
func (x NullBool) IsZero() bool {
	return !x.Valid
}

// MarshalBSONValue method is called by bson.Marshal in Mongo for type = NullString
func (x *NullBool) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !x.Valid {
		return bsontype.Null, nil, nil
	}

	valueType, b, err := bson.MarshalValue(x.Bool)
	return valueType, b, err
}

// UnmarshalBSONValue method is called by bson.Unmarshal in Mongo for type = NullString
func (x *NullBool) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	raw := bson.RawValue{Type: t, Value: data}

	var uBool bool
	if err := raw.Unmarshal(&uBool); err != nil {
		return err
	}

	if raw.Value == nil {
		x.Valid = false
		return nil
	}

	x.Valid = true
	x.Bool = uBool
	return nil
}
