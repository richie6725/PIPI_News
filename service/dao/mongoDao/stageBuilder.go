package mongoDao

import "go.mongodb.org/mongo-driver/bson"

func NewStageBuilder() *StageBuilder {
	return &StageBuilder{
		pipeline: []bson.D{},
	}
}

type StageBuilder struct {
	pipeline []bson.D
}

func (s *StageBuilder) AddMatch(queries []bson.E) *StageBuilder {
	if len(queries) > 0 {
		s.pipeline = append(s.pipeline, bson.D{{"$match", queries}})
	}
	return s
}

func (s *StageBuilder) AddSort(queries []bson.E) *StageBuilder {
	if len(queries) > 0 {
		s.pipeline = append(s.pipeline, bson.D{{"sort", queries}})
	}
	return s
}

func (s *StageBuilder) AddSearch(queries []bson.E) *StageBuilder {
	if len(queries) > 0 {
		s.pipeline = append(s.pipeline, bson.D{{"$search", queries}})
	}
	return s
}

func (s *StageBuilder) AddCount() *StageBuilder {
	s.pipeline = append(s.pipeline, bson.D{{"$count", "count"}})
	return s
}

func (s *StageBuilder) AddGroupBy(queries []bson.E) *StageBuilder {
	if len(queries) > 0 {
		s.pipeline = append(s.pipeline, bson.D{{"$group", queries}})
	}
	return s
}

// AddCustomQueries is utilized to specify a custom query method for cases where other requirements do not apply.
func (s *StageBuilder) AddCustomQueries(queries []bson.D) *StageBuilder {
	if len(queries) > 0 {
		s.pipeline = append(s.pipeline, queries...)
	}
	return s
}

func (s *StageBuilder) Generate() []bson.D {
	return s.pipeline
}
