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

	// fmt.Printf("DEBUG: created valve %+v\n", v)

	// add to graph
	graph[v.name] = v

	return nil
}

// Constants from the problem constraints
const (
	START_VALVE   = "AA"
	TIME_ALLOTTED = 30
)

var (
	// representing the valves with a non-zero flow rate
	VALVES_WITH_VALUE []string
	// represents the sum of those flow rates
	SUM_FLOW_RATES int
)

func minutes_pressure(opened_valves map[string]bool, graph map[string]*Valve) int {
	// Add up pressure released during a minute
	pressure := 0
	for valve := range opened_valves {
		pressure += graph[valve].flow_rate
	}
	return pressure
}

// returns the sum of the remaining valuable flow rates
func remaining_valuable_flow_rates(opened_valves map[string]bool, graph map[string]*Valve) int {
	remaining_sum := 0
	for _, v := range graph {
		if _, is_open := opened_valves[v.name]; !is_open {
			// for non-valuable flow rates, 0 is just added
			remaining_sum += v.flow_rate
		}
	}

	return remaining_sum
}

// returns true if DFS is now done
func update_best_if_DFS_done(time_passed int, curr_released_pressure int, max_released_pressure *int) bool {
	if time_passed >= TIME_ALLOTTED {
		// fmt.Printf("DEBUG: finished accumulating pressure: %v\n", curr_released_pressure)

		if curr_released_pressure > *max_released_pressure {
			*max_released_pressure = curr_released_pressure
		}
		return true
	}

	return false
}

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

func FollowValveOrder(
	order []string,
	graph map[string]*Valve,
	paths_between_valuable_valves map[string]map[string][]string,
	max_released_pressure *int) {
	// Keep track of how much pressure we've released, and what time it would be
	minute := 1
	open_amount := 0
	released := 0

	// Keep track of where we are, and what valve we want next
	next_valuable_valve := 0
	curr := START_VALVE

	// Following the given order...
Order_Loop:
	for next_valuable_valve < len(order) && minute <= TIME_ALLOTTED {
		// Follow the most efficient path, increasing time as we travel
		path := paths_between_valuable_valves[curr][order[next_valuable_valve]]

		for _, valve := range path[1:] {
			// fmt.Printf("DEBUG: %vm: move to %s from %s\n", minute, valve, curr)

			released += open_amount
			minute++
			curr = valve

			if minute > TIME_ALLOTTED {
				break Order_Loop
			}
		}

		// Open the valuable valve, increasing time for that operation
		released += open_amount
		// fmt.Printf("DEBUG: %vm: open valve %s\n", minute, curr)
		minute++
		open_amount += graph[curr].flow_rate

		// Look toward the next valve goal
		next_valuable_valve++

		// fmt.Printf("DEBUG: at the start of %vm we have released %v\n", minute, released)
		// Check if the path we're headed down will be able to surpass the current max value; if not, exit early
		if *max_released_pressure >= released+(TIME_ALLOTTED-minute+1)*SUM_FLOW_RATES {
			// fmt.Printf("DEBUG: in %vm: returning early from this path b/c our current pressure %v will not surpass current max %v\n", minute, released, *max_released_pressure)
			return
		}
	}

	// Once we've visited all the valuable valves, just chill and count up pressure
	for ; minute <= TIME_ALLOTTED; minute++ {
		released += open_amount
	}

	// Once we are done, check if we beat the released pressure
	// fmt.Printf("DEBUG: finished order %q with %v pressure\n", order, released)
	if released > *max_released_pressure {
		*max_released_pressure = released
	}
}

// TODO If I could back to this, this isn't fast enough: generate the permutations so that they order (i.e. sort VALVES_WITH_VALUE with sort.Slice) the biggest flow rates first or the ones closest to AA first
func Part1DijkstraToPermutedPathMethod(graph map[string]*Valve) int {
	// Compute Dijkstra's shortest path between every valve with a nonzero flow rate (and AA)
	paths_between_valuable_valves := make(map[string]map[string][]string)
	paths_between_valuable_valves[START_VALVE] = Dijkstra(START_VALVE, graph)

	// Track which valued valves have been opened
	// note: all valves start in closed-state
	for _, valve := range graph {
		if valve.flow_rate > 0 {
			VALVES_WITH_VALUE = append(VALVES_WITH_VALUE, valve.name)
			SUM_FLOW_RATES += valve.flow_rate

			paths_between_valuable_valves[valve.name] = Dijkstra(valve.name, graph)
		}
	}

	/*fmt.Printf("DEBUG: there are %v useful valves (flow rate > 0)\n", len(VALVES_WITH_VALUE))

	for source, paths := range paths_between_valuable_valves {
		fmt.Printf("DEBUG: shortest paths to valuable nodes starting at %s:\n", source)
		for dest, path := range paths {
			fmt.Printf(" to %s: ", dest)
			for _, v := range path {
				fmt.Printf("%s, ", v)
			}
			fmt.Println()
		}
		fmt.Println()
	}*/

	max_released_pressure := 1575 // set a lower limit on answer a little bit after debugging full input

	// List all permutations of the valves with nonzero flow rate, and run those paths (starting from AA, which has 0 flow)

	// Try all possible orders of the important valves
	var Heaps_algo func(arr []string, n int) // declare separately to support recursion
	Heaps_algo = func(arr []string, n int) {
		if n == 1 {
			// run the order without storing it, since permutations may be many
			// fmt.Printf("DEBUG: trying order %q\n", arr)
			FollowValveOrder(arr, graph, paths_between_valuable_valves, &max_released_pressure)
		} else {
			for i := 0; i < n; i++ {
				Heaps_algo(arr, n-1)
				if n%2 == 1 {
					// swap this elem with the "end" of the array
					arr[i], arr[n-1] = arr[n-1], arr[i]
				} else {
					// swap the start and "end" of the array
					arr[0], arr[n-1] = arr[n-1], arr[0]
				}
			}
		}
	}

	Heaps_algo(VALVES_WITH_VALUE, len(VALVES_WITH_VALUE))
	/*for _, perm := range valuable_valve_permutations {
		fmt.Printf("DEBUG: permutation is: ")
		for _, v := range perm {
			fmt.Printf("%s,", v)
		}
		fmt.Println()
	}
	fmt.Println()*/

	return max_released_pressure
}

// Performs DFS on the transformed graph of valves with nonzero flow rate (and AA) as nodes and shortest path between them as edges
// Returns the maximum pressure that can be released from this path
func DFSValuableValveDistance(
	valve_tunnel_graph map[string]*Valve,
	valuable_valve_distance_graph map[string]map[string]int,
	visited map[string]bool,
	minute int,
	valve string,
	pressure_released int,
	open_flow_rate int) int {
	// Mark this node as visited
	visited[valve] = true

	// If this valve adds value to be opened, open it
	if valve_tunnel_graph[valve].flow_rate > 0 {
		pressure_released += open_flow_rate
		minute++
		// fmt.Printf("DEBUG: %vm: valve %s is now open\n", minute, valve)
		open_flow_rate += valve_tunnel_graph[valve].flow_rate
	}

	// note: the opening should not bring us to the end of our time; we check that before visiting

	// TODO if this is too slow, we could add back in the check logic for cutting off a path if it's not going to reach the maximum, but all these paths are better paths now, so I won't bother for this first iteration

	// If we didn't visit any more neighbors, set the bar for max at just accumulating pressure from existing open rates
	max_released := pressure_released + (TIME_ALLOTTED-minute+1)*open_flow_rate

	// Try visiting each unvisited neighbor to continue the path
	for n, path_len := range valuable_valve_distance_graph[valve] {
		// do not visit if we've already visited this valve (i.e. opened it)
		if has_been_visited, in_map := visited[n]; in_map && has_been_visited {
			continue
		}

		// do not bother to visit if we would run out of time by going there and opening its valve
		if minute+path_len+1 > TIME_ALLOTTED {
			continue
		}

		// Check if the maximum pressure can be released down this path
		max_pressure_with_n := DFSValuableValveDistance(
			valve_tunnel_graph,
			valuable_valve_distance_graph,
			visited,
			minute+path_len,
			n,
			pressure_released+path_len*open_flow_rate,
			open_flow_rate)

		if max_pressure_with_n > max_released {
			max_released = max_pressure_with_n
		}
	}

	// Backtrack:
	// Unmark this node as visited
	visited[valve] = false

	// fmt.Printf("DEBUG: down %s's path, %v is the maximum pressure released\n", valve, max_released)

	return max_released
}

func main() {
	// valve flow units: pressure per minute in open state
	// tunnels between valve
	input_lines, err := fileutil.GetLinesFromFile("example_input.txt") // TODO NEXT run main input
	if err != nil {
		panic(err)
	}

	// Store the original as a graph of valves as nodes and tunnels as edges
	valve_tunnel_graph := make(map[string]*Valve)
	for _, line := range input_lines {
		err := CreateValveForGraph(line, valve_tunnel_graph)
		if err != nil {
			panic(err)
		}
	}

	// Transform the original graph into a graph of valves with nonzero flow rate (and AA) as nodes and shortest path between them as edges:

	// Compute Dijkstra's shortest path between every valve with a nonzero flow rate (and AA)
	paths_between_valuable_valves := make(map[string]map[string][]string)
	// compute all paths first so we can know the list of valuable valves
	for _, valve := range valve_tunnel_graph {
		if valve.flow_rate > 0 || strings.Compare(valve.name, START_VALVE) == 0 {
			VALVES_WITH_VALUE = append(VALVES_WITH_VALUE, valve.name)
			SUM_FLOW_RATES += valve.flow_rate

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
	max_released_pressure := DFSValuableValveDistance(
		valve_tunnel_graph,
		valuable_valve_distance_graph,
		make(map[string]bool),
		1,
		START_VALVE,
		0,
		0)

	// NOTE no negative flow rates in either input
	// NOTE: all flow rates are unique and <30 so it would be really easy to track which flows we had already taken in an array, but they're not all primes, so we couldn't just factor the 30-minute value so far
	// NOTE: the majority of valves have a flow rate of 0, so are just a junction point, so most of the complexity should come from it being a graph problem

	// What is the most pressure you could release in 30 minutes?
	fmt.Printf("\nPart 1 answer: %v\n", max_released_pressure)
	// TODO FINALLY when both parts are working for both inputs, delete unused functions

	// IDEA graph flow rate problem?? On further reading, I don't believe this is applicable in its current form, but perhaps we could modify the form so it would be; if a first layer was /29 minutes, and then split off to all of AA's neighbors, but the per-minute flow numbers have nothing to do with how many minutes it's on for, so we couldn't cap incoming flow by minute

	// IDEA graph problem
	// ** IDEA track time during traversal, so that value added to a path is dependent on time, then track the max value of all possible visits? probably looks like DFS

	// IDEA dynamic programming? we are bound to start at a particular valve, and we must travel only to connected valves, and we can choose to open valve or not; if we had a top-down recursive function and had a base case of time=0 when we have to return 0 pressure released, then we could save a memo of [time][start_valve] = max_pressure, where [0][all_valves]=0 BUT we can only open each valve once so [opened_valves] would also have to be in the memo

	// IDEA knapsack problem with 0/1 choice, also time-dependent; 3d memo array: m[i, w, time]? could we define i as current valve rather than considering the first i items? but then it still wouldn't be deterministic which of the previously-visited valves we had taken

	// IDEA since all the flow rates are <30, is there a certain point at which we always know that any remaining closed valves aren't worth moving to?

	// IDEA keep track of how far away all closed valves are, and their flow rates, maybe as a single number (the max value we could obtain by turning on that valve next) and greedily take - but the issue is we could be moving further away from the next-next option

	// IDEA preprocess by creating a Dijkstra's dist value for every valve starting at every other valve, so we now know the min number of steps to take to get from any valve to any other valve; then start at time zero/every valve and work backwards (if a path can't reach AA by time 30, the path is impossible), but we still don't have a heuristic/way of knowing for sure which valve is best to take next, so we'd have to try all options, and this approach is just improving time to calculate time decrementing

	// TODO HERE
	fmt.Printf("\nPart 2 answer: %v\n", max_released_pressure)
}
