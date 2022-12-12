/*
Package main solves Day 11 of Advent of Code 2022
main.go: Laura Galbraith
What is the level of monkey business after 20 rounds of stuff-slinging simian shenanigans?
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"bufio"
	"container/list" // https://pkg.go.dev/container/list
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetLinesFromFile(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", name, err)
		return nil, err
	}
	defer f.Close()

	lines := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	if err := s.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", name, err)
		return nil, err
	}

	return lines, nil
}

type Operation func(old int64) (new int64)

func Add(old, literal int64) (new int64)    { return old + literal }
func Multiply(old, other int64) (new int64) { return old * other }

type Monkey struct {
	id             int
	items          *list.List
	inspect_op     Operation
	divisible_test int64
	true_dest      int
	false_dest     int

	inspections_performed uint
}

func (m *Monkey) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%d: [", m.id)

	for e := m.items.Front(); e != nil; e = e.Next() {
		fmt.Fprintf(&sb, "%+v, ", e.Value)
	}
	sb.WriteString("]")

	return sb.String()
}

func ParseMonkeysFromInput(lines []string) ([]*Monkey, error) {
	// Initialize monkeys
	monkeys := make([]*Monkey, 8) // cheat a little, see input

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
		m.items = list.New()
		for _, item := range items {
			num, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				fmt.Printf("Could not get item number from '%s': %+v\n", item, err)
				return nil, err
			}

			m.items.PushBack(int64(num))
		}

		i++

		// get operation
		// fmt.Printf("DEBUG: line is *%s*\n", lines[i])
		op_details := operation_line_re.FindAllStringSubmatch(lines[i], -1)
		if len(op_details[0]) < 3 {
			fmt.Printf("Unable to extract operation details from line '%s': %+v\n", lines[i], op_details)
			return nil, err
		}

		if strings.Compare(op_details[0][1], "*") == 0 {
			if strings.Compare(op_details[0][2], "old") == 0 {
				m.inspect_op = func(old int64) (new int64) { return Multiply(old, old) }
			} else {
				literal, err := strconv.ParseInt(op_details[0][2], 10, 64)
				if err != nil {
					fmt.Printf("Unable to extract operation literal from line '%s': %+v\n", lines[i], err)
					return nil, err
				}

				m.inspect_op = func(old int64) (new int64) { return Multiply(old, literal) }
			}
		} else {
			// addition
			literal, err := strconv.ParseInt(op_details[0][2], 10, 64)
			if err != nil {
				fmt.Printf("Unable to extract operation literal from line '%s': %+v\n", lines[i], err)
				return nil, err
			}

			m.inspect_op = func(old int64) (new int64) { return Add(old, literal) }
		}

		i++

		// get test condition
		test_details := test_line_re.FindAllStringSubmatch(lines[i], -1)
		if len(test_details[0]) < 2 {
			fmt.Printf("Unable to extract test details from line '%s': %+v\n", lines[i], test_details)
			return nil, err
		}

		m.divisible_test, err = strconv.ParseInt(test_details[0][1], 10, 64)
		if err != nil {
			fmt.Printf("Unable to extract division literal from line '%s': %+v\n", lines[i], err)
			return nil, err
		}

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
		// fmt.Printf("DEBUG: Have monkey %d: %+v\n", m.id, m)
	}

	return monkeys, nil
}

func main() {
	// Get input from file
	lines, err := GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	monkeys, err := ParseMonkeysFromInput(lines)
	if err != nil {
		panic(err)
	}

	// Emulate monkeys throwing for 20 rounds
	for round := 1; round <= 20; round++ {
		for _, m := range monkeys {
			num_items := m.items.Len()
			for i := 0; i < num_items; i++ {
				// monkey inspects item
				item := m.items.Front()
				worry_level := item.Value.(int64)
				m.items.Remove(item)
				m.inspections_performed++

				// worry level goes through operation
				worry_level = m.inspect_op(worry_level)

				// monkey gets bored with item; worry level is divided by 3
				worry_level /= 3

				// test worry level against monkey's condition
				dest_monkey := m.false_dest
				if worry_level%(m.divisible_test) == 0 {
					dest_monkey = m.true_dest
				}

				// throw item
				// fmt.Printf("DEBUG: monkey %d throwing item %v to monkey %d\n", m.id, worry_level, dest_monkey)
				monkeys[dest_monkey].items.PushBack(worry_level)
			}
		}
	}

	// Determine two most active monkeys
	var most, second_most uint
	most = 0
	second_most = 0

	for _, m := range monkeys {
		// fmt.Printf("DEBUG: monkey %d has inspected items %d times\n", m.id, m.inspections_performed)
		if m.inspections_performed > most {
			second_most = most
			most = m.inspections_performed
		} else if m.inspections_performed > second_most {
			second_most = m.inspections_performed
		}
	}

	// Calculate the level of monkey business
	monkey_business_level := most * second_most
	fmt.Printf("Part 1 answer: %v\n", monkey_business_level)
}
