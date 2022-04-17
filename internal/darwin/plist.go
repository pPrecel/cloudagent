package darwin

import "fmt"

const (
	plistBodyFormat = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:</string>
	</dict>
	<key>Label</key>
	<string>com.pPrecel.gardenagent.agent.plist</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
		<string>serve</string>
%s
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>/tmp/gardener-agent.stdout</string>
</dict>
</plist>`
)

func PlistBody(programPath string, args []string) []byte {
	programArguments := ""
	for i := range args {
		programArguments += fmt.Sprintf("		<string>%s</string>\n", args[i])
	}

	return []byte(fmt.Sprintf(plistBodyFormat, programPath, programArguments))
}
