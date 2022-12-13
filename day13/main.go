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
	"sort"
	"strconv"
	"strings"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

type Value interface {
	// initialize the value for starting comparison with another value
	StartComparison()
	// returns the next value to compare, nil if they have run out
	Next() Value
	// returns true if this Value represents a divider packet
	IsDivider() bool
	// utility function so sort can swap elements safely
	DeepCopy() Value
}

// IntegerValue represents a Value of type Integer
type IntegerValue struct {
	val        int
	is_divider bool
}

func (iv *IntegerValue) StartComparison() {} // do nothing; caller should change type
func (iv *IntegerValue) Next() Value      { return iv }
func (iv *IntegerValue) IsDivider() bool  { return iv.is_divider }
func (iv *IntegerValue) DeepCopy() Value {
	new_iv := new(IntegerValue)

	new_iv.val = iv.val
	new_iv.is_divider = iv.is_divider

	return new_iv
}
func (iv *IntegerValue) String() string { return fmt.Sprint(iv.val) }

// ListValue represents a Value of type List
type ListValue struct {
	vals        []Value
	compare_pos int
	is_divider  bool
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
func (lv *ListValue) IsDivider() bool { return lv.is_divider }
func (lv *ListValue) DeepCopy() Value {
	new_lv := new(ListValue)

	for _, val := range lv.vals {
		new_lv.vals = append(new_lv.vals, val.DeepCopy())
	}

	new_lv.compare_pos = lv.compare_pos
	new_lv.is_divider = lv.is_divider

	return new_lv
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
func ParseListFromPacket(packet string, is_divider bool) (*ListValue, error) {
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

			inner_list, err := ParseListFromPacket(packet[start_bracket:end_bracket+1], false) // inner lists of divider are not, themselves, a divider
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

	if is_divider {
		list.is_divider = true
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
		switch rv := original_right.(type) {
		case *IntegerValue:
			// must convert right to list to compare
			right_list := new(ListValue)
			right_list.vals = []Value{rv}
			right = right_list
		}
	case *IntegerValue:
		switch rv := original_right.(type) {
		case *ListValue:
			// must convert left to list to compare
			left_list := new(ListValue)
			left_list.vals = []Value{lv}
			left = left_list
		case *IntegerValue:
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

	return result
}

// Define a type that sort can use
type Packets []Value

func (p Packets) Len() int { return len(p) }
func (p Packets) Swap(i, j int) {
	temp := p[i].DeepCopy()
	p[i] = p[j].DeepCopy()
	p[j] = temp.DeepCopy()
}
func (p Packets) Less(i, j int) bool { return Compare(p[i], p[j]) != Incorrect }

func main() {
	// Get input
	received_packets, err := fileutil.GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	// Parse packets into values
	var values []Value
	for _, packet := range received_packets {
		// skip blank lines
		if len(packet) <= 0 {
			continue
		}

		list_value, err := ParseListFromPacket(packet, false)
		if err != nil {
			panic(err)
		}

		values = append(values, list_value)
	}

	// Compare pairs of packets
	correct_order_sum := 0
	for pair_index := 1; pair_index <= len(values)/2; pair_index++ {
		left_value_index := (pair_index - 1) * 2
		right_value_index := (pair_index-1)*2 + 1

		if Compare(values[left_value_index], values[right_value_index]) == Correct {
			correct_order_sum += pair_index
		}
	}

	fmt.Printf("\nPart 1 answer: %+v\n", correct_order_sum)

	// Part 2
	// add divider packets into list
	divider2, err := ParseListFromPacket("[[2]]", true)
	if err != nil {
		panic(err)
	}
	values = append(values, divider2)

	divider6, err := ParseListFromPacket("[[6]]", true)
	if err != nil {
		panic(err)
	}
	values = append(values, divider6)

	// sort all of the packets into the correct order, including the dividers
	sort.Sort(Packets(values))

	// answer is the indeces of the divider packets multiplied together
	decoder_key := 1
	for i, packet := range values {
		if packet.IsDivider() {
			decoder_key *= i + 1 // 1-indexed number from 0-indexed list
		}
	}

	fmt.Printf("\nPart 2 answer: %v\n", decoder_key)
}
