package liveprobe

import (
	"dogalarm/dao"
	"dogalarm/notice"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cast"
)
type probe struct {	
	client http.Client
}

var Probe = &probe{}

func getAppInfo(appCopy map[string]any) (port int, ips []string, appName string, liveProbePath string) {
	ips = strings.Split(cast.ToString(appCopy["hosts"]), ",")
	appName = cast.ToString(appCopy["app"])
	port = cast.ToInt(appCopy["port"])
	liveProbePath = cast.ToString(appCopy["probe_url"])
	return port, ips, appName, liveProbePath
}


func (p *probe) Enable() {
	go func ()  {
		for {
			apps := dao.DeployInfo.All()
			for _,app := range apps {
				port, ips, appName, liveProbePath := getAppInfo(app)
				for _,ip := range ips {
					checkUrl := fmt.Sprintf("http://%s:%d/%s", ip, port, liveProbePath)
					p.checkLive(checkUrl,appName, ip,app["phone_webhook"].(string),app["phone"].(string), 3)
				}
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()
}

func (p *probe) checkLive(checkUrl, appName, host, alarmUrl, phone string, retryCnt int) bool {
	for i := 0; i < retryCnt; i++ {
        resp, err := p.client.Get(checkUrl)
        if err != nil {
            time.Sleep(time.Second)
            continue
        }

        if resp != nil && resp.Body != nil {
            if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
                resp.Body.Close()
                return true
            }
            resp.Body.Close()
        }

        time.Sleep(time.Second)
    }
	notice.Alarmer.Send(notice.Phone, fmt.Sprintf("app=%s host=%s 探测失败，服务宕机", appName, host), alarmUrl, phone)
    return false
}