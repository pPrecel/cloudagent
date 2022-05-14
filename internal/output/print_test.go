package output

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testData = struct {
		Field1 string
		Field2 int
		Field3 interface{}
	}{
		Field1: "field1",
		Field2: 2,
		Field3: nil,
	}

	testDataJSON = "{\n    \"Field1\": \"field1\",\n    \"Field2\": 2,\n    \"Field3\": null\n}"
	testDataYAML = "field1: field1\nfield2: 2\nfield3: null\n"

	testHeaders = []string{"NAME", "DESCRIPTION"}
	testRows    = [][]string{
		{"testName1", "description1"},
		{"testName2", "description2"},
	}
	testWrongRows = [][]string{
		{"testName1", "description1"},
		{"testName2"},
	}

	testRowsRaw = "    NAME    | DESCRIPTION   \n------------|---------------\n  testName1 | description1  \n  testName2 | description2  \n"
)

type mockWriter struct {
	val    []byte
	rInt   int
	rError error
}

func (w *mockWriter) Write(b []byte) (int, error) {
	w.val = append(w.val, b...)
	return w.rInt, w.rError
}

func Test_printText(t *testing.T) {
	t.Run("print text", func(t *testing.T) {
		w := &bytes.Buffer{}
		assert.NoError(t, printText(w, "simple text"))
		assert.Equal(t, "simple text", w.String())
	})

	t.Run("write error", func(t *testing.T) {
		w := &mockWriter{
			rError: errors.New("sample error"),
		}
		assert.Error(t, printText(w, "simple text"))
	})
}

func Test_printJson(t *testing.T) {
	t.Run("print json", func(t *testing.T) {
		w := &bytes.Buffer{}
		assert.NoError(t, printJson(w, &testData))
		assert.Equal(t, testDataJSON, w.String())
	})

	t.Run("write error", func(t *testing.T) {
		w := &mockWriter{
			rError: errors.New("sample error"),
		}
		assert.Error(t, printJson(w, &testData))
	})

	t.Run("marshal error", func(t *testing.T) {
		w := &bytes.Buffer{}
		assert.Error(t, printJson(w, new(chan int)))
	})
}

func Test_printYaml(t *testing.T) {
	t.Run("print yaml", func(t *testing.T) {
		w := &bytes.Buffer{}
		assert.NoError(t, printYaml(w, &testData))
		assert.Equal(t, testDataYAML, w.String())
	})

	t.Run("write error", func(t *testing.T) {
		w := &mockWriter{
			rError: errors.New("sample error"),
		}
		assert.Error(t, printYaml(w, &testData))
	})

	t.Run("marshal error", func(t *testing.T) {
		w := &bytes.Buffer{}
		assert.Panics(t, func() {
			printYaml(w, make(chan int))
		})
	})
}

func Test_printTable(t *testing.T) {
	t.Run("print table", func(t *testing.T) {
		w := &bytes.Buffer{}
		assert.NoError(t, printTable(w, testHeaders, testRows))
		assert.Equal(t, testRowsRaw, w.String())
	})

	t.Run("validation error", func(t *testing.T) {
		w := &bytes.Buffer{}
		assert.Error(t, printTable(w, testHeaders, testWrongRows))
	})
}
