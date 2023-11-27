package filters

import (
	"github.com/zsmartex/pkg/v2/mpa"
	"go.mongodb.org/mongo-driver/bson"
)

func WithFieldEqual(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: value,
		}
	}
}

func WithFieldNotEqual(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$ne": value},
		}
	}
}

func WithFieldGreaterThan(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$gt": value},
		}
	}
}

func WithFieldLessThan(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$lt": value},
		}
	}
}

func WithFieldGreaterThanOrEqualTo(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$gte": value},
		}
	}
}

func WithFieldLesThanOrEqualTo(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$lte": value},
		}
	}
}

func WithFieldIn(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$in": value},
		}
	}
}

func WithFieldNotIn(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$nin": bson.A{value}},
		}
	}
}

func WithFieldLike(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$regex": value},
		}
	}
}

func WithFieldNotLike(key string, value interface{}) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$not": bson.M{"$regex": value}},
		}
	}
}

func WithFieldIsNull(key string) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$exists": false},
		}
	}
}

func WithFieldNotNull(key string) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   key,
			Value: bson.M{"$exists": true},
		}
	}
}
