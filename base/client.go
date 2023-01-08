package base

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/bytedance/douyincloud-configcenter-sdk-go/cache"
	error2 "github.com/bytedance/douyincloud-configcenter-sdk-go/error"
	"github.com/bytedance/douyincloud-configcenter-sdk-go/http"
	"github.com/bytedance/douyincloud-configcenter-sdk-go/openapi"
)

type Client interface {
	Get(key string) (string, error)
	RefreshConfig() error
}

type internalClient struct {
	cache    *cache.Cache
	ccClient *http.Client
	ticker   *cache.Ticker
}

func Start() (Client, error) {
	config := NewClientConfig()
	return StartWithConfig(config)
}

func StartWithConfig(config *clientConfig) (Client, error) {
	client := &internalClient{}

	client.ccClient = http.NewClient(func(options *http.Options) {
		options.Timeout = config.GetTimeout()
	})
	client.cache = cache.NewCache()

	err := updateCache(client.cache, client.ccClient)
	if err != nil {
		log.Printf("first update cache err, err = %v", err)
		return nil, err
	}

	client.ticker = cache.NewTicker(client.cache, client.ccClient, config.GetUpdateInterval())

	log.Println("sdk start finished!")

	return client, nil
}

func (c *internalClient) Get(key string) (string, error) {
	v, _, err := c.getWithCache(key)
	if err != nil {
		return "", err
	}
	value := v.Object.(string)
	return value, nil
}

func (c *internalClient) RefreshConfig() error {
	err := updateCache(c.cache, c.ccClient)
	if err != nil {
		log.Printf("update cache err, err = %v", err)
		return err
	}
	return nil
}

func (c *internalClient) getWithCache(key string) (*cache.Item, bool, error) {
	item, exist := c.cache.Get(key)
	if !exist {
		return nil, false, errors.New("item not exist")
	}
	return item, true, nil

}

func updateCache(cache2 *cache.Cache, ccClient *http.Client) error {
	configVersion := cache2.GetVersion()
	if configVersion == "" {
		configVersion = "0"
	}
	bodyStruct := openapi.GetConfigListReqBody{Version: configVersion}
	jsonByte, _ := json.Marshal(bodyStruct)
	body := string(jsonByte)

	respBody, _, err := ccClient.CtxHttpPostRaw(context.Background(), body, nil)
	if err != nil {
		log.Printf("resp err in updateCache, err: %v", err)
		return err
	}

	var resp openapi.GetConfigListResponse
	var httpResult openapi.HttpResp
	err = json.Unmarshal(respBody, &httpResult)
	if err != nil {
		log.Printf("json unmarshal err in init config, err: %v", err)
		return err
	}
	resp = httpResult.Data
	code := httpResult.Code
	msg := httpResult.Msg

	if code != 0 {
		if code == error2.ConfigServiceNotExist {
			log.Println("Please check whether the config center is opened in douyin cloud.")
			return errors.New("please check whether the config center is opened in douyin cloud")
		}
		if code == error2.AuthenticationFailed {
			log.Println("Permission denied. You have no permission to access the dyc config center. Please check whether the program is running in dyc cloud or in the ide with dyc plugin.")
			return errors.New("permission denied. You have no permission to access the dyc config center. Please check whether the program is running in dyc cloud or in the ide with dyc plugin")
		}
		log.Printf("Get config from config center err, err: %v", msg)
		return errors.New("Get config from config center err, err: " + msg)
	}

	localVersion, _ := strconv.Atoi(configVersion)
	onlineVersion, _ := strconv.Atoi(resp.Version)

	if localVersion >= onlineVersion {
		return nil
	}

	items := make(map[string]*cache.Item, len(resp.Kvs))
	for _, v := range resp.Kvs {
		items[v.Key] = &cache.Item{
			Object: v.Value,
			Type:   v.Type,
		}
	}
	cache2.Set(items)
	cache2.SetVersion(resp.Version)
	return nil
}
