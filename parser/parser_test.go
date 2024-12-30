package parser

import (
	"testing"

	"github.com/YangTaeyoung/yaml-diff-reporter/domain"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name    string
		args    args
		want    domain.ParserResult
		wantErr bool
	}{
		{
			name: "성공",
			args: args{
				config: Config{LHSPath: "./test.yml", RHSPath: "./test.yml"},
			},
			want: domain.ParserResult{
				LHS: map[string]any{
					"hello": map[string]any{
						"name": map[string]any{
							"firstName": "John",
							"lastName":  "Doe",
						},
						"age":     30,
						"address": "1234 Elm St.",
						"phones": []any{
							map[string]any{
								"type":   "home",
								"number": "123-456-7890",
							},
							map[string]any{
								"type":   "office",
								"number": "098-765-4321",
							},
						},
						"cars": []any{
							"porche",
							"ferrari",
							"lamborghini",
						},
					},
				},
				RHS: map[string]any{
					"hello": map[string]any{
						"name": map[string]any{
							"firstName": "John",
							"lastName":  "Doe",
						},
						"age":     30,
						"address": "1234 Elm St.",
						"phones": []any{
							map[string]any{
								"type":   "home",
								"number": "123-456-7890",
							},
							map[string]any{
								"type":   "office",
								"number": "098-765-4321",
							},
						},
						"cars": []any{
							"porche",
							"ferrari",
							"lamborghini",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			p := New(tt.args.config)
			got, err := p.Parse()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
