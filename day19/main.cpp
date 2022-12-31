// main.cpp: Laura Galbraith
// Description: Solve Day 19 of Advent of Code 2022
// Not Enough Minerals
// Compile and run: rm main.out ; g++ -std=c++17 -static-liblsan -fsanitize=leak ../util/cppfileutil/fileutil.cpp main.cpp -o main.out -Wall -Werror -Wextra -pedantic -Wshadow -Wconversion -fmax-errors=2 && ./main.out

#include "../util/cppfileutil/fileutil.hpp" // ReadLinesFromFile

#include <algorithm>  // move, max
#include <iostream>  // cout, endl, ostream
#include <regex>  // regex, regex_search, smatch
#include <stdexcept>  // invalid_argument
#include <string>  // string, stoi
#include <vector>  // vector

using namespace std;

// order in reverse to prioritize building harder robots first
enum class Resource {kCrackedGeode, kObsidian, kClay, kOre, kResourceLimit};

std::ostream& operator<<(std::ostream& os, const Resource& r) {
  switch(r) {
    case Resource::kCrackedGeode:
      os << "CrackedGeode";
      break;
    case Resource::kObsidian:
      os << "Obsidian";
      break;
    case Resource::kClay:
      os << "Clay";
      break;
    case Resource::kOre:
      os << "Ore";
      break;
    default:
      break;
  }
  return os;
}

class Blueprint {
public:
  const int id_;

  static Blueprint ParseBlueprint(const string& input) {
    // Parse out from regex
    smatch sm;
    regex_search(input, sm, regex(kInputRegex));
    if (sm.size() < 8) {
        throw invalid_argument("input line in unexpected format");
    }

    // Get resource costs of each robot type
    vector<vector<int>> costs(static_cast<size_t>(Resource::kResourceLimit), vector<int>(static_cast<size_t>(Resource::kResourceLimit), 0));
    
    costs[static_cast<size_t>(Resource::kOre)][static_cast<size_t>(Resource::kOre)] = stoi(sm[2]);
    
    costs[static_cast<size_t>(Resource::kClay)][static_cast<size_t>(Resource::kOre)] = stoi(sm[3]);

    costs[static_cast<size_t>(Resource::kObsidian)][static_cast<size_t>(Resource::kOre)] = stoi(sm[4]);
    costs[static_cast<size_t>(Resource::kObsidian)][static_cast<size_t>(Resource::kClay)] = stoi(sm[5]);

    costs[static_cast<size_t>(Resource::kCrackedGeode)][static_cast<size_t>(Resource::kOre)] = stoi(sm[6]);
    costs[static_cast<size_t>(Resource::kCrackedGeode)][static_cast<size_t>(Resource::kObsidian)] = stoi(sm[7]);

    return Blueprint(stoi(sm[1]), move(costs));
  }

  Blueprint(
    const int id_number,
    const vector<vector<int>>&& costs): id_(id_number), robot_costs_(costs)
  {
    // Memoize the max number of robots of each type it would be useful to build
    this->max_useful_robots_ = vector<int>(static_cast<size_t>(Resource::kResourceLimit), 0);

    size_t ore_idx = static_cast<size_t>(Resource::kOre);
    for (auto robot_cost:this->robot_costs_) {
      this->max_useful_robots_[ore_idx] = max(this->max_useful_robots_[ore_idx], robot_cost[ore_idx]);
    }

    this->max_useful_robots_[static_cast<size_t>(Resource::kClay)] = this->robot_costs_[static_cast<size_t>(Resource::kObsidian)][static_cast<size_t>(Resource::kClay)];
    this->max_useful_robots_[static_cast<size_t>(Resource::kObsidian)] = this->robot_costs_[static_cast<size_t>(Resource::kCrackedGeode)][static_cast<size_t>(Resource::kObsidian)];

    this->max_useful_robots_[static_cast<size_t>(Resource::kCrackedGeode)] = 1000;  // (always build geode robots if possible)
  }

  // Determine the maximum quality possible
  int Quality(const int time_limit) const {
    // Start with one ore-collecting robot (free) and a robot factory
    vector<int> robots(static_cast<size_t>(Resource::kResourceLimit), 0);
    robots[static_cast<size_t>(Resource::kOre)] = 1;

    vector<int> collected(static_cast<size_t>(Resource::kResourceLimit), 0);

    int max_geodes_cracked = 0;
    this->MaxGeodesCracked(robots, collected, 0, time_limit, max_geodes_cracked);
    return this->id_ * max_geodes_cracked;
  }

  // returns true if the given resources are enough to build the given robot type
  bool Buildable(const vector<int>& resources, Resource robot_type) const {
    for (Resource r_type = Resource::kCrackedGeode;
      r_type != Resource::kResourceLimit;
      r_type = static_cast<Resource>(static_cast<int>(r_type) + 1))
    {
      size_t resource_type = static_cast<size_t>(r_type);
      if (resources[resource_type] < this->robot_costs_[static_cast<size_t>(robot_type)][resource_type]) {
        return false;
      }
    }

    return true;
  }

  // Do not build another robot if the number of built robots of that type produce enough resources every minute to build anything we want
  bool UnnecessaryToBuild(const vector<int>& built_robots, Resource robot_type) const {
    size_t r_type = static_cast<size_t>(robot_type);
    return built_robots[r_type] >= this->max_useful_robots_[r_type];
  }

  void SubtractRobotCost(vector<int>& resources, Resource robot_type) const {
    for (Resource r_type = Resource::kCrackedGeode;
      r_type != Resource::kResourceLimit;
      r_type = static_cast<Resource>(static_cast<int>(r_type) + 1))
    {
      size_t resource_type = static_cast<size_t>(r_type);
      resources[resource_type] -= this->robot_costs_[static_cast<size_t>(robot_type)][resource_type];
    }
  }

private:
  // vector which effectively maps Resource to a vector of costs (which effectively maps each Resource type to how much of that type must be spent)
  vector<vector<int>> robot_costs_;
  // vector which stores the max number of robots (per resource type) it is useful to build
  // (i.e. the max cost in that resource for any robot) - excepting geode robots
  vector<int> max_useful_robots_;

  inline static const string kInputRegex = "Blueprint (\\d+): Each ore robot costs (\\d+) ore. Each clay robot costs (\\d+) ore. Each obsidian robot costs (\\d+) ore and (\\d+) clay. Each geode robot costs (\\d+) ore and (\\d+) obsidian.";

  void MaxGeodesCracked(
    vector<int>& robots,
    vector<int>& collected,
    int minute,
    const int time_limit,
    int& final_max_geodes) const
  {
    // Check if we have run out of time
    if (minute >= time_limit) {
      final_max_geodes = max(final_max_geodes, collected[static_cast<size_t>(Resource::kCrackedGeode)]);
      return;
    }

    // Prune the path early if we would not be able to reach the current-max geode count, even if we could build a new geode robot for the remaining minutes
    auto natural_number_sum = [](int n){ return n * (n+1) / 2; };
    int current_geode_robots = robots[static_cast<size_t>(Resource::kCrackedGeode)];

    // ex. if we currently have 5 geode robots during minute 22, we can build 2 more, resulting in 5+6 cracked geodes = 6*7/2 - 4*5/2 = nat_num_sum(5+time_limit-minute-1) - nat_num_sum(5-1)
    int best_geodes_possible = collected[static_cast<size_t>(Resource::kCrackedGeode)] + natural_number_sum(current_geode_robots + time_limit - minute - 1) - natural_number_sum(current_geode_robots - 1);

    if (best_geodes_possible <= final_max_geodes) {
      return;
    }

    // Decision point: which type of robot to build this turn
    for (size_t robot_type = static_cast<size_t>(Resource::kCrackedGeode);
      robot_type <= static_cast<size_t>(Resource::kResourceLimit);
      ++robot_type)
    {
      vector<int> collected_copy(collected);
      vector<int> robots_copy(robots);

      Resource robot_to_build = static_cast<Resource>(robot_type);

      // Consider building robot this minute
      if (robot_to_build != Resource::kResourceLimit) {
        if (!this->Buildable(collected, robot_to_build)) {
          continue;
        }

        // Do not build another robot if the number of built robots of that type produce enough resources every minute to build anything we want
        if (this->UnnecessaryToBuild(robots, robot_to_build)) {
          continue;
        }

        // Spend resources to build robot: takes one minute for the robot factory to construct any type of robot, although it consumes the necessary resources available when construction begins
        this->SubtractRobotCost(collected_copy, robot_to_build);
      }

      // Collect resources from all existing robots: each robot can collect 1 of its resource type per minute
      for (size_t resource_i = static_cast<size_t>(Resource::kCrackedGeode);
        resource_i != static_cast<size_t>(Resource::kResourceLimit);
        ++resource_i)
      {
        collected_copy[resource_i] += robots[resource_i];
      }

      if (robot_to_build != Resource::kResourceLimit) {
        // Add built robot to robots list
        // cout << "DEBUG: build robot type " << robot_to_build << " after " << minute << " minutes" << endl;
        robots_copy[robot_type] += 1;
      }

      // Run another minute, which updates total max geodes count
      this->MaxGeodesCracked(robots_copy, collected_copy, minute + 1, time_limit, final_max_geodes);

      // Backtrack this robot build
      // (nothing to do because array copies were created)
    }
  }
};

int main() {
  // Parse input
  vector<string> input = ReadLinesFromFile("example_input.txt");
  vector<Blueprint> blueprints;
  for (auto line:input) {
    blueprints.push_back(Blueprint::ParseBlueprint(line));
  }

  // Part 1
  // maximize the number of opened geodes after 24 minutes, per blueprint
  int kTimeLimit = 24;
  // What do you get if you add up the quality level of all of the blueprints in your list?
  int quality_sum = 0;

  for (auto b:blueprints) {
    int b_quality = b.Quality(kTimeLimit);
    cout << "DEBUG: blueprint " << b.id_ << " has quality " << b_quality << endl;
    quality_sum += b_quality;
  }

  /*
  IDEAS
  - backtracking per blueprint to decide when/which robot is built
    - O(5^24) which is too slow
  - try building more costly robots first
  - create subclasses for robot types, and an enum class for the resource types
  - DP? with what minute each necessary robot is built? with how soon a single geode robot can be built?
  - cut off path if we know we will not meet the current maximum
  - treat as graph problem, with edges between different-robot-count states, and do dijkstra's with lowest-minute distances saved... then treat those as edges and do a dfs where each walk counts up geodes and return the best geode count
  - work backwards from cost of geode robot to cost of the necessary robots by minute 24
  - build the most expensive robot we don't have yet or build geode robot; or else calculate what we are lacking for that and either build the robot that we need a higher projection of resources from or simply wait
    - the next robot can get built at minute X if we just wait; or at minute Y if we build a robot of its first resource type now, or at minute Z if we build a robot of its second resource type now; take the minimum of these minutes as a decision

  - from https://github.com/biggysmith/advent_of_code_2022/blob/master/src/day19/day19.cpp :
    - do not build another robot of a given type if the number of robots of that type already built produce enough resources every minute to build anything we want (i.e. the max cost in that resource for any robot) - excepting geode robots (always build geode robots if possible)
    - prune a path early if we would not be able to reach the current-max geode count even if we could build a new geode robot for the remaining minutes

  - from https://github.com/vss2sn/advent_of_code/blob/master/2022/cpp/day_19a.cpp :
    - do not build a robot if it is not "helpful"; a build action is "helpful" if the number of robots has not reached the point where you can cover the cost every turn (see above point about every minute)
  */
  
  cout << endl << "Part 1 answer: " << quality_sum << endl;

  // Part 2
  kTimeLimit = 32;
  int max_geodes_multiplied = 1;

  for (size_t b_i = 0; b_i < 3 && b_i < blueprints.size(); ++b_i) {
    int b_max_geodes = blueprints[b_i].Quality(kTimeLimit) / blueprints[b_i].id_;
    cout << "DEBUG: blueprint index " << b_i << " can get " << b_max_geodes << " max geodes" << endl;
    max_geodes_multiplied *= b_max_geodes;
  }

  // What do you get if you multiply these numbers together?
  // TODO NEXT save the answers from part 1 (in the blueprint class?) to start max_geodes at for part 2, to enable quicker pruning for part 2
  cout << endl << "Part 2 answer: " << max_geodes_multiplied << endl;

  return 0;
}
