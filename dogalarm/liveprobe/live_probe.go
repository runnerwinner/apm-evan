package liveprobe

import (
	"dogalarm/dao"
	"dogalarm/metric"
	"dogalarm/notice"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"
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
				metric.LiveProbeGuage.WithLabelValues(appName, host).Set(float64(metric.Living))
                return true
            }
            resp.Body.Close()
        }

        time.Sleep(time.Second)
    }

	remoteStartApp(appName, host)
	time.Sleep(2*time.Second)
	resp, err := p.client.Get(checkUrl)
	if err != nil || resp.StatusCode != http.StatusOK {
		notice.Alarmer.Send(notice.Phone, fmt.Sprintf("app=%s host=%s 探测失败，服务宕机", appName, host), alarmUrl, phone)
		metric.LiveProbeGuage.WithLabelValues(appName, host).Set(float64(metric.ShutDown))
		return false
	}else{
		return true
	}
	
}

func remoteStartApp(app string, ip string) {

	sshUser := "root"
	sshKeyPath := "/Users/evan/.ssh/id_rsa"
	sshProt :=10022 // 映射远程docker的ssh 22端口
	config := &ssh.ClientConfig{
		User: sshUser,
		Timeout: time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, sshProt), config)
	if err != nil {
		return
	}
	defer sshClient.Close()

	session,err := sshClient.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	// 执行远程命令
	// cmd := fmt.Sprintf("cd /root/%s && ./run.sh %s 1>>nohup.out 2>&1 &", app, app)
	// err = session.Run(cmd)
	combo, err := session.CombinedOutput(fmt.Sprintf("cd /root/%s && ./run.sh %s 1>>nohup.out 2>&1 &", app, app))
	if err != nil {
		log.Println("远程执行cmd 失败", err, string(combo))
		return
	}
	log.Println("命令输出", string(combo))

}

func publicKeyAuthFunc(keypath string) ssh.AuthMethod {
	keypath, err := homedir.Expand(keypath)
	if err != nil { 
		panic(err)
	}
	key, err := ioutil.ReadFile(keypath)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}