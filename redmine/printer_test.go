package redmine

import (
	"reflect"
	"testing"
)

func TestTablePrinter_Print(t *testing.T) {
	type args struct {
		issue            Issue
		printDescription bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "without description",
			args: args{
				issue: Issue{
					AssignedTo: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					Author: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					Description: "",
					Project: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					SpentHours: 0,
					Status: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					Subject: "",
					Tracker: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
				},
			},
			want: []string{
				"*ЗАДАЧА #0*: ",
				"`" + `+-----------+-----+
| ПРОЕКТ    |     |
+-----------+-----+
| СТАТУС    |     |
+-----------+-----+
| ТРЕКЕР    |     |
+-----------+-----+
| ПРИОРИТЕТ |     |
+-----------+-----+
| АВТОР     |     |
+-----------+-----+
| НАЗНАЧЕНО |     |
+-----------+-----+
| ЧАСЫ      | 0.0 |
+-----------+-----+
` + "`",
			},
		},
		{
			name: "with description",
			args: args{
				issue: Issue{
					AssignedTo: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					Author: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					Description: "SSSS",
					Project: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					SpentHours: 1,
					Status: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					Subject: "",
					Tracker: struct {
						ID   int    `json:"id"`
						Name string `json:"name"`
					}{},
					ID: 2,
				},
				printDescription: true,
			},
			want: []string{
				"*ЗАДАЧА #2*: ",
				"`" + `+-----------+-----+
| ПРОЕКТ    |     |
+-----------+-----+
| СТАТУС    |     |
+-----------+-----+
| ТРЕКЕР    |     |
+-----------+-----+
| ПРИОРИТЕТ |     |
+-----------+-----+
| АВТОР     |     |
+-----------+-----+
| НАЗНАЧЕНО |     |
+-----------+-----+
| ЧАСЫ      | 1.0 |
+-----------+-----+
` + "`",
				`ОПИСАНИЕ:

SSSS`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TablePrinter{}
			if got := tr.Print(tt.args.issue, tt.args.printDescription); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TablePrinter.Print()\n%v, want\n%v", got, tt.want)
			}
		})
	}
}
