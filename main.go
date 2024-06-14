package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const startYear = 2015
const outputFileName = "README.md"
const tableHeader = "| **Day** | **Language(s)** |\n| --- | --- |"

var languageMap = map[string]string{
	"go": "Go",
	"py": "Python",
	"ts": "TypeScript",
}

type Day struct {
	day       int
	languages []string
}

type Year struct {
	year int
	days []Day
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: <command> <path/to/aoc/directory>")
		os.Exit(1)
	}

	path := os.Args[1]

	generateMarkdownFile(path, parseDirectory(path))
}

func parseDirectory(path string) []Year {
	var years []Year

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, entry := range entries {
		value, err := strconv.Atoi(entry.Name())
		if err != nil || value < startYear || value > time.Now().Year() {
			continue
		}
		years = append(
			years,
			Year{
				year: value,
				days: getDays(filepath.Join(path, entry.Name())),
			},
		)
	}

	return years
}

func getDays(path string) []Day {
	var days []Day
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, entry := range entries {
		value, err := strconv.Atoi(entry.Name())
		if err != nil || value < 1 || value > 25 {
			continue
		}
		days = append(days, Day{
			day:       value,
			languages: detectLanguages(filepath.Join(path, entry.Name())),
		})
	}

	sort.Slice(days, func(a, b int) bool {
		return days[a].day < days[b].day
	})

	return days
}

func detectLanguages(path string) []string {
	languageSet := make(map[string]bool)

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, entry := range entries {
		ext := strings.TrimLeft(filepath.Ext(entry.Name()), ".")
		if language, ok := languageMap[ext]; ok {
			languageSet[language] = true
		}
	}

	var languages []string
	for language := range languageSet {
		languages = append(languages, language)
	}

	sort.Strings(languages)

	return languages
}

func generateMarkdownFile(path string, years []Year) {
	file, err := os.Create(filepath.Join(path, outputFileName))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	file.WriteString("# Advent of Code\n\n")
	file.WriteString("My solutions to the [Advent of Code](https://adventofcode.com/) challenges.\n\n")
	file.WriteString("## Years attempted\n\n")

	for _, year := range years {
		file.WriteString(fmt.Sprintf("### %d\n\n", year.year))
		file.WriteString("<table>\n<tr>\n")

		for i := 0; i < len(year.days); i += 1 {
			if i%5 == 0 && i != 0 {
				file.WriteString("</td>\n")
			}
			if i%5 == 0 {
				file.WriteString("<td>\n\n")
				file.WriteString(tableHeader + "\n")
			}
			file.WriteString(fmt.Sprintf("| [%d](./%d/%d/) | %s |\n", year.days[i].day, year.year, year.days[i].day, strings.Join(year.days[i].languages, ", ")))
			if i == len(year.days)-1 {
				file.WriteString("</td>\n")
			}
		}

		file.WriteString("</tr>\n</table>\n\n")
	}

	fmt.Println("Generated", filepath.Join(path, outputFileName))
}
