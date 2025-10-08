// take the tsv and extract all the team names, a good start.

const readline = require('readline');

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false
});

const teams = new Set();


rl.on('line', (line) => {
    // split by tab.
    const split = line.split("\t")
    // home and away on parts 6 and 8.
    const [home, away] = [split[6]?.trim(), split[8]?.trim()]
    if (home && home !== "Home Team") { // ignore the "header" row
        teams.add(home)
    }
    if (away && away !== "Away Team") {
        teams.add(away)
    }
});

rl.once('close', () => {
    // end of input
    const t = Array.from(teams);
    t.sort();
    process.stdout.write(`// initial file generated, but do not regenerate between seasons, edit the file
package main

import "fmt"

var teams = map[int][]string{
${t.map((team, i) => `\t${i + 1}:${i < 9 ? " " : ""} {${JSON.stringify(team)}},`).join('\n')}
}

var invertedIndex = map[string]int{}

func init() {
\tfor i, list := range teams {
\t\tfor _, team := range list {
\t\t\tinvertedIndex[team] = i
\t\t}
\t}
}

func getIDFromName(name string) int {
\tn, ok := invertedIndex[name]
\tif !ok {
\t\tpanic(fmt.Errorf("no ID found for %q", name))
\t}
\treturn n
}

`)

});