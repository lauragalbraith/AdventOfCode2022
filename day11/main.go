/*
Package main solves Day 11 of Advent of Code 2022
main.go: Laura Galbraith
What is the level of monkey business after 20 rounds of stuff-slinging simian shenanigans?
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"container/list" // https://pkg.go.dev/container/list
	"fmt"
	"regexp"
	"strconv"
	"strings"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil" // GetLinesFromFile
)

type Operation func(old uint64) (new uint64)

type Monkey struct {
	id             int
	reduced_items  *list.List
	inspect_op     Operation
	divisible_test uint64
	true_dest      int
	false_dest     int

	inspections_performed uint64
}

func (m *Monkey) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%d: [", m.id)

	for e := m.reduced_items.Front(); e != nil; e = e.Next() {
		fmt.Fprintf(&sb, "%+v, ", e.Value)
	}
	sb.WriteString("]")

	return sb.String()
}

func (m *Monkey) DeepCopy() *Monkey {
	if m == nil {
		return nil
	}

	new_m := new(Monkey)

	new_m.id = m.id
	new_m.reduced_items = list.New()
	for e := m.reduced_items.Front(); e != nil; e = e.Next() {
		original_e := e.Value.(uint64)
		new_m.reduced_items.PushBack(original_e)
	}
	new_m.inspect_op = m.inspect_op
	new_m.divisible_test = m.divisible_test
	new_m.true_dest = m.true_dest
	new_m.false_dest = m.false_dest
	new_m.inspections_performed = m.inspections_performed

	return new_m
}

func ParseMonkeysFromInput(lines []string) ([]*Monkey, error) {
	// Initialize monkeys
	monkeys := make([]*Monkey, (len(lines)+1)/7)

	operation_line_re := regexp.MustCompile(`^  Operation: new = old (.) (.+)$`)
	test_line_re := regexp.MustCompile(`^  Test: divisible by (\d+)$`)
	action_line_re := regexp.MustCompile(`^    If (true|false): throw to monkey (\d)$`)

	for i := 0; i < len(lines); i++ {
		// Create monkey from block of lines
		m := new(Monkey)
		m.inspections_performed = 0

		// get id
		var err error
		m.id, err = strconv.Atoi(lines[i][len(lines[i])-2 : len(lines[i])-1])
		if err != nil {
			fmt.Printf("Could not convert '%v' to int\n", lines[i][len(lines[i])-2:len(lines[i])-1])
			return nil, err
		}

		i++

		// get starting items
		items_str := lines[i][18:]
		items := strings.Split(items_str, ", ")
		m.reduced_items = list.New()
		for _, item := range items {
			num, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				fmt.Printf("Could not get item number from '%s': %+v\n", item, err)
				return nil, err
			}

			m.reduced_items.PushBack(uint64(num))
		}

		i++

		// get operation
		op_details := operation_line_re.FindAllStringSubmatch(lines[i], -1)
		if len(op_details[0]) < 3 {
			fmt.Printf("Unable to extract operation details from line '%s': %+v\n", lines[i], op_details)
			return nil, err
		}

		if strings.Compare(op_details[0][1], "*") == 0 {
			if strings.Compare(op_details[0][2], "old") == 0 {
				m.inspect_op = func(old uint64) (new uint64) { return old * old }
			} else {
				literal, err := strconv.ParseInt(op_details[0][2], 10, 64)
				if err != nil {
					fmt.Printf("Unable to extract operation literal from line '%s': %+v\n", lines[i], err)
					return nil, err
				}

				m.inspect_op = func(old uint64) (new uint64) { return old * uint64(literal) }
			}
		} else {
			// addition
			literal, err := strconv.ParseInt(op_details[0][2], 10, 64)
			if err != nil {
				fmt.Printf("Unable to extract operation literal from line '%s': %+v\n", lines[i], err)
				return nil, err
			}

			m.inspect_op = func(old uint64) (new uint64) { return old + uint64(literal) }
		}

		i++

		// get test condition
		test_details := test_line_re.FindAllStringSubmatch(lines[i], -1)
		if len(test_details[0]) < 2 {
			fmt.Printf("Unable to extract test details from line '%s': %+v\n", lines[i], test_details)
			return nil, err
		}

		div_64, err := strconv.ParseInt(test_details[0][1], 10, 64)
		if err != nil {
			fmt.Printf("Unable to extract division literal from line '%s': %+v\n", lines[i], err)
			return nil, err
		}

		m.divisible_test = uint64(div_64)

		i++

		// get true case
		true_details := action_line_re.FindAllStringSubmatch(lines[i], -1)
		if len(true_details[0]) < 3 {
			fmt.Printf("Unable to extract true details from line '%s': %+v\n", lines[i], true_details)
			return nil, err
		}

		m.true_dest, err = strconv.Atoi(true_details[0][2])
		if err != nil {
			fmt.Printf("Unable to extract destination monkey id from line '%s': %+v\n", lines[i], err)
			return nil, err
		}

		i++

		// get false case
		false_details := action_line_re.FindAllStringSubmatch(lines[i], -1)
		if len(false_details[0]) < 3 {
			fmt.Printf("Unable to extract false details from line '%s': %+v\n", lines[i], false_details)
			return nil, err
		}

		m.false_dest, err = strconv.Atoi(false_details[0][2])
		if err != nil {
			fmt.Printf("Unable to extract destination monkey id from line '%s': %+v\n", lines[i], err)
			return nil, err
		}

		i++

		// save it to the array
		monkeys[m.id] = m
	}

	return monkeys, nil
}

func MonkeyBusiness(monkeys []*Monkey, rounds uint, undamaged_bonus bool) uint64 {
	// Compute product of all divisors
	var divisor_product uint64
	divisor_product = 1
	for _, m := range monkeys {
		divisor_product *= m.divisible_test
	}

	// Emulate monkeys throwing for X rounds
	var round uint
	for round = 1; round <= rounds; round++ {
		for _, m := range monkeys {
			num_items := m.reduced_items.Len()
			for i := 0; i < num_items; i++ {
				// monkey inspects item
				item := m.reduced_items.Front()
				worry_level := item.Value.(uint64)
				m.reduced_items.Remove(item)
				m.inspections_performed++

				// worry level goes through operation
				worry_level = m.inspect_op(worry_level)

				// monkey gets bored with item; worry level may be divided by 3
				if undamaged_bonus {
					worry_level /= 3
				} else {
					// reduce to remainder after product
					worry_level %= divisor_product
				}

				// test worry level against monkey's condition
				dest_monkey := m.false_dest
				if worry_level%m.divisible_test == 0 {
					dest_monkey = m.true_dest
				}

				// throw item
				monkeys[dest_monkey].reduced_items.PushBack(worry_level)
			}
		}
	}

	// Determine two most active monkeys
	var most, second_most uint64
	most = 0
	second_most = 0

	for _, m := range monkeys {
		if m.inspections_performed > most {
			second_most = most
			most = m.inspections_performed
		} else if m.inspections_performed > second_most {
			second_most = m.inspections_performed
		}
	}

	// Calculate the level of monkey business
	return most * second_most
}

func main() {
	// Get input from file
	lines, err := fileutil.GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	initial_monkeys, err := ParseMonkeysFromInput(lines)
	if err != nil {
		panic(err)
	}

	// Part 1
	// worry is divided by 3 each inspection, 20 rounds
	monkeys := make([]*Monkey, len(initial_monkeys))
	for i, m := range initial_monkeys {
		monkeys[i] = m.DeepCopy()
	}

	monkey_business_level := MonkeyBusiness(monkeys, 20, true)
	fmt.Printf("\nPart 1 answer: %v\n\n", monkey_business_level)

	// Part 2
	// Starting again from the initial state in your puzzle input, what is the level of monkey business after 10000 rounds?
	monkey_business_level = MonkeyBusiness(initial_monkeys, 10000, false)
	fmt.Printf("\nPart 2 answer: %v\n", monkey_business_level)
}
