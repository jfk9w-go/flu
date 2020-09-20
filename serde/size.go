package serde

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

var SizeRegexp = regexp.MustCompile(`^(\d+)\s?(\w+)?$`)

var SizeUnits = map[string]int64{
	"b":  1,
	"Kb": 1 << 10,
	"Mb": 1 << 20,
	"Gb": 1 << 30,
	"Tb": 1 << 40,
}

type Size struct {
	Bytes int64
}

func (s Size) String() string {
	var whole int64
	var finalUnit string
	for unit, divisor := range SizeUnits {
		quotient := s.Bytes / divisor
		if quotient > 0 {
			if whole == 0 || quotient < whole {
				whole = quotient
				finalUnit = unit
			}
		}
	}

	return fmt.Sprintf("%d%s", whole, finalUnit)
}

func (s *Size) FromString(str string) error {
	groups := SizeRegexp.FindStringSubmatch(str)
	if len(groups) < 2 {
		return errors.Errorf("no size match: %s", str)
	}

	amount, err := strconv.ParseInt(groups[1], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "parse size: %s", groups[1])
	}

	unit := "b"
	if groups[2] != "" {
		unit = groups[2]
		if len(unit) == 2 {
			unit = strings.Title(strings.ToLower(unit))
		}
	}

	if multiplier, ok := SizeUnits[unit]; ok {
		s.Bytes = amount * multiplier
		return nil
	}

	return errors.Errorf("unknown unit: %s", unit)
}

func (s *Size) UnmarshalYAML(node *yaml.Node) error {
	return s.FromString(node.Value)
}

func (s *Size) UnmarshalJSON(bytes []byte) error {
	var str string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return errors.Wrap(err, "unmarshal string")
	}

	return s.FromString(str)
}
