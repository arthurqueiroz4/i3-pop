package main

import (
	"fmt"
	"log"
	"strings"
)

type DoubleStack[T any] struct {
	front Stack[T]
	back  Stack[T]
}

var onupdate = make(chan interface{})

func NewDoubleStack[T any](limit int) *DoubleStack[T] {
	ds := &DoubleStack[T]{
		front: NewStackWithLimit[T](limit),
		back:  NewStackWithLimit[T](limit),
	}
	go ds.monitor()
	return ds
}

func (ds *DoubleStack[T]) PushOnBack(data T) {
	ds.back.push(data)
	onupdate <- struct{}{}
}

func (ds *DoubleStack[T]) PopOnBack() *T {
	v := ds.back.pop()
	onupdate <- struct{}{}
	return v
}

func (ds *DoubleStack[T]) PushOnFront(data T) {
	ds.front.push(data)
	onupdate <- struct{}{}
}

func (ds *DoubleStack[T]) PopOnFront() *T {
	v := ds.front.pop()
	onupdate <- struct{}{}
	return v
}

func (ds *DoubleStack[T]) PopOnFrontAndPutOnBack() *T {
	// Pop on front stack and Put on back stack
	onupdate <- "popped on front"
	return swapPeek[T](ds.front, ds.back)
}

func (ds *DoubleStack[T]) PopOnBackAndPutOnFront() *T {
	// Pop on back stack and Put on back stack
	onupdate <- "popped on back"
	return swapPeek[T](ds.back, ds.front)
}

func swapPeek[T any](pop, put Stack[T]) *T {
	if d := pop.pop(); d != nil {
		put.push(*d)
		return d
	}
	return nil
}

func (ds *DoubleStack[T]) monitor() {
	log.Println("Starting stack monitor")
	for {
		<-onupdate
		log.Printf("[Back Stack] Lenght: %d -- [ %s ]", ds.back.length(), formatStack[string](ds.back))
		log.Printf("[Front Stack] Lenght: %d -- [ %s ]", ds.front.length(), formatStack[string](ds.front))
	}
}

func formatStack[T any](stack any) string {
	var sb strings.Builder
	cur := stack.(*linkedlist[T]).head

	for cur != nil {
		sb.WriteString(fmt.Sprintf("%v ", cur.value))
		cur = cur.next
	}

	return sb.String()
}

type Stack[T any] interface {
	push(data T)
	pop() *T
	length() int
	peek() *T
}

var _ Stack[any] = (*linkedlist[any])(nil)

func NewStackWithLimit[T any](limit int) Stack[T] {
	return &linkedlist[T]{
		head:  nil,
		len:   0,
		limit: limit,
	}
}

type linkedlist[T any] struct {
	head  *node[T]
	tail  *node[T]
	len   int
	limit int
}

type node[T any] struct {
	value T
	next  *node[T]
	prev  *node[T]
}

func (ll *linkedlist[T]) push(data T) {
	n := &node[T]{
		value: data,
		next:  ll.head,
		prev:  nil,
	}

	if ll.head != nil {
		ll.head.prev = n
	}
	ll.head = n

	if ll.len == 0 {
		ll.tail = n
	}

	ll.len++

	if ll.limit > 0 && ll.len > ll.limit {
		ll.tail = ll.tail.prev
		if ll.tail != nil {
			ll.tail.next = nil
		}
		ll.len--
	}
}

func (ll *linkedlist[T]) pop() *T {
	if ll.len == 0 || ll.head == nil {
		return nil
	}
	n := ll.head
	ll.head = ll.head.next
	ll.len--
	if ll.len == 0 {
		ll.tail = nil
	}
	return &n.value
}

func (ll *linkedlist[T]) length() int {
	return ll.len
}

func (ll *linkedlist[T]) peek() *T {
	if ll.head == nil {
		return nil
	}
	return &ll.head.value
}
