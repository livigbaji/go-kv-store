package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var GlobalStore = make(map[string]string)

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
			GlobalStore[key] = value
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

func Get(key string, T *TransactionStack) {
	ActiveTransaction := T.Peek()
	if ActiveTransaction == nil {
		if val, ok := GlobalStore[key]; ok {
			fmt.Printf("%s\n", val)
		} else {
			fmt.Printf("%s not set\n", key)
		}
	} else {
		if val, ok := ActiveTransaction.store[key]; ok {
			fmt.Printf("%s\n", val)
		} else {
			fmt.Printf("%s not set\n", key)
		}
	}
}

func Set(key string, value string, T *TransactionStack) {
	ActiveTransaction := T.Peek()
	if ActiveTransaction == nil {
		GlobalStore[key] = value
	} else {
		ActiveTransaction.store[key] = value
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	items := &TransactionStack{}
	for {
		fmt.Printf("> ")
		text, _ := reader.ReadString('\n')
		// split the text into operation strings
		operation := strings.Fields(text)
		switch operation[0] {
		case "BEGIN":
			items.PushTransaction()
		case "ROLLBACK":
			items.RollbackTransaction()
		case "COMMIT":
			items.Commit()
			items.PopTransaction()
		case "END":
			items.PopTransaction()
		case "SET":
			Set(operation[1], operation[2], items)
		case "GET":
			Get(operation[1], items)
		//case "DELETE":
		//	Delete(operation[1], items)
		//case "COUNT":
		//	Count(operation[1], items)
		case "STOP":
			os.Exit(0)
		default:
			fmt.Printf("ERROR: Unrecognised Operation %s\n", operation[0])
		}
	}
}
