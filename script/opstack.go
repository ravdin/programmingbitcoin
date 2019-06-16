package script

type opStack struct {
	stack    [][]byte
	Length   int
	Capacity int
}

func newOpStack(cmds [][]byte) *opStack {
	capacity := 16
	for capacity < len(cmds) {
		capacity += 16
	}
	stack := make([][]byte, capacity)
	result := opStack{stack: stack, Length: 0, Capacity: capacity}
	for _, item := range cmds {
		result.push(item)
	}
	return &result
}

func (stack *opStack) push(item []byte) {
	if stack.Length == stack.Capacity {
		stack.Capacity += 16
		newStack := make([][]byte, stack.Capacity)
		copy(newStack, stack.stack)
		stack.stack = newStack
	}
	stack.stack[stack.Length] = item
	stack.Length++
}

func (stack *opStack) pop() []byte {
	if stack.Length == 0 {
		panic("Can't pop from empty stack!")
	}
	stack.Length--
	return stack.stack[stack.Length]
}

func (stack *opStack) peek() []byte {
	return stack.stack[stack.Length-1]
}
