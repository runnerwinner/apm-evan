package api

import (
	"context"
	"crypto/md5"
	"dogalarm/dao"
	"dogalarm/model"
	"dogalarm/notice"
	"dogapm"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/spf13/cast"
)

type alarm struct {
}

var Alarm = &alarm{}

func (a *alarm) MetricWebHook(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	dogapm.Logger.Info(context.TODO(), "receiveGrafanaMsg", map[string]any{
		"data": string(data),
	})

	alertModel := &model.AlertModel{}
	err = json.Unmarshal(data, alertModel)
	if err != nil {
		dogapm.Logger.Error(context.TODO(), "json unmarshal alert model fail", map[string]any{
			"data": string(data),
		},err)
		return
	}

	for _, alert := range alertModel.Alerts {
		app,host,msg := alert.Labels["app"],alert.Labels["host"],alert.Labels["alertname"]
		msg = fmt.Sprintf("content=%s host=%s app=%s value=%v", msg, host, app, alert.Values)
		deployInfo := dao.DeployInfo.GetInfoByApp(app)
		phoneWebhook := deployInfo["phone_webhook"].(string)
		phone := deployInfo["phone"].(string)
		if len(phoneWebhook)!=0 {
			success,_ := dogapm.Infra.Rdb.SetNX(context.TODO(), fmt.Sprintf("%s:phone:%s", "dogalarm", phone), 1, time.Minute).Result()
			if success {
				notice.Alarmer.Send(notice.Phone,msg,phoneWebhook, phone)
			}
		}
	}
}

func (a *alarm) LogWebHook(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}
	dogapm.Logger.Info(context.TODO(), "receiveLogStashMsg", map[string]any{
		"data": string(data),
	})

	contentMsg := make(map[string]any)
	err = json.Unmarshal(data, &contentMsg)
	if err != nil {
		dogapm.Logger.Error(context.TODO(), "json unmarshal logstash msg fail", map[string]any{
			"data": string(data),
		},err)
		return
	}

	app,host,msg := cast.ToString(contentMsg["app"]),cast.ToString(contentMsg["hostname"]),cast.ToString(contentMsg["message"])
	deployInfo := dao.DeployInfo.GetInfoByApp(app)
	if deployInfo["dingding_webhook"]==nil {
		return
	}
	dingdingWebhook := deployInfo["dingding_webhook"].(string)
	msg = fmt.Sprintf("content=%s host=%s app=%s", msg, host, app)
	msglen := len(msg)
	encryptStr := ""
	if msglen<1024 {
		encryptStr = getMD5(msg)
	} else {
		encryptStr = getMD5(msg[:1024])
	}
	limited := dogapm.RedisLimiter.IsLimit(dogapm.Infra.Rdb, fmt.Sprintf("%s:ding:%s:%d", "dogalarm",encryptStr,msglen),10,60)	
	if !limited {
		notice.Alarmer.Send(notice.DingDing,msg,dingdingWebhook,"")
	}
}

func getMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}