/*
Package main solves Day 13 of Advent of Code 2022
main.go: Laura Galbraith
Distress Signal
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"fmt"
	"strconv"
	"strings"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

type Value interface {
	// initialize the value for starting comparison with another value
	StartComparison() // TODO FINALLY remove if we don't need it to reset state for Part 2
	// returns the next value to compare, nil if they have run out
	Next() Value
}

// IntegerValue represents a Value of type Integer
type IntegerValue struct {
	val int
}

func (iv *IntegerValue) StartComparison() {} // do nothing; caller should change type
func (iv *IntegerValue) Next() Value      { return iv }
func (iv *IntegerValue) String() string   { return fmt.Sprint(iv.val) }

// ListValue represents a Value of type List
type ListValue struct {
	vals        []Value
	compare_pos int
}

func (lv *ListValue) StartComparison() { lv.compare_pos = 0 }
func (lv *ListValue) Next() Value {
	if lv.compare_pos >= len(lv.vals) {
		return nil
	}

	val := lv.vals[lv.compare_pos]
	lv.compare_pos++
	return val
}
func (lv *ListValue) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, val := range lv.vals {
		if i == lv.compare_pos {
			sb.WriteString("*")
		}

		fmt.Fprint(&sb, val)
		if i < len(lv.vals)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")

	return sb.String()
}

// takes a packet in form like "[1,[2,[3,[4,[5,6,7]]]],8,9]"
// input should always start with '[' and end with ']'
func ParseListFromPacket(packet string) (*ListValue, error) {
	list := new(ListValue)

	for i := 1; i < len(packet)-1; i++ {
		var curr Value

		// move past commas to next value
		if packet[i] == ',' {
			continue
		}

		if packet[i] == '[' {
			// parse out inner list
			start_bracket := i

			// find closing ']'
			i++
			for count_inner_brackets := 0; !(count_inner_brackets == 0 && packet[i] == ']'); i++ {
				if packet[i] == '[' {
					count_inner_brackets++
				} else if packet[i] == ']' {
					count_inner_brackets--
				}
			}

			end_bracket := i
			// fmt.Printf("DEBUG: inner packet found at %d:%d: *%s*\n", start_bracket, end_bracket, packet[start_bracket:end_bracket+1])

			inner_list, err := ParseListFromPacket(packet[start_bracket : end_bracket+1])
			if err != nil {
				return list, err
			}

			curr = inner_list
		} else {
			// find character that ends int
			int_start := i
			for ; packet[i] >= '0' && packet[i] <= '9'; i++ {
			}
			val, err := strconv.Atoi(packet[int_start:i])
			if err != nil {
				fmt.Printf("Failed to parse int from '%s': %+v\n", packet[int_start:i], err)
				return list, err
			}

			int_val := new(IntegerValue)
			int_val.val = val
			curr = int_val
		}

		list.vals = append(list.vals, curr)
	}

	return list, nil
}

// order result options
const (
	Incorrect = iota
	Correct
	Continue
)

func Compare(original_left, original_right Value) int {
	// Check if one or both of the values are nil, which indicates the end of a list
	if original_left == nil && original_right == nil {
		return Continue
	} else if original_left == nil {
		return Correct
	} else if original_right == nil {
		return Incorrect
	}

	left := original_left
	right := original_right

	// Match up what types original_left, original_right are
	switch lv := original_left.(type) {
	case *ListValue:
		// fmt.Printf("DEBUG: left is list: %v\n", lv)
		switch rv := original_right.(type) {
		case *IntegerValue:
			// fmt.Printf("DEBUG: right is int: %v\n", rv)
			// must convert right to list to compare
			right_list := new(ListValue)
			right_list.vals = []Value{rv}
			right = right_list
		}
	case *IntegerValue:
		// fmt.Printf("DEBUG: left is int: %v\n", lv)
		switch rv := original_right.(type) {
		case *ListValue:
			// fmt.Printf("DEBUG: right is list: %v\n", rv)
			// must convert left to list to compare
			left_list := new(ListValue)
			left_list.vals = []Value{lv}
			left = left_list
		case *IntegerValue:
			// fmt.Printf("DEBUG: right is int: %v\n", rv)
			// just compare the two ints and return!
			if lv.val < rv.val {
				return Correct
			} else if lv.val > rv.val {
				return Incorrect
			} else {
				return Continue
			}
		}
	}

	// Compare each value in the lists
	left.StartComparison()
	l := left.Next()
	right.StartComparison()
	r := right.Next()

	// continue until we have a definitive result, or we reach the end of one of the lists
	result := Compare(l, r)
	for result == Continue && l != nil && r != nil {
		l = left.Next()
		r = right.Next()
		result = Compare(l, r)
	}

	// fmt.Printf("DEBUG: comparing %v to %v resulted in %v\n", l, r, result)

	return result
}

func main() {
	// Get input
	received_packets, err := fileutil.GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	// Parse packets into values
	var values []*ListValue
	for _, packet := range received_packets {
		// skip blank lines
		if len(packet) <= 0 {
			continue
		}

		// fmt.Printf("DEBUG: creating packet %d... ", i)

		list_value, err := ParseListFromPacket(packet)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("DEBUG: created value %d: %+v\n", i, list_value)

		values = append(values, list_value)
	}

	// Compare pairs of packets
	correct_order_sum := 0
	for pair_index := 1; pair_index <= len(values)/2; pair_index++ {
		left_value_index := (pair_index - 1) * 2
		right_value_index := (pair_index-1)*2 + 1

		// fmt.Printf("DEBUG: comparing %+v to %+v\n", values[left_value_index], values[right_value_index])
		if Compare(values[left_value_index], values[right_value_index]) == Correct {
			// fmt.Println("DEBUG: Correct order")
			correct_order_sum += pair_index
		}
	}

	fmt.Printf("\nPart 1 answer: %+v\n", correct_order_sum)
}
