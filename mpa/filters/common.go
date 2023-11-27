package filters

import (
	"time"

	"github.com/zsmartex/pkg/v2/mpa"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WithID(id string) mpa.Filter {
	return func() bson.E {
		_id, _ := primitive.ObjectIDFromHex(id)

		return bson.E{
			Key:   "_id",
			Value: _id,
		}
	}
}

func WithIDs(ids ...string) mpa.Filter {
	return func() bson.E {
		var _ids []primitive.ObjectID

		for _, id := range ids {
			_id, _ := primitive.ObjectIDFromHex(id)
			_ids = append(_ids, _id)
		}

		return bson.E{
			Key:   "_id",
			Value: bson.M{"$in": _ids},
		}
	}
}

func WithCreatedAtBy(created_at time.Time) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   "created_at",
			Value: created_at,
		}
	}
}

func WithUpdatedAtBy(updated_at time.Time) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   "updated_at",
			Value: updated_at,
		}
	}
}

func WithCreatedAtAfter(created_at time.Time) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   "created_at",
			Value: bson.M{"$gt": created_at},
		}
	}
}

func WithCreatedAtBefore(created_at time.Time) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   "created_at",
			Value: bson.M{"$lt": created_at},
		}
	}
}

func WithUpdatedAtAfter(updated_at time.Time) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   "updated_at",
			Value: bson.M{"$gt": updated_at},
		}
	}
}

func WithUpdatedAtBefore(updated_at time.Time) mpa.Filter {
	return func() bson.E {
		return bson.E{
			Key:   "updated_at",
			Value: bson.M{"$lt": updated_at},
		}
	}

}
