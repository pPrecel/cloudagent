package darwin

import (
	"reflect"
	"testing"
)

func TestPlistBody(t *testing.T) {
	type args struct {
		programPath string
		args        []string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "without additional args",
			args: args{
				programPath: "/tmp/path",
				args:        []string{},
			},
			want: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:</string>
	</dict>
	<key>Label</key>
	<string>com.pPrecel.cloudagent.agent.plist</string>
	<key>ProgramArguments</key>
	<array>
		<string>/tmp/path</string>
		<string>serve</string>

	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>/tmp/cloud-agent.stdout</string>
</dict>
</plist>`),
		},
		{
			name: "without additional args",
			args: args{
				programPath: "/tmp/path",
				args: []string{
					"--any-flag=any",
					"--sec-flag=other",
				},
			},
			want: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:</string>
	</dict>
	<key>Label</key>
	<string>com.pPrecel.cloudagent.agent.plist</string>
	<key>ProgramArguments</key>
	<array>
		<string>/tmp/path</string>
		<string>serve</string>
		<string>--any-flag=any</string>
		<string>--sec-flag=other</string>

	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>/tmp/cloud-agent.stdout</string>
</dict>
</plist>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PlistBody(tt.args.programPath, tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlistBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
