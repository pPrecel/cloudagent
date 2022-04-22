package output

import (
	"errors"
	"testing"
)

type testData struct {
	Name      string `json:"name" header:"name"`
	Namespace string `json:"namespace" header:"namespace"`
}

var (
	name1Filter = "#(name==Name1)#"

	testDataTab = []testData{
		{
			Name:      "Name0",
			Namespace: "Namespace0",
		},
		{
			Name:      "Name1",
			Namespace: "Namespace1",
		},
		{
			Name:      "Name1",
			Namespace: "Namespace2",
		},
		{
			Name:      "Name1",
			Namespace: "Namespace3",
		},
		{
			Name:      "Name2",
			Namespace: "Namespace2",
		},
	}

	testDataString = `[
    {
        "name": "Name0",
        "namespace": "Namespace0"
    },
    {
        "name": "Name1",
        "namespace": "Namespace1"
    },
    {
        "name": "Name1",
        "namespace": "Namespace2"
    },
    {
        "name": "Name1",
        "namespace": "Namespace3"
    },
    {
        "name": "Name2",
        "namespace": "Namespace2"
    }
]`

	testDataName1String = `[{
        "name": "Name1",
        "namespace": "Namespace1"
    },{
        "name": "Name1",
        "namespace": "Namespace2"
    },{
        "name": "Name1",
        "namespace": "Namespace3"
    }]`

	testDataTable = ` ─────────── ──────────── 
│ NAME (5)  │ NAMESPACE  │
 ─────────── ──────────── 
│ Name0     │ Namespace0 │
│ Name1     │ Namespace1 │
│ Name1     │ Namespace2 │
│ Name1     │ Namespace3 │
│ Name2     │ Namespace2 │
 ─────────── ──────────── 
`
	testDataName1Table = ` ─────── ──────────── 
│ NAME  │ NAMESPACE  │
 ─────── ──────────── 
│ Name1 │ Namespace1 │
│ Name1 │ Namespace2 │
│ Name1 │ Namespace3 │
 ─────── ──────────── 
`
)

type writer struct {
	Val string
	Err error
}

func (w *writer) Write(b []byte) (int, error) {
	w.Val += string(b)
	return len(w.Val), w.Err
}

func TestPrintText(t *testing.T) {
	type args struct {
		w *writer
		v interface{}
		o TextOptions
		f []string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "print text",
			args: args{
				w: &writer{},
				v: testDataTab,
				o: TextOptions{
					Format: "%r/%h/%u/%a",
					RPath:  "#(name==Name0)#|#",
					HPath:  "#(name==Namespace2)#|#",
					UPath:  "#(name==Name3)#|#",
					APath:  "#",
				},
			},
			wantW:   "1/0/0/5",
			wantErr: false,
		},
		{
			name: "filters",
			args: args{
				w: &writer{},
				v: testDataTab,
				o: TextOptions{
					Format: "%r/%h/%u/%a",
					RPath:  "#(name==Name0)#|#",
					HPath:  "#(name==Namespace2)#|#",
					UPath:  "#(name==Name1)#|#",
					APath:  "#(name==Name1)#|#",
				},
				f: []string{
					"#(name==Name1)#",
					"#(namespace=Namespace2)#",
				},
			},
			wantW:   "0/0/1/3",
			wantErr: false,
		},
		{
			name: "marshal error",
			args: args{
				w: &writer{},
				v: make(chan int),
			},
			wantW:   "",
			wantErr: true,
		},
		{
			name: "write error",
			args: args{
				w: &writer{
					Err: errors.New("test error"),
				},
			},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PrintText(tt.args.w, tt.args.v, tt.args.o, tt.args.f...); (err != nil) != tt.wantErr {
				t.Errorf("PrintText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := tt.args.w.Val; gotW != tt.wantW {
				t.Errorf("PrintText() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPrintErrorText(t *testing.T) {
	type args struct {
		w *writer
		o ErrorOptions
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "print error",
			args: args{
				w: &writer{},
				o: ErrorOptions{
					Format: "ERR: %e",
					Error:  "test error",
				},
			},
			wantW:   "ERR: test error",
			wantErr: false,
		},
		{
			name: "write error",
			args: args{
				w: &writer{
					Err: errors.New("test error"),
				},
			},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PrintErrorText(tt.args.w, tt.args.o); (err != nil) != tt.wantErr {
				t.Errorf("PrintErrorText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := tt.args.w.Val; gotW != tt.wantW {
				t.Errorf("PrintErrorText() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestPrintJson(t *testing.T) {
	type args struct {
		w *writer
		v interface{}
		f []string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "print json",
			args: args{
				w: &writer{},
				v: testDataTab,
			},
			wantErr: false,
			wantW:   testDataString,
		},
		{
			name: "filter",
			args: args{
				w: &writer{},
				v: testDataTab,
				f: []string{
					name1Filter,
				},
			},
			wantErr: false,
			wantW:   testDataName1String,
		},
		{
			name: "marshal error",
			args: args{
				w: &writer{},
				v: make(chan int),
			},
			wantW:   "",
			wantErr: true,
		},
		{
			name: "write error",
			args: args{
				w: &writer{
					Err: errors.New("test error"),
				},
			},
			wantW:   "null",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PrintJson(tt.args.w, tt.args.v, tt.args.f...); (err != nil) != tt.wantErr {
				t.Errorf("PrintJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := tt.args.w.Val; gotW != tt.wantW {
				t.Errorf("PrintJson() = '%v', want '%v'", gotW, tt.wantW)
			}
		})
	}
}

func TestPrintTable(t *testing.T) {
	type args struct {
		w *writer
		v interface{}
		f []string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "print table",
			args: args{
				w: &writer{},
				v: testDataTab,
			},
			wantW:   testDataTable,
			wantErr: false,
		},
		{
			name: "filter",
			args: args{
				w: &writer{},
				v: testDataTab,
				f: []string{
					name1Filter,
				},
			},
			wantW:   testDataName1Table,
			wantErr: false,
		},
		{
			name: "marshal error",
			args: args{
				w: &writer{},
				v: make(chan int),
			},
			wantW:   "",
			wantErr: true,
		},
		{
			name: "unmarshal error",
			args: args{
				w: &writer{},
				v: testDataTab,
				f: []string{
					"#(", // wrong filter
				},
			},
			wantW:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PrintTable(tt.args.w, tt.args.v, tt.args.f...); (err != nil) != tt.wantErr {
				t.Errorf("PrintTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := tt.args.w.Val; gotW != tt.wantW {
				t.Errorf("PrintTable() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
