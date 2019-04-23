package script

type OpStack struct {
	stack    [][]byte
	Length   int
	Capacity int
}

func NewOpStack(cmds [][]byte) *OpStack {
	capacity := 16
	for capacity < len(cmds) {
		capacity += 16
	}
	stack := make([][]byte, capacity)
	result := OpStack{stack: stack, Length: 0, Capacity: capacity}
	for _, item := range cmds {
		result.Push(item)
	}
	return &result
}

func (self *OpStack) Push(item []byte) {
	if self.Length == self.Capacity {
		self.Capacity += 16
		newStack := make([][]byte, self.Capacity)
		copy(newStack, self.stack)
		self.stack = newStack
	}
	self.stack[self.Length] = item
	self.Length++
}

func (self *OpStack) Pop() []byte {
	if self.Length == 0 {
		panic("Can't pop from empty stack!")
	}
	self.Length--
	return self.stack[self.Length]
}

func (self *OpStack) Peek() []byte {
	return self.stack[self.Length-1]
}
