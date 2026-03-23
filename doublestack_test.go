package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopAndPushOperationOnSigleStack(t *testing.T) {
	s := NewStackWithLimit[int](4)
	for _, v := range []int{1, 2, 3, 4} {
		s.push(v)
	}

	for _, v := range []int{4, 3, 2, 1} {
		assert.Equal(t, v, *s.pop())
	}
}

func TestStackWithLimit(t *testing.T) {
	s := NewStackWithLimit[int](2)

	for _, v := range []int{1, 2, 3, 4} {
		s.push(v)
	}

	assert.Equal(t, 2, s.length())
	assert.Equal(t, 4, *s.pop())
	assert.Equal(t, 3, *s.pop())
}

func TestDoubleStack(t *testing.T) {
	ds := NewDoubleStack[int](3)

	ds.PushOnBack(1)
	ds.PushOnBack(2)
	ds.PushOnBack(3)

	assert.Equal(t, 3, *ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 2, *ds.PeekOnBack())
	assert.Equal(t, 3, *ds.PeekOnFront())

	assert.Equal(t, 2, *ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 1, *ds.PeekOnBack())
	assert.Equal(t, 2, *ds.PeekOnFront())

	assert.Equal(t, 1, *ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 1, *ds.PeekOnFront())
	assert.Nil(t, ds.PeekOnBack())

	assert.Nil(t, ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 3, ds.FrontLength())
	assert.Equal(t, 1, *ds.PeekOnFront())
}
func TestDoubleStackWithLimit(t *testing.T) {
	ds := NewDoubleStack[int](2)

	ds.PushOnBack(1)
	ds.PushOnBack(2)
	ds.PushOnBack(3)

	assert.Equal(t, 3, *ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 2, *ds.PopOnBackAndPutOnFront())
	assert.Nil(t, ds.PopOnBackAndPutOnFront())

	assert.Equal(t, 2, *ds.PopOnFront())
	assert.Equal(t, 3, *ds.PopOnFront())
	assert.Nil(t, ds.PopOnFront())
}
