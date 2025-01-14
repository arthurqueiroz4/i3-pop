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
	assert.Equal(t, 2, *ds.back.peek())
	assert.Equal(t, 3, *ds.front.peek())

	assert.Equal(t, 2, *ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 1, *ds.back.peek())
	assert.Equal(t, 2, *ds.front.peek())

	assert.Equal(t, 1, *ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 1, *ds.front.peek())
	assert.Nil(t, ds.back.peek())

	assert.Nil(t, ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 3, ds.front.length())
	assert.Equal(t, 3, *ds.front.peek())
}
func TestDoubleStackWithLimit(t *testing.T) {
	ds := NewDoubleStack[int](2)

	ds.PushOnBack(1)
	ds.PushOnBack(2)
	ds.PushOnBack(3)

	assert.Equal(t, 3, *ds.PopOnBackAndPutOnFront())
	assert.Equal(t, 2, *ds.PopOnBackAndPutOnFront())
	assert.Nil(t, ds.PopOnBackAndPutOnFront())

	assert.Equal(t, 2, *ds.front.pop())
	assert.Equal(t, 3, *ds.front.pop())
	assert.Nil(t, ds.front.pop())
}
