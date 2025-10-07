package mongoDao

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

type fieldName interface {
	String() string
}

func NewMatchBuilder() *MatchBuilder {
	return &MatchBuilder{queries: []bson.E{}}
}

type MatchBuilder struct {
	queries []bson.E
}

func (q *MatchBuilder) AddEqual(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: value})
	}
	return q
}

func (q *MatchBuilder) AddNotEqual(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: bson.M{"$ne": value}})
	}
	return q
}

func (q *MatchBuilder) AddIn(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{key.String(), bson.M{"$in": value}})
	}
	return q
}

func (q *MatchBuilder) AddNotIn(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{key.String(), bson.M{"$nin": value}})
	}
	return q
}

func (q *MatchBuilder) AddGreaterThan(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: bson.M{"$gt": value}})
	}
	return q
}

func (q *MatchBuilder) AddGreaterThanEqual(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: bson.M{"$gte": value}})
	}
	return q
}

func (q *MatchBuilder) AddLessThanEqual(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: bson.M{"lte": value}})
	}
	return q
}

func (q *MatchBuilder) AddLessThan(key fieldName, value interface{}) *MatchBuilder {
	val := reflect.ValueOf(value)
	if !val.IsZero() {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: bson.M{"lt": value}})
	}
	return q
}

func (q *MatchBuilder) AddBoolEqual(key fieldName, value *bool) *MatchBuilder {
	if value == nil {
		return q
	}
	if *value {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: value})
	} else {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: bson.M{"$ne": true}})
	}
	return q
}

func (q *MatchBuilder) AddOr(value bson.A) *MatchBuilder {
	if len(value) == 0 {
		return q
	}
	q.queries = append(q.queries, bson.E{"$or", value})

	return q
}

func (q *MatchBuilder) AddQueries(value []bson.E) *MatchBuilder {
	if len(value) == 0 {
		return q
	}
	q.queries = append(q.queries, value...)

	return q
}

func (q *MatchBuilder) AddBetween(key fieldName, start, end interface{}) *MatchBuilder {
	rangeCond := bson.D{}
	if start != nil {
		rangeCond = append(rangeCond, bson.E{Key: "$gte", Value: start})
	}
	if end != nil {
		rangeCond = append(rangeCond, bson.E{Key: "$lte", Value: end})
	}
	if len(rangeCond) > 0 {
		q.queries = append(q.queries, bson.E{Key: key.String(), Value: rangeCond})
	}
	return q
}

func (q *MatchBuilder) Generate() []bson.E {
	return q.queries
}
