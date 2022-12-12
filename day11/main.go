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
	"math/big" // https://pkg.go.dev/math/big
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

type Operation func(old *big.Int) (new *big.Int)

type Monkey struct {
	id             int
	items          *list.List
	inspect_op     Operation
	divisible_test *big.Int
	true_dest      int
	false_dest     int

	inspections_performed uint64
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

func (m *Monkey) DeepCopy() *Monkey {
	if m == nil {
		return nil
	}

	new_m := new(Monkey)

	new_m.id = m.id
	new_m.items = list.New()
	for e := m.items.Front(); e != nil; e = e.Next() {
		var original_e *big.Int
		original_e = e.Value.(*big.Int)
		if !original_e.IsInt64() {
			panic(fmt.Errorf("Item value %v inside monkey %d cannot be copied\n", original_e, m.id))
		}

		new_m.items.PushBack(big.NewInt(original_e.Int64()))
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
	// fmt.Printf("There are %v lines so there are %v monkeys\n", len(lines), len(monkeys))

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

			m.items.PushBack(big.NewInt(num))
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
				m.inspect_op = func(old *big.Int) (new *big.Int) { return old.Mul(old, old) }
			} else {
				literal, err := strconv.ParseInt(op_details[0][2], 10, 64)
				if err != nil {
					fmt.Printf("Unable to extract operation literal from line '%s': %+v\n", lines[i], err)
					return nil, err
				}

				m.inspect_op = func(old *big.Int) (new *big.Int) { return old.Mul(old, big.NewInt(literal)) }
			}
		} else {
			// addition
			literal, err := strconv.ParseInt(op_details[0][2], 10, 64)
			if err != nil {
				fmt.Printf("Unable to extract operation literal from line '%s': %+v\n", lines[i], err)
				return nil, err
			}

			m.inspect_op = func(old *big.Int) (new *big.Int) { return old.Add(old, big.NewInt(literal)) }
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

		m.divisible_test = big.NewInt(div_64)

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

func MonkeyBusiness(monkeys []*Monkey, rounds uint, undamaged_bonus bool) uint64 {
	// for _, m := range monkeys {
	// fmt.Printf("DEBUG: monkey %+v has inspected items %d times\n", m, m.inspections_performed)
	// }

	// Emulate monkeys throwing for X rounds
	var round uint
	for round = 1; round <= rounds; round++ {
		if round%100 == 0 {
			fmt.Printf("Round %d\n", round)
		}
		for _, m := range monkeys {
			num_items := m.items.Len()
			for i := 0; i < num_items; i++ {
				// monkey inspects item
				item := m.items.Front()
				worry_level := item.Value.(*big.Int)
				m.items.Remove(item)
				m.inspections_performed++

				// worry level goes through operation
				worry_level = m.inspect_op(worry_level)

				// monkey gets bored with item; worry level may be divided by 3
				if undamaged_bonus {
					worry_level, _ = worry_level.DivMod(worry_level, big.NewInt(3), big.NewInt(0))
				}

				// test worry level against monkey's condition
				dest_monkey := m.false_dest

				remainder := big.NewInt(0)
				quotient := big.NewInt(0)
				quotient, remainder = quotient.DivMod(worry_level, m.divisible_test, remainder)
				if remainder.IsInt64() && remainder.Int64() == 0 {
					dest_monkey = m.true_dest
				}

				// throw item
				// fmt.Printf("DEBUG: monkey %d throwing item %v to monkey %d\n", m.id, worry_level, dest_monkey)
				monkeys[dest_monkey].items.PushBack(worry_level)
			}
		}

		if round == 1 || round == 20 || round%1000 == 0 {
			fmt.Printf("DEBUG: after round %d...\n", round)
			for _, m := range monkeys {
				fmt.Printf("DEBUG: monkey %+v inspected items %d times\n", m.id, m.inspections_performed)
			}
		}
	}

	// Determine two most active monkeys
	var most, second_most uint64
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
	return most * second_most
}

func main() {
	// Get input from file
	lines, err := GetLinesFromFile("example_input.txt")
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
		// fmt.Printf("DEBUG: Copying monkey %+v\n", m)
		monkeys[i] = m.DeepCopy()
	}

	monkey_business_level := MonkeyBusiness(monkeys, 20, true)
	fmt.Printf("Part 1 answer: %v\n", monkey_business_level)

	// Part 2
	// Starting again from the initial state in your puzzle input, what is the level of monkey business after 10000 rounds?
	monkey_business_level = MonkeyBusiness(initial_monkeys, 10000, false)
	fmt.Printf("Part 2 answer: %v\n", monkey_business_level) // TODO 26823519255 is too high; round 20 numbers are correct, but rund 1000 is where it starts to differ (all monkeys differ)
	// TODO NEXT debug ideas: see if any true/false branches are not taken prior to round 21; see if there's list limits that are being hit; see if there are any variables that are changing during emulation that shouldn't be
	// TODO HERE our int type is too small to hold some of the numbers; numbers are becoming negative
}
