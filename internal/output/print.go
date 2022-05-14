package output

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

func printText(w io.Writer, text string) error {
	_, err := w.Write([]byte(text))

	return err
}

func printJson(w io.Writer, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	_, err = w.Write(b)

	return err
}

func printYaml(w io.Writer, v interface{}) error {
	// there is no case when the marshal method will return an error
	b, _ := yaml.Marshal(v)

	_, err := w.Write(b)

	return err
}

func printTable(w io.Writer, headers []string, data [][]string) error {
	table := tablewriter.NewWriter(w)
	table.SetHeader(headers)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetBorder(false)

	for i := range data {
		if len(data[i]) != len(headers) {
			return errors.New("number of headers is not equal to number of columns")
		}
		table.Append(data[i])
	}

	table.Render()
	return nil
}
