package formater

import (
	"github.com/pPrecel/cloudagent/internal/output"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
)

var _ output.Formater = &checkFormater{}

var (
	checkHeaders = []string{"PROJECT", "STATUS", "MESSAGE", "LAST UPDATE"}
)

type checkFormater struct {
	resp *cloud_agent.GardenerResponse
}

func NewCheck(resp *cloud_agent.GardenerResponse) output.Formater {
	return &checkFormater{
		resp: resp,
	}
}

func (f *checkFormater) YAML() interface{} {
	if f.resp == nil {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"shoots": mergeShoots(f.resp.ShootList),
	}
}

func (f *checkFormater) JSON() interface{} {
	if f.resp == nil {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"shoots": mergeShoots(f.resp.ShootList),
	}
}

func (f *checkFormater) Table() ([]string, [][]string) {
	rows := [][]string{}
	for i := range f.resp.ShootList {
		s := f.resp.ShootList[i]

		status := "OK"
		message := "OK"
		if s.Error != "" {
			status = "ERROR"
			message = s.Error
		}

		rows = append(rows, []string{i, status, message, s.Time.AsTime().String()})
	}
	return checkHeaders, rows
}

func (f *checkFormater) Text(outFormat, errFormat string) string {
	return ""
}
