/*
Package main solves Day 16 of Advent of Code 2022
main.go: Laura Galbraith
Proboscidea Volcanium
Compile and run: rm main.out; go clean; FMT_NEEDED=$(gofmt -e -d main.go | wc -l); if [ $FMT_NEEDED = 0 ]; then go build -o main.out main && ./main.out; else gofmt -e -d main.go; fi
Go 1.19 used
*/
package main

import (
	"container/heap"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	fileutil "github.com/lauragalbraith/AdventOfCode2022/util/gofileutil"
)

type Valve struct {
	name      string
	flow_rate int
	tunnels   []string
}

func (v *Valve) String() string {
	return fmt.Sprintf("valve '%s' has flow rate '%d' and tunnel list '%q'", v.name, v.flow_rate, v.tunnels)
}

var (
	input_re = regexp.MustCompile(`^Valve (.+) has flow rate=(\d+); tunnel[s]{0,1} lead[s]{0,1} to valve[s]{0,1} (.+)$`)
)

func CreateValveForGraph(input string, graph map[string]*Valve) error {
	info := input_re.FindAllStringSubmatch(input, -1)
	if len(info) < 1 || len(info[0]) < 4 {
		return fmt.Errorf("unexpected input format: '%s'", input)
	}

	v := new(Valve)

	// file in info from input
	v.name = info[0][1]

	var err error
	v.flow_rate, err = strconv.Atoi(info[0][2])
	if err != nil {
		return err
	}

	v.tunnels = strings.Split(info[0][3], ", ")

	// add to graph
	graph[v.name] = v

	return nil
}

// Constants from the problem constraints
const (
	START_VALVE   = "AA"
	TIME_ALLOTTED = 30
)

// support storing Valves in a priority queue with the heap type
type ValvePQ struct {
	valve_names          []string
	path_len_from_source map[string]int
}

func (pq *ValvePQ) Len() int {
	return len(pq.valve_names)
}

// needed for sort.Interface
func (pq *ValvePQ) Less(i, j int) bool {
	// return true if i has higher priority than j
	return pq.path_len_from_source[pq.valve_names[i]] < pq.path_len_from_source[pq.valve_names[j]]
}

// needed for sort.Interface
func (pq *ValvePQ) Swap(i, j int) {
	pq.valve_names[i], pq.valve_names[j] = pq.valve_names[j], pq.valve_names[i]
}

func (pq *ValvePQ) Push(x any) {
	new_valve_name := x.(string)
	pq.valve_names = append(pq.valve_names, new_valve_name)
}

func (pq *ValvePQ) Pop() any {
	ret := pq.valve_names[len(pq.valve_names)-1]
	pq.valve_names = pq.valve_names[0 : len(pq.valve_names)-1]
	return ret
}

// returns a list of valve-nodes on the shortest path from source to all other valve-nodes
func Dijkstra(source string, graph map[string]*Valve) map[string][]string {
	// Form a priority queue of cells to try next
	var pq ValvePQ

	// Track the minimum distance to source found
	pq.path_len_from_source = make(map[string]int)
	for v_name, _ := range graph {
		pq.path_len_from_source[v_name] = len(graph) + 1 // set to be longer than possible
	}
	pq.path_len_from_source[source] = 0

	heap.Init(&pq)
	heap.Push(&pq, source)

	// Track the predecessors of valves along best paths
	previous := make(map[string]string)
	previous[source] = source

	// Until we can no longer improve, keep trying the best option in the PQ
	for len(pq.valve_names) > 0 {
		curr_valve := heap.Pop(&pq).(string)

		for _, n := range graph[curr_valve].tunnels {
			// check if visiting from current is an improvement
			n_dist := pq.path_len_from_source[curr_valve] + 1
			if pq.path_len_from_source[n] > n_dist {
				// update paths
				pq.path_len_from_source[n] = n_dist
				previous[n] = curr_valve

				// add neighbor to PQ
				// do not bother removing old value in PQ; it should not amount to anything b/c of improvement comparison
				heap.Push(&pq, n)
			}
		}
	}

	// Record the paths from all valves to the source
	min_paths := make(map[string][]string)

	for v_name, _ := range graph {
		// create the list, dest to source
		min_path := make([]string, 1)
		min_path[0] = v_name

		curr_valve := v_name
		for strings.Compare(curr_valve, source) != 0 {
			curr_valve = previous[curr_valve]
			min_path = append(min_path, curr_valve)
		}

		// reverse the list so it reads source to dest
		for i, j := 0, len(min_path)-1; i < j; i, j = i+1, j-1 {
			min_path[i], min_path[j] = min_path[j], min_path[i]
		}

		min_paths[v_name] = min_path
	}

	return min_paths
}

// Performs DFS on the transformed graph of valves with nonzero flow rate (and AA) as nodes and shortest path between them as edges
// Returns the maximum pressure that can be released from this path
func DFSValuableValveDistance(
	valve_tunnel_graph map[string]*Valve,
	valuable_valve_distance_graph map[string]map[string]int,
	visited map[string]bool,
	minutes []int,
	valves []string,
	pressure_released int,
	open_flow_rates []int) int {

	// If this valve adds value to be opened, open it
	for i, valve := range valves {
		if valve_tunnel_graph[valve].flow_rate > 0 {
			pressure_released += open_flow_rates[i]
			minutes[i]++
			open_flow_rates[i] += valve_tunnel_graph[valve].flow_rate
		}
	}

	// note: the opening should not bring us to the end of our time; we check that before visiting

	// If we didn't visit any more neighbors, set the bar for max at just accumulating pressure from existing open rates
	max_released := pressure_released
	for i, flow := range open_flow_rates {
		max_released += (TIME_ALLOTTED - minutes[i] + 1) * flow
	}

	// Try visiting each unvisited neighbor to continue the path
	for n_0, path_len_0 := range valuable_valve_distance_graph[valves[0]] { // human
		// do not visit if we've already visited this valve (i.e. opened it)
		if has_been_visited, in_map := visited[n_0]; in_map && has_been_visited {
			continue
		}

		// do not bother to visit if we would run out of time by going there and opening its valve
		if minutes[0]+path_len_0+1 > TIME_ALLOTTED {
			continue
		}

		// Mark neighbor as visited before we call it
		visited[n_0] = true

		// NOTE: since we are using slices instead of plain ints, they are passed by reference, and we must make copies of the data to pass so we are not updated
		n_minutes := []int{minutes[0] + path_len_0}
		n_valves := []string{n_0}
		n_flow_rates := []int{open_flow_rates[0]}
		n_pressure_released := pressure_released + (path_len_0 * open_flow_rates[0])

		for n_1, path_len_1 := range valuable_valve_distance_graph[valves[len(valves)-1]] { // elephant
			// do not visit if we've already visited this valve (i.e. opened it)
			if has_been_visited, in_map := visited[n_1]; in_map && has_been_visited && len(valves) > 1 {
				continue
			}

			// do not bother to visit if we would run out of time by going there and opening its valve
			if len(valves) > 1 && minutes[1]+path_len_1+1 > TIME_ALLOTTED {
				continue
			}

			// Add on parameters if we have the elephant
			if len(valves) > 1 {
				// create new slices without referencing the old ones so outer loop values are not overwritten by callee
				n_minutes = []int{minutes[0] + path_len_0, minutes[1] + path_len_1}
				n_valves = []string{n_0, n_1}
				n_flow_rates = []int{open_flow_rates[0], open_flow_rates[1]}

				n_pressure_released = pressure_released + (path_len_0 * open_flow_rates[0]) + (path_len_1 * open_flow_rates[1])

				// Mark neighbor as visited before we call it
				visited[n_1] = true
			}

			// Check if the maximum pressure can be released down this path
			max_pressure_with_n := DFSValuableValveDistance(
				valve_tunnel_graph,
				valuable_valve_distance_graph,
				visited,
				n_minutes,
				n_valves,
				n_pressure_released,
				n_flow_rates)

			if max_pressure_with_n > max_released {
				max_released = max_pressure_with_n
			}

			// If we're dealing with a human only, treat this as a temporary block
			if len(valves) < 2 {
				break
			}

			// Backtrack: unmark neighbor as visited
			visited[n_1] = false
		}

		// Since one creature can stay stil at a juncture longer than another creature, see if that would result in a max
		if len(valves) > 1 && len(n_valves) == 1 {
			// elephant is helping but did not get to move
			// since creatures are interchangeable, human not moving but elephant being able to move on the same valves will only be counted once here (and skipped past when the places are flipped)

			// -> set elephant to stay in one place (whose valve cannot be opened so regular code flow can run)
			n_minutes = append(n_minutes, minutes[1])
			n_valves = append(n_valves, START_VALVE)
			n_flow_rates = append(n_flow_rates, open_flow_rates[1])
			// n_pressure_released = pressure_released + (path_len_0 * open_flow_rates[0]) // stays the same from before

			// Check if the maximum pressure can be released by moving human but not elephant
			max_pressure_with_n := DFSValuableValveDistance(
				valve_tunnel_graph,
				valuable_valve_distance_graph,
				visited,
				n_minutes,
				n_valves,
				n_pressure_released,
				n_flow_rates)

			if max_pressure_with_n > max_released {
				max_released = max_pressure_with_n
			}
		}

		// Backtrack: unmark neighbor as visited
		visited[n_0] = false
	}

	return max_released
}

func main() {
	// valve flow units: pressure per minute in open state
	input_lines, err := fileutil.GetLinesFromFile("input.txt")
	if err != nil {
		panic(err)
	}

	// NOTE no negative flow rates in either input
	// NOTE: all flow rates are unique and <30 but they're not all primes, so we couldn't just factor the 30-minute value so far

	// Store the original as a graph of valves as nodes and tunnels as edges
	valve_tunnel_graph := make(map[string]*Valve)
	for _, line := range input_lines {
		err := CreateValveForGraph(line, valve_tunnel_graph)
		if err != nil {
			panic(err)
		}
	}

	// Transform the original graph into a graph of valves with nonzero flow rate (and AA) as nodes and shortest path between them as edges:
	// This is because the majority of valves have a flow rate of 0, so are just a junction point

	// Compute Dijkstra's shortest path between every valve with a nonzero flow rate (and AA)
	paths_between_valuable_valves := make(map[string]map[string][]string)

	// compute all paths before transforming the graph so we can easily know the list of valuable valves
	for _, valve := range valve_tunnel_graph {
		if valve.flow_rate > 0 || strings.Compare(valve.name, START_VALVE) == 0 {
			paths_between_valuable_valves[valve.name] = Dijkstra(valve.name, valve_tunnel_graph)
		}
	}

	// Create the transformed graph
	valuable_valve_distance_graph := make(map[string]map[string]int)
	// save path length as distance for all valves of note (start valve and valves with flow rates > 0)
	for source_v, min_paths := range paths_between_valuable_valves {
		valuable_valve_distance_graph[source_v] = make(map[string]int)
		for dest_v, path := range min_paths {
			if _, dest_v_is_noteable := paths_between_valuable_valves[dest_v]; dest_v_is_noteable {
				valuable_valve_distance_graph[source_v][dest_v] = len(path) - 1 // path includes both source and dest, but we only need to spend one minute per tunnel between two valves
			}
		}
	}

	// Perform a DFS on the transformed graph to find the path resulting in the maximum released pressure
	visited := make(map[string]bool)
	visited[START_VALVE] = true
	max_released_pressure := DFSValuableValveDistance(
		valve_tunnel_graph,
		valuable_valve_distance_graph,
		visited,
		[]int{1},
		[]string{START_VALVE},
		0,
		[]int{0})

	// What is the most pressure you could release in 30 minutes?
	fmt.Printf("\nPart 1 answer: %v\n", max_released_pressure)

	// Reset visited status for Part 2
	visited = make(map[string]bool)
	visited[START_VALVE] = true

	// Run DFS for both a human and elephant, starting later and both at AA
	max_released_pressure = DFSValuableValveDistance(
		valve_tunnel_graph,
		valuable_valve_distance_graph,
		visited,
		[]int{5, 5},
		[]string{START_VALVE, START_VALVE},
		0,
		[]int{0, 0})

	fmt.Printf("\nPart 2 answer: %v\n", max_released_pressure)
}
