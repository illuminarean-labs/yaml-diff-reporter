package parser

import (
	"os"

	"github.com/illuminarean-labs/yaml-diff-reporter/domain"

	"gopkg.in/yaml.v3"
)

type Parser interface {
	Parse() (domain.ParserResult, error)
}

type parser struct {
	config Config
}

func New(config Config) Parser {
	return parser{config: config}
}

type Config struct {
	LHSPath string
	RHSPath string
}

func (p parser) Parse() (domain.ParserResult, error) {
	var (
		lhs map[string]any
		rhs map[string]any
	)

	lhsFile, err := os.ReadFile(p.config.LHSPath)
	if err != nil {
		return domain.ParserResult{}, err
	}

	if err = yaml.Unmarshal(lhsFile, &lhs); err != nil {
		return domain.ParserResult{}, err
	}

	rhsFile, err := os.ReadFile(p.config.RHSPath)
	if err != nil {
		return domain.ParserResult{}, err
	}

	if err = yaml.Unmarshal(rhsFile, &rhs); err != nil {
		return domain.ParserResult{}, err
	}

	return domain.ParserResult{LHS: lhs, RHS: rhs}, nil
}
