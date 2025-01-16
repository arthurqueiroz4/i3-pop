package main

import (
	"fmt"
	"log"
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
	// Put on front stack
	ds.back.push(data)
	onupdate <- fmt.Sprint("pushed on back")
}

func (ds *DoubleStack[T]) PopOnFrontAndPutOnBack() *T {
	// Pop on front stack and Put on back stack
	onupdate <- fmt.Sprint("popped on front")
	return swapPeek[T](ds.front, ds.back)
}

func (ds *DoubleStack[T]) PopOnBackAndPutOnFront() *T {
	// Pop on back stack and Put on back stack
	onupdate <- fmt.Sprint("popped on back")
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
		log.Printf("[Back Stack] Lenght: %d", ds.back.length())
		log.Printf("[Front Stack] Lenght: %d", ds.front.length())
	}
}

type Stack[T any] interface {
	push(data T)
	pop() *T
	length() int
	peek() *T
}

var _ Stack[int] = (*linkedlist[int])(nil)

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
