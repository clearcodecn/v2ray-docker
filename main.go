package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Vmess struct {
	V    int    `json:"v"`
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	Id   string `json:"id"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	Path string `json:"path"`
}

var config string

func init() {
	flag.StringVar(&config, "c", "config.json", "config file")
}

func main() {
	flag.Parse()

	ip := getIP()

	data, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal(err)
	}
	var vc V2rayConfig
	err = json.Unmarshal(data, &vc)
	if err != nil {
		log.Fatal(err)
	}
	var vmess Vmess
	vmess.V = 2
	vmess.Ps = "v2ray-" + ip + "-" + fmt.Sprintf("%d", vc.Inbounds[0].Port)
	vmess.Add = ip
	vmess.Port = fmt.Sprintf("%d", vc.Inbounds[0].Port)
	vmess.Id = vc.Inbounds[0].Settings.Clients[0].Id
	vmess.Aid = fmt.Sprintf("%d", vc.Inbounds[0].Settings.Clients[0].AlterId)
	vmess.Net = vc.Inbounds[0].StreamSettings.Network
	vmess.Type = vc.Inbounds[0].Protocol
	vmess.Path = vc.Inbounds[0].Sniffing.DestOverride[0]

	data, _ = json.Marshal(vmess)
	fmt.Println("vmess://" + string(data))
}

func getIP() string {
	resp, _ := http.Get("http://ipinfo.io")
	if resp != nil {
		data, _ := ioutil.ReadAll(resp.Body)
		var info IPInfo
		json.Unmarshal(data, &info)
		return info.Ip
	}
	return ""
}

type IPInfo struct {
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

/**
  is_vmess_url=$(jq -c '{v:2,ps:'\"233boy-${net}-$is_addr\"',add:'\"$is_addr\"',port:'\"$port\"',
id:'\"$uuid\"',aid:"0",net:'\"$net\"',type:'\"$header_type\"',path:'\"$kcp_seed\"'}' <<<{})
  is_url=vmess://$(echo -n $is_vmess_url | base64 -w 0)
*/

type V2rayConfig struct {
	Log struct {
		Loglevel string `json:"loglevel"`
	} `json:"log"`
	Inbounds []struct {
		Port     int    `json:"port"`
		Protocol string `json:"protocol"`
		Settings struct {
			Clients []struct {
				Id      string `json:"id"`
				Level   int    `json:"level"`
				AlterId int    `json:"alterId"`
			} `json:"clients"`
		} `json:"settings"`
		StreamSettings struct {
			Network string `json:"network"`
		} `json:"streamSettings"`
		Sniffing struct {
			Enabled      bool     `json:"enabled"`
			DestOverride []string `json:"destOverride"`
		} `json:"sniffing"`
	} `json:"inbounds"`
	Outbounds []struct {
		Protocol string `json:"protocol"`
		Settings struct {
			DomainStrategy string `json:"domainStrategy,omitempty"`
		} `json:"settings"`
		Tag string `json:"tag"`
	} `json:"outbounds"`
	Dns struct {
		Servers []string `json:"servers"`
	} `json:"dns"`
	Routing struct {
		DomainStrategy string `json:"domainStrategy"`
		Rules          []struct {
			Type        string   `json:"type"`
			Ip          []string `json:"ip,omitempty"`
			OutboundTag string   `json:"outboundTag"`
			Domain      []string `json:"domain,omitempty"`
			Protocol    []string `json:"protocol,omitempty"`
		} `json:"rules"`
	} `json:"routing"`
	Transport struct {
		KcpSettings struct {
			UplinkCapacity   int  `json:"uplinkCapacity"`
			DownlinkCapacity int  `json:"downlinkCapacity"`
			Congestion       bool `json:"congestion"`
		} `json:"kcpSettings"`
	} `json:"transport"`
}
