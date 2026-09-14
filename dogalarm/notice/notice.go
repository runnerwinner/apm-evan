package notice

import (
	"dogapm"
	"encoding/json"
	"net/http"
	"strings"
)

type alarm struct {
	client http.Client
}

var Alarmer = &alarm{}

type NoticeType int 

const (
	DingDing NoticeType = 1
	Phone NoticeType = 2
)

func (a *alarm) Send(noticeType NoticeType, msg string, webhookUrl, phone string) {
	switch noticeType {
	case DingDing:
		a.SendDingDing(msg, webhookUrl)
	case Phone:
		a.SendPhone(msg, webhookUrl, phone)
	}
}

func (a *alarm) SendDingDing(msg string, webhookUrl string) {
	jsonBody := map[string]any{
		"msgtype": "text",
		"text": map[string]string{
			"content": msg,
		},
	}
	data, _ := json.Marshal(jsonBody)
	req, _ := http.NewRequest("POST", webhookUrl, strings.NewReader(string(data)))
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		dogapm.Logger.Error(nil, "noticeDingDing", map[string]any{
			"url": webhookUrl,
		}, err)
		return
	}
	defer resp.Body.Close()
}

func (a *alarm) SendPhone(msg string, webhookUrl string, phone string) {
	jsonBody := map[string]any{
		"receiver": phone,
		"type":"phone",
		"title":"电话告警",
		"content":msg,
	}
	data, _:= json.Marshal(jsonBody)
	req, _ := http.NewRequest("POST", webhookUrl, strings.NewReader(string(data)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Servicetoken", "")
	resp,err := a.client.Do(req)
	if err != nil {
		dogapm.Logger.Error(nil, "noticePhone", map[string]any{
			"url": webhookUrl,
		},err)
		return
	}
	defer resp.Body.Close()
}