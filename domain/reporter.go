package domain

type ReportFormat string

type ReportLanguage string

type ReportOutputType string

type Report struct {
	Key         string    `json:"key"`
	ErrorCode   ErrorCode `json:"errorCode"`
	Description string    `json:"description"`
}

type ReportResponse struct {
	Reports []Report `json:"reports"`
}
