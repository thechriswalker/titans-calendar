// initial file generated, but do not regenerate between seasons, edit the file
package main

import "fmt"

var teams = map[int][]string{
	1:  {"Bayside Sharks"},
	2:  {"Bridgwater Blackhawks II"},
	3:  {"Bridgwater Buccaneers"},
	4:  {"Bridport Evolution"},
	5:  {"Devon Peelers"},
	6:  {"Exeter Raptors"},
	7:  {"Exeter Spartans"},
	8:  {"Exmouth Jesters"},
	9:  {"Exmouth Jesters I"},
	10: {"Exmouth Jesters II"},
	11: {"Exmouth Jesters III"},
	12: {"Exmouth Knights"},
	13: {"North Devon"},
	14: {"North Devon II"},
	15: {"North Devon III"},
	16: {"North Devon IV"},
	17: {"Spartan Nomads"},
	18: {"Spartan Nomads II"},
	19: {"Taunton Huish Tigers"},
	20: {"Tiverton Titans"},
	21: {"Tiverton Titans II"},
	22: {"Torbay Tigers"},
	23: {"Torbay Tigers III"},
	24: {"Torbay Tigers Too"},
	25: {"Torbay Tigresses"},
	26: {"UTB Bucks"},
	27: {"UTB Phoenix"},
	28: {"UTB Pirates"},
	29: {"University of Exeter I"},
	30: {"University of Exeter II"},
	31: {"University of Exeter III"},
	32: {"University of Exeter IV"},
}

var invertedIndex = map[string]int{}

func init() {
	for i, list := range teams {
		for _, team := range list {
			invertedIndex[team] = i
		}
	}
}

func getIDFromName(name string) int {
	n, ok := invertedIndex[name]
	if !ok {
		panic(fmt.Errorf("no ID found for %q", name))
	}
	return n
}

