/**
 * @File: acm
 * @Author: Shaw
 * @Date: 2020/6/27 2:27 PM
 * @Desc

 */

package acm

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
	"strings"
)

type Acm struct {
	constant.ClientConfig

	client config_client.IConfigClient

	DefaultGroup string //默认的分组

	Enable bool //是否启用，调用了DefaultAcm表示启用，不开启的情况返回默认值
}

func DefaultDisableAcm() (acm *Acm) {
	acm = &Acm{
		Enable: false,
	}
	return acm
}

func DefaultAcm(accessKey, secretKey, namespaceId, defaultGroup string) (acm *Acm, err error) {
	if accessKey == "" || secretKey == "" || namespaceId == "" || defaultGroup == "" {
		return nil, fmt.Errorf("illegal arguments")
	}
	// 从控制台命名空间管理的"命名空间详情"中拷贝 End Point、命名空间 ID
	var endpoint = "acm.aliyun.com"

	acm = &Acm{
		ClientConfig: constant.ClientConfig{
			Endpoint:       endpoint + ":8080",
			NamespaceId:    namespaceId,
			AccessKey:      accessKey,
			SecretKey:      secretKey,
			TimeoutMs:      5 * 1000,
			ListenInterval: 30 * 1000,
		},
		DefaultGroup: defaultGroup,
		Enable:       true,
	}

	acm.client, err = clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": acm.ClientConfig,
	})
	if err != nil {
		logs.Error(err)
		return
	}

	return
}

//根据key(即dataId)获取，如果为空则设置，并返回默认值
//没有开启Acm情况下直接返回默认值
func (a *Acm) GetSetString(key string, defaultConfig string) (val string, err error) {
	if !a.Enable {
		return defaultConfig, nil
	}

	val, err = a.client.GetConfig(vo.ConfigParam{
		DataId: key,
		Group:  a.DefaultGroup})
	if err != nil && !strings.Contains(err.Error(), "config not found") {
		val = defaultConfig
		logs.Error(err)
		return
	}
	if defaultConfig != "" && err != nil && strings.Contains(err.Error(), "config not found") {
		logs.Debug("set acm key: %s value: %s", key, defaultConfig)
		val = defaultConfig
		var success bool
		success, err = a.client.PublishConfig(vo.ConfigParam{
			DataId:  key,
			Group:   a.DefaultGroup,
			Content: defaultConfig,
		})
		if err != nil {
			val = defaultConfig
			logs.Error(err)
			return
		}
		if !success {
			err = fmt.Errorf("set Group:%s DataId:%s failed", a.DefaultGroup, key)
			logs.Error(err)
			return
		}
	}
	return
}

func (a *Acm) GetSetInt64(dataId string, defaultConfig int64) (val int64, err error) {
	if !a.Enable {
		return defaultConfig, nil
	}
	content, err := a.GetSetString(dataId, fmt.Sprint(defaultConfig))
	if err != nil {
		val = defaultConfig
		return
	}
	val, err = strconv.ParseInt(content, 10, 64)
	if err != nil {
		val = defaultConfig
		return
	}
	return
}

func (a *Acm) GetClient() (client config_client.IConfigClient, err error) {
	if !a.Enable {
		return nil, fmt.Errorf("not enabled acm")
	}
	client = a.client
	return
}
