package cache

import (
	"context"
	"encoding/json"
	"github.com/bytedance/douyincloud-configcenter-sdk-go/openapi"
	"log"
	"strconv"
	"time"

	error2 "github.com/bytedance/douyincloud-configcenter-sdk-go/error"
	"github.com/bytedance/douyincloud-configcenter-sdk-go/http"
)

type Ticker struct {
	StopChan chan bool
}

func NewTicker(cache *Cache, ccClient *http.Client, interval time.Duration) *Ticker {
	ticker := time.NewTicker(interval)
	stopChan := make(chan bool)
	restartChan := make(chan bool)
	go func(ticker *time.Ticker) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error: %v", err)
			}
			restartChan <- true
		}()
		for {
			select {
			case <-ticker.C:
				UpdateCache(cache, ccClient)
			case stop := <-stopChan:
				if stop {
					log.Println("Ticker Stop")
					return
				}
			}
		}
	}(ticker)

	go func(ticker *time.Ticker) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error: %v", err)
			}
			ticker.Stop()
		}()
		for {
			select {
			case <-restartChan:
				for {
					select {
					case <-ticker.C:
						UpdateCache(cache, ccClient)
					case stop := <-stopChan:
						if stop {
							log.Println("Ticker Stop")
							return
						}
					}
				}
			}
		}
	}(ticker)

	return &Ticker{StopChan: stopChan}
}

func UpdateCache(cache *Cache, ccClient *http.Client) {
	configVersion := cache.GetVersion()
	if configVersion == "" {
		configVersion = "0"
	}
	bodyStruct := openapi.GetConfigListReqBody{Version: configVersion}
	jsonByte, _ := json.Marshal(bodyStruct)
	body := string(jsonByte)

	respBody, _, err := ccClient.CtxHttpPostRaw(context.Background(), body, nil)
	if err != nil {
		log.Printf("resp err, err = %v", err)
		return
	}

	var resp openapi.GetConfigListResponse
	var httpResult openapi.HttpResp
	err = json.Unmarshal(respBody, &httpResult)
	if err != nil {
		log.Printf("json unmarshal err in refresh config, err: %v", err)
		return
	}
	resp = httpResult.Data
	code := httpResult.Code
	msg := httpResult.Msg

	if code != 0 {
		if code == error2.ConfigServiceNotExist {
			log.Println("Please check whether the config center is opened in douyin cloud.")
			return
		}
		if code == error2.AuthenticationFailed {
			log.Println("Permission denied. You have no permission to access the dyc config center. Please check whether the program is running in dyc cloud or in the ide with dyc plugin.")
			return
		}
		log.Printf("Get config from config center err, err: %v", msg)
		return
	}

	localVersion, _ := strconv.Atoi(configVersion)
	onlineVersion, _ := strconv.Atoi(resp.Version)

	if localVersion >= onlineVersion {
		return
	}

	items := make(map[string]*Item, len(resp.Kvs))
	for _, v := range resp.Kvs {
		items[v.Key] = &Item{
			Object: v.Value,
			Type:   v.Type,
		}
	}
	cache.Set(items)
	cache.SetVersion(resp.Version)
}
