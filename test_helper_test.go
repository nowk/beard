package beard

import (
	"testing"
)

type StepFunc func(testing.TB, Context)

type Context map[string]interface{}

func (c Context) Set(k string, v interface{}) {
	c[k] = v
}

func (c Context) Get(k string) interface{} {
	return c[k]
}

type Asser struct {
	testing.TB
}

func (a Asser) Given(fn StepFunc) *Asserter {
	as := &Asserter{
		t:   a.TB,
		ctx: make(Context),
	}

	fn(as.t, as.ctx)

	return as
}

type Asserter struct {
	t   testing.TB
	ctx Context
}

func (a *Asserter) Given(fn StepFunc) *Asserter {
	fn(a.t, a.ctx)

	return a
}

func (a *Asserter) Then(fn StepFunc) *Asserter {
	fn(a.t, a.ctx)

	return a
}

func (a *Asserter) And(fn StepFunc) *Asserter {
	fn(a.t, a.ctx)

	return a
}
