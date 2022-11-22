package db

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
)

type MigrationDirection string

const (
	MigrationDirectionUP   MigrationDirection = "up"
	MigrationDirectionDown MigrationDirection = "down"
)

func MigrationGenerate(migrationDir string, from, to int) (MigrationDirection, int, error) {
	if from == to {
		return "", 0, errors.New("'from' equals to 'to'")
	}
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return "", 0, err
	}
	migrationSteps := []int{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		r := regexp.MustCompile(`^([0-9]{14})_changes\.([a-z]{2,4})\.sql$`)
		if !r.Match([]byte(f.Name())) {
			continue
		}
		matches := r.FindSubmatch([]byte(f.Name()))
		if len(matches) != 3 {
			return "", 0, errors.New(fmt.Sprintf("invalid match %s", f.Name()))
		}
		timestampStr := string(matches[1])
		timestamp, err := strconv.Atoi(timestampStr)
		if err != nil {
			return "", 0, err
		}

		if !contains(migrationSteps, timestamp) {
			migrationSteps = append(migrationSteps, timestamp)
		}
	}
	sort.Slice(migrationSteps, func(i, j int) bool { return migrationSteps[i] < migrationSteps[j] })
	fromIndex := indexOf(migrationSteps, from)
	toIndex := indexOf(migrationSteps, to)
	if fromIndex == -1 || toIndex == -1 {
		return "", 0, errors.New(fmt.Sprintf("error from/to name %v, %v", from, to))
	}
	if fromIndex < toIndex {
		return MigrationDirectionUP, toIndex - fromIndex, nil
	}
	return MigrationDirectionDown, fromIndex - toIndex, nil
}

func contains(list []int, i int) bool {
	for _, l := range list {
		if l == i {
			return true
		}
	}
	return false
}

func indexOf(list []int, i int) int {
	for index, l := range list {
		if l == i {
			return index
		}
	}
	return -1
}
