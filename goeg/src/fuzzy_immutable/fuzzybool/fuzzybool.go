package fuzzybool

import "fmt"

type FuzzyBool float32

func New(value interface{}) (FuzzyBool, error) {
    var fuzzy float32
    switch value := value.(type) { // shadow variable
    case float32:
        fuzzy = value
    case float64:
        fuzzy = float32(value)
    case int:
        fuzzy = float32(value)
    case bool:
        fuzzy = 0
        if value {
            fuzzy = 1
        }
    default:
        return FuzzyBool(0), fmt.Errorf("fuzzybool.New(): %v is not a " +
            "number or boolean\n", value)
    }
    if fuzzy < 0 {
        fuzzy = 0
    } else if fuzzy > 1 {
        fuzzy = 1
    }
    return FuzzyBool(fuzzy), nil
}

func (fuzzy FuzzyBool) Copy() FuzzyBool {
    return FuzzyBool(fuzzy)
}

func (fuzzy FuzzyBool) String() string {
    return fmt.Sprintf("%.0f%%", 100*float32(fuzzy))
}

func (fuzzy FuzzyBool) Not() FuzzyBool {
    return FuzzyBool(1 - float32(fuzzy))
}

func (fuzzy FuzzyBool) And(first FuzzyBool,
    rest ...FuzzyBool) FuzzyBool {
    minimum := fuzzy
    rest = append(rest, first)
    for _, other := range rest {
        if minimum > other {
            minimum = other
        }
    }
    return FuzzyBool(minimum)
}

func (fuzzy FuzzyBool) Or(first FuzzyBool,
    rest ...FuzzyBool) FuzzyBool {
    maximum := fuzzy
    rest = append(rest, first)
    for _, other := range rest {
        if maximum < other {
            maximum = other
        }
    }
    return FuzzyBool(maximum)
}

func (fuzzy FuzzyBool) Less(other FuzzyBool) bool {
    return fuzzy < other
}

func (fuzzy FuzzyBool) Equal(other FuzzyBool) bool {
    return fuzzy == other
}

func (fuzzy FuzzyBool) Bool() bool {
    return float32(fuzzy) >= .5
}

func (fuzzy FuzzyBool) Float() float64 {
    return float64(fuzzy)
}
