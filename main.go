package main

import (
	"context"
	"os"

	"github.com/illuminarean-labs/yaml-diff-reporter/comparer"
	"github.com/illuminarean-labs/yaml-diff-reporter/domain"
	"github.com/illuminarean-labs/yaml-diff-reporter/parser"
	"github.com/illuminarean-labs/yaml-diff-reporter/reporter"

	"github.com/urfave/cli/v3"
)

func main() {
	var (
		lhsPath    string
		rhsPath    string
		outputPath string
		modes      []string

		lhsAlias string
		rhsAlias string

		ignoredKeys []string
		outputType  string
		format      string
		language    string
	)

	cmd := &cli.Command{
		Name:  "yaml-diff-reporter",
		Usage: "Compare two yaml files and generate a report",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "lhs-path",
				Usage:       "Path to the left-hand-side yaml file",
				Required:    true,
				Destination: &lhsPath,
				Aliases:     []string{"l"},
			},
			&cli.StringFlag{
				Name:        "rhs-path",
				Usage:       "Path to the right-hand-side yaml file",
				Required:    true,
				Destination: &rhsPath,
				Aliases:     []string{"r"},
			},
			&cli.StringFlag{
				Name:        "output-path",
				Usage:       "Path to the output file",
				Required:    false,
				Destination: &outputPath,
				Aliases:     []string{"o"},
			},
			&cli.StringSliceFlag{
				Name:        "modes",
				Usage:       "Compare modes (type, key, index, value)",
				Aliases:     []string{"M"},
				Required:    false,
				Value:       []string{"type", "key", "index", "value"},
				Destination: &modes,
			},
			&cli.StringSliceFlag{
				Name:        "ignored-keys",
				Usage:       "Ignored keys",
				Aliases:     []string{"I"},
				Required:    false,
				Value:       []string{},
				Destination: &ignoredKeys,
			},
			&cli.StringFlag{
				Name:        "output-type",
				Usage:       "Output type (stdout, file)",
				Aliases:     []string{"ot"},
				Required:    false,
				Value:       "stdout",
				Destination: &outputType,
			},
			&cli.StringFlag{
				Name:        "format",
				Usage:       "Report format (json, markdown, plain)",
				Aliases:     []string{"f"},
				Required:    false,
				Value:       "json",
				Destination: &format,
			},
			&cli.StringFlag{
				Name:        "language",
				Usage:       "Report language (en: english, ko: korean)",
				Aliases:     []string{"lang"},
				Required:    false,
				Value:       "en",
				Destination: &language,
			},
			&cli.StringFlag{
				Name:        "lhs-alias",
				Usage:       "Alias for the left-hand-side yaml",
				Aliases:     []string{"la"},
				Required:    false,
				Value:       "lhs",
				Destination: &lhsAlias,
			},
			&cli.StringFlag{
				Name:        "rhs-alias",
				Usage:       "Alias for the right-hand-side yaml",
				Aliases:     []string{"ra"},
				Required:    false,
				Value:       "rhs",
				Destination: &rhsAlias,
			},
		},

		Action: func(ctx context.Context, command *cli.Command) error {
			p := parser.New(parser.Config{
				LHSPath: lhsPath,
				RHSPath: rhsPath,
			})
			yamls, err := p.Parse()
			if err != nil {
				return err
			}

			c := comparer.New(comparer.Config{
				IgnoredKeys: ignoredKeys,
				Modes:       domain.NewCompareModes(modes),
			})

			c.Compare("", yamls.LHS, yamls.RHS)

			r := reporter.New(reporter.Config{
				Format:     domain.ReportFormat(format),
				Language:   domain.ReportLanguage(language),
				LHSAlias:   lhsAlias,
				RHSAlias:   rhsAlias,
				OutputPath: &outputPath,
				OutputType: domain.ReportOutputType(outputType),
			})

			if err = r.Report(*c.Results()); err != nil {
				return err
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
