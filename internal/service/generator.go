package service

import squids "github.com/sqids/sqids-go"

type simpleIdGenerator struct {
	s *squids.Sqids
}

func NewSimpleIdGenerator() *simpleIdGenerator {
	s, _ := squids.New()
	return &simpleIdGenerator{
		s: s,
	}
}

func (g *simpleIdGenerator) Generate(count uint64) (string, error) {
	return g.s.Encode([]uint64{count, count / 10, count / 100})
}
