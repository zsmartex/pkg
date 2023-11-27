package filters

import (
	"github.com/zsmartex/pkg/v2/mpa"
	"go.mongodb.org/mongo-driver/bson"
)

func WithFieldEqual(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, value
	}
}

func WithFieldNotEqual(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$ne": value}
	}
}

func WithFieldGreaterThan(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$gt": value}
	}
}

func WithFieldLessThan(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$lt": value}
	}
}

func WithFieldGreaterThanOrEqualTo(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$gte": value}
	}
}

func WithFieldLesThanOrEqualTo(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$lte": value}
	}
}

func WithFieldIn(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$in": value}
	}
}

func WithFieldNotIn(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$nin": value}
	}
}

func WithFieldLike(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$regex": value}
	}
}

func WithFieldNotLike(key string, value interface{}) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$not": bson.M{"$regex": value}}
	}
}

func WithFieldIsNull(key string) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$exists": false}
	}
}

func WithFieldNotNull(key string) mpa.Filter {
	return func() (k string, v interface{}) {
		return key, bson.M{"$exists": true}
	}
}
