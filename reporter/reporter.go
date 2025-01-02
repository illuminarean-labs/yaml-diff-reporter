package reporter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/illuminarean-labs/yaml-diff-reporter/domain"

	"github.com/samber/lo"
)

const (
	JSON     domain.ReportFormat = "json"
	Markdown domain.ReportFormat = "markdown"
	Plain    domain.ReportFormat = "plain"
)

const (
	EN domain.ReportLanguage = "en"
	KO domain.ReportLanguage = "ko"
)

const (
	Stdout domain.ReportOutputType = "stdout"
	File   domain.ReportOutputType = "file"
)

type Reporter interface {
	Report(results domain.ErrorResults) error
}

type Config struct {
	Format     domain.ReportFormat
	Language   domain.ReportLanguage
	LHSAlias   string
	RHSAlias   string
	OutputPath *string
	OutputType domain.ReportOutputType
}

type reporter struct {
	config Config
}

func New(config Config) Reporter {
	return reporter{config: config}
}

func (r reporter) Report(results domain.ErrorResults) error {
	var (
		report string
		err    error
	)

	if len(results) == 0 {
		fmt.Println("No differences found")
		return nil
	}

	switch r.config.Format {
	case JSON:
		report, err = r.generateJsonReport(results)
		if err != nil {
			return err
		}
	case Markdown:
		report, err = r.generateMarkdownReport(results)
		if err != nil {
			return err
		}
	case Plain:
		report, err = r.generatePlainTextReport(results)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported report mode")
	}

	switch r.config.OutputType {
	case Stdout:
		r.printReport(report)
	case File:
		if err = r.writeReportFile(report); err != nil {
			return err
		}
	default:
		return errors.New("unsupported output type")
	}

	return nil
}

func (r reporter) generatePlainTextReport(results domain.ErrorResults) (string, error) {
	plainText := ""

	var descriptionMap map[domain.ReportLanguage]string
	for _, result := range results {

		switch result.ErrorCode {
		case domain.ErrorTypeUnmatched:
			descriptionMap = map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("- [%s]키의 타입이 일치하지 않습니다. %s: %s, %s: %s\n", result.Key, r.config.LHSAlias, result.LHS.Type, r.config.RHSAlias, result.RHS.Type),
				EN: fmt.Sprintf("- [%s]Type unmatched. %s: %s, %s: %s\n", result.Key, r.config.LHSAlias, result.LHS.Type, r.config.RHSAlias, result.RHS.Type),
			}
		case domain.ErrorValueUnmatched:
			descriptionMap = map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("- [%s]키의 값이 일치하지 않습니다. %s: (%s)%s, %s: (%s)%s\n", result.Key, r.config.LHSAlias, result.LHS.Type, result.LHS.Value, r.config.RHSAlias, result.RHS.Type, result.RHS.Value),
				EN: fmt.Sprintf("- [%s]Value unmatched. %s: (%s)%s, %s: (%s)%s\n", result.Key, r.config.LHSAlias, result.LHS.Type, result.LHS.Value, r.config.RHSAlias, result.RHS.Type, result.RHS.Value),
			}
		case domain.ErrorKeyNotFound:
			var sideAlias string
			if result.FindNilSide() == "LHS" {
				sideAlias = r.config.LHSAlias
			} else {
				sideAlias = r.config.RHSAlias
			}

			descriptionMap = map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("- %s에서 [%s]키가 존재하지 않습니다.\n", sideAlias, result.Key),
				EN: fmt.Sprintf("- Key not found in %s. key:[%s]\n", sideAlias, result.Key),
			}
		case domain.ErrorIndexNotFound:
			var sideAlias string
			if result.FindNilSide() == "LHS" {
				sideAlias = r.config.LHSAlias
			} else {
				sideAlias = r.config.RHSAlias
			}

			descriptionMap = map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("- %s에서 [%s]인덱스가 존재하지 않습니다.\n", sideAlias, result.Key),
				EN: fmt.Sprintf("- Index not found in %s. [%s]\n", sideAlias, result.Key),
			}
		default:
			return "", errors.New("unsupported error code")
		}

		plainText += descriptionMap[r.config.Language]
	}

	return plainText, nil
}

func (r reporter) generateJsonReport(results domain.ErrorResults) (string, error) {
	reports := make([]domain.Report, 0, len(results))

	for _, result := range results {
		switch result.ErrorCode {
		case domain.ErrorTypeUnmatched:
			DescriptionMap := map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("타입이 일치하지 않습니다. %s: %s, %s: %s",
					r.config.LHSAlias, result.LHS.Type,
					r.config.RHSAlias, result.RHS.Type,
				),
				EN: fmt.Sprintf("Type unmatched. %s: %s, %s: %s",
					r.config.LHSAlias, result.LHS.Type,
					r.config.RHSAlias, result.RHS.Type,
				),
			}

			reports = append(reports, domain.Report{
				Key:         result.Key,
				ErrorCode:   result.ErrorCode,
				Description: DescriptionMap[r.config.Language],
			})
		case domain.ErrorValueUnmatched:
			DescriptionMap := map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("값이 일치하지 않습니다. %s: (%s)%s, %s: (%s)%s",
					r.config.LHSAlias, result.LHS.Type, result.LHS.Value,
					r.config.RHSAlias, result.RHS.Type, result.RHS.Value,
				),
				EN: fmt.Sprintf("Value unmatched. %s: (%s)%s, %s: (%s)%s",
					r.config.LHSAlias, result.LHS.Type, result.LHS.Value,
					r.config.RHSAlias, result.RHS.Type, result.RHS.Value,
				),
			}

			reports = append(reports, domain.Report{
				Key:         result.Key,
				ErrorCode:   result.ErrorCode,
				Description: DescriptionMap[r.config.Language],
			})
		case domain.ErrorKeyNotFound:
			var sideAlias string
			if result.FindNilSide() == "LHS" {
				sideAlias = r.config.LHSAlias
			} else {
				sideAlias = r.config.RHSAlias
			}

			DescriptionMap := map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("키가 존재하지 않습니다. %s", sideAlias),
				EN: fmt.Sprintf("Key not found. %s", sideAlias),
			}

			reports = append(reports, domain.Report{
				Key:         result.Key,
				ErrorCode:   result.ErrorCode,
				Description: DescriptionMap[r.config.Language],
			})

		case domain.ErrorIndexNotFound:
			var sideAlias string
			if result.FindNilSide() == "LHS" {
				sideAlias = r.config.LHSAlias
			} else {
				sideAlias = r.config.RHSAlias
			}

			DescriptionMap := map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("인덱스가 존재하지 않습니다. %s", sideAlias),
				EN: fmt.Sprintf("Index not found. %s", sideAlias),
			}

			reports = append(reports, domain.Report{
				Key:         result.Key,
				ErrorCode:   result.ErrorCode,
				Description: DescriptionMap[r.config.Language],
			})
		default:
			return "", errors.New("unsupported error code")
		}
	}

	reportJson, err := json.Marshal(domain.ReportResponse{Reports: reports})
	if err != nil {
		return "", err
	}

	return string(reportJson), nil
}

func (r reporter) generateMarkdownReport(results domain.ErrorResults) (string, error) {
	report := "## Difference Report\n\n"

	report += fmt.Sprintf("| Key | Error Code | %s | %s | Description |\n",
		r.config.LHSAlias, r.config.RHSAlias,
	)
	report += "| --- | --- | --- | --- | --- |\n"

	var descriptionMap map[domain.ReportLanguage]string
	for _, result := range results {
		switch result.ErrorCode {
		case domain.ErrorTypeUnmatched:
			descriptionMap = map[domain.ReportLanguage]string{
				KO: "타입이 일치하지 않습니다. ",
				EN: "Type unmatched.",
			}
		case domain.ErrorValueUnmatched:
			descriptionMap = map[domain.ReportLanguage]string{
				KO: "값이 일치하지 않습니다.",
				EN: "Value unmatched.",
			}
		case domain.ErrorKeyNotFound:
			descriptionMap = map[domain.ReportLanguage]string{
				KO: "키가 존재하지 않습니다.",
				EN: "Key not found.",
			}

		case domain.ErrorIndexNotFound:
			descriptionMap = map[domain.ReportLanguage]string{
				KO: fmt.Sprintf("인덱스가 존재하지 않습니다."),
				EN: fmt.Sprintf("Index not found."),
			}
		default:
			return "", errors.New("unsupported error code")
		}

		report += fmt.Sprintf("| `%s` | `%s` | `(%s)%s` | `(%s)%s` | %s |\n",
			result.Key, result.ErrorCode, result.LHS.Type, result.LHS.Value, result.RHS.Type, result.RHS.Value, descriptionMap[r.config.Language],
		)
	}

	return report, nil
}

func (r reporter) printReport(report string) {
	fmt.Println(report)
}

func (r reporter) writeReportFile(report string) error {
	if lo.FromPtr(r.config.OutputPath) == "" {
		return errors.New("output path is required")
	}

	// 디렉토리 존재 확인 후 없는 경우 생성
	dirPath := path.Dir(*r.config.OutputPath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, 0755); err != nil {
			return err
		}
	}

	// 파일 쓰기
	if err := os.WriteFile(*r.config.OutputPath, []byte(report), 0644); err != nil {
		return err
	}

	fmt.Printf("Report has been saved to %s", *r.config.OutputPath)

	return nil
}
