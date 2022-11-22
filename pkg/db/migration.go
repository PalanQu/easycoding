package db

import (
	"easycoding/pkg/errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
)

type MigrationDirection string

const (
	MigrationDirectionUP   MigrationDirection = "up"
	MigrationDirectionDown MigrationDirection = "down"
)

func MigrationGenerate(migrationDir string, from, to string) (MigrationDirection, int, error) {
	if from == to {
		return "", 0, errors.ErrInvalidRaw("'from' equals to 'to'")
	}
	files, err := ioutil.ReadDir(migrationDir)
	if err != nil {
		return "", 0, errors.ErrInvalid(err)
	}
	migrationSteps := []string{}
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
			return "", 0, errors.ErrInvalidRaw(fmt.Sprintf("invalid match %s", f.Name()))
		}
		timestamp := string(matches[1])
		if !contains(migrationSteps, timestamp) {
			migrationSteps = append(migrationSteps, timestamp)
		}
	}
	sort.Strings(migrationSteps)
	fromIndex := indexOf(migrationSteps, from)
	toIndex := indexOf(migrationSteps, to)
	if fromIndex == -1 || toIndex == -1 {
		return "", 0, errors.ErrInternalRaw(fmt.Sprintf("error from/to name %s, %s", from, to))
	}
	if fromIndex < toIndex {
		return MigrationDirectionUP, toIndex - fromIndex, nil
	}
	return MigrationDirectionDown, fromIndex - toIndex, nil
}

func contains(list []string, s string) bool {
	for _, l := range list {
		if l == s {
			return true
		}
	}
	return false
}

func indexOf(list []string, s string) int {
	for i, l := range list {
		if l == s {
			return i
		}
	}
	return -1
}
