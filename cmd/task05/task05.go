package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Rules map[int]map[int]struct{}

func main() {
	rulesList, updates, err := readInput()
	if err != nil {
		return
	}

	rules := make(Rules)
	for _, rule := range rulesList {
		if _, ok := rules[rule[1]]; !ok {
			rules[rule[1]] = make(map[int]struct{})
		}

		rules[rule[1]][rule[0]] = struct{}{}
	}

	println(sumMiddleParts(rules, updates))
	println(sumCorrectedParts(rules, updates))
}

func sumMiddleParts(rules Rules, updates [][]int) int {
	validUpdates := findUpdates(rules, updates, true)
	sum := 0

	for u := 0; u < len(validUpdates); u++ {
		update := validUpdates[u]
		middle := update[len(update)/2]
		sum += middle
	}

	return sum
}

func sumCorrectedParts(rules Rules, updates [][]int) int {
	invalidUpdates := findUpdates(rules, updates, false)
	sum := 0

	for u := 0; u < len(invalidUpdates); u++ {
		update := invalidUpdates[u]
		for i := 0; i < len(update); i++ {
			r, hasRule := rules[update[i]]
			if hasRule {
				for j := i + 1; j < len(update); j++ {
					if _, o := r[update[j]]; o {
						update[j], update[i] = update[i], update[j]
						i = -1
						break
					}
				}
			}
		}

		sum += update[len(update)/2]
	}

	return sum
}

func findUpdates(rules Rules, updates [][]int, findValid bool) [][]int {
	var result [][]int

	for u := 0; u < len(updates); u++ {
		if isValid(rules, updates[u]) {
			if findValid {
				result = append(result, updates[u])
			}
		} else {
			if !findValid {
				result = append(result, updates[u])
			}
		}
	}

	return result
}

func isValid(rules Rules, update []int) bool {
	valid := true

	for i := 0; i < len(update) && valid; i++ {
		r, hasRule := rules[update[i]]
		if hasRule {
			for j := i + 1; j < len(update); j++ {
				if _, o := r[update[j]]; o {
					valid = false
					break
				}
			}
		}
	}

	return valid
}

func readInput() ([][]int, [][]int, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		println("Failed to open input.txt: " + err.Error())
		return nil, nil, err
	}
	defer file.Close()

	var (
		rules        [][]int
		updates      [][]int
		readingRules = true
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if readingRules {
			if strings.TrimSpace(line) == "" {
				readingRules = false
				continue
			}

			ruleParts := strings.Split(line, "|")

			a, err := strconv.Atoi(ruleParts[0])
			if err != nil {
				panic(err)
			}

			b, err := strconv.Atoi(ruleParts[1])
			if err != nil {
				panic(err)
			}

			rules = append(rules, []int{a, b})
		} else {
			updateParts := strings.Split(line, ",")
			var update []int
			for _, updatePart := range updateParts {
				updateValue, err := strconv.Atoi(updatePart)
				if err != nil {
					panic(err)
				}

				update = append(update, updateValue)
			}
			updates = append(updates, update)
		}
	}

	return rules, updates, nil
}
