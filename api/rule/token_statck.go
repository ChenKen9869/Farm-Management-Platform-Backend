package rule

import (
	"container/list"
	"sync"
)

type Stack struct {
	list *list.List
	lock *sync.RWMutex
}

func NewStack() *Stack {
	list := list.New()
	l := &sync.RWMutex{}
	return &Stack{list, l}
}

func (stack *Stack) Push(value Token) {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	stack.list.PushBack(value)
}

func (stack *Stack) Pop() Token {
	stack.lock.Lock()
	defer stack.lock.Unlock()
	e := stack.list.Back()
	if e != nil {
		stack.list.Remove(e)
		return e.Value
	}
	return nil
}

func (stack *Stack) Peak() Token {
	e := stack.list.Back()
	if e != nil {
		return e.Token
	}

	return nil
}

func (stack *Stack) Len() int {
	return stack.list.Len()
}

func (stack *Stack) Empty() bool {
	return stack.list.Len() == 0
}
