package formater

import (
	"strconv"
	"strings"

	"github.com/pPrecel/cloudagent/internal/output"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
)

const (
	CheckTextAllFormat        = "$a"
	CheckTextErrorCountFormat = "$e"
	CheckTextHealthyFormat    = "$h"
	CheckTextErrorFormat      = "$E"
)

var _ output.Formater = &checkFormater{}

var (
	checkHeaders = []string{"PROJECT", "STATUS", "MESSAGE", "LAST UPDATE", "PROVIDER"}

	checkDirectives = checkDirectiveMap{
		CheckTextAllFormat: func(r *cloud_agent.GardenerResponse) string {
			return strconv.Itoa(len(r.ShootList))
		},
		CheckTextErrorCountFormat: func(r *cloud_agent.GardenerResponse) string {
			e := 0
			for i := range r.ShootList {
				if r.ShootList[i].Error != "" {
					e++
				}
			}

			return strconv.Itoa(e)
		},
		CheckTextHealthyFormat: func(r *cloud_agent.GardenerResponse) string {
			h := 0
			for i := range r.ShootList {
				if r.ShootList[i].Error == "" {
					h++
				}
			}

			return strconv.Itoa(h)
		},
		CheckTextErrorFormat: func(r *cloud_agent.GardenerResponse) string {
			e := []string{}
			if r.GeneralError != "" {
				e = append(e, r.GeneralError)
			}

			for i := range r.ShootList {
				err := r.ShootList[i].Error
				if err != "" {
					e = append(e, r.ShootList[i].Error)
				}

			}

			return strings.Join(e, ", ")
		},
	}
)

type checkDirectiveMap map[string]func(*cloud_agent.GardenerResponse) string

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
		"generalError": f.resp.GeneralError,
		"shoots":       mergeShoots(f.resp.ShootList),
	}
}

func (f *checkFormater) JSON() interface{} {
	if f.resp == nil {
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"generalError": f.resp.GeneralError,
		"shoots":       mergeShoots(f.resp.ShootList),
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

		rows = append(rows, []string{
			i,
			status,
			message,
			s.Time.AsTime().Local().Format("2006-01-02 15:04:05"),
			gardenerProvider,
		})
	}
	return checkHeaders, rows
}

func (f *checkFormater) Text(outFormat, errFormat string) string {
	if f.resp == nil {
		err := "nil response"
		return strings.ReplaceAll(errFormat, CheckTextErrorFormat, err)
	}

	str := outFormat
	for key, val := range checkDirectives {
		str = strings.ReplaceAll(str, key, val(f.resp))
	}

	return str
}
