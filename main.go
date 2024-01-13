package main

import "fmt"

var GLobalStore = make(map[string]string)

type Transaction struct {
	store map[string]string
	next  *Transaction
}

type TransactionStack struct {
	top  *Transaction
	size int
}

func (ts *TransactionStack) PushTransaction() {
	temp := Transaction{store: make(map[string]string)}
	temp.next = ts.top
	ts.top = &temp
	ts.size++
}

func (ts *TransactionStack) PopTransaction() {
	if ts.top == nil {
		fmt.Printf("Error: No Active Transactions\n")
	} else {
		node := &Transaction{}
		ts.top = ts.top.next
		node.next = nil
		ts.size--
	}
}

func (ts *TransactionStack) Peek() *Transaction {
	return ts.top
}

func (ts *TransactionStack) Commit() {
	ActiveTransaction := ts.Peek()
	if ActiveTransaction != nil {
		for key, value := range ActiveTransaction.store {
			GLobalStore[key] = value
			if ActiveTransaction != nil {
				ActiveTransaction.next.store[key] = value
			}
		}
	} else {
		fmt.Printf("Info: Nothing to commit\n")
	}
	// TODO write to disk
}

func (ts *TransactionStack) RollbackTransaction() {
	if ts.top == nil {
		fmt.Printf("Error: No Active Transaction\n")
	} else {
		for key := range ts.top.store {
			delete(ts.top.store, key)
		}
	}
}

func main() {

}
