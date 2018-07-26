package rmtcfg

import (
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Read(rmtcfg string) error {
	// 解析url
	u, err := url.Parse(rmtcfg)
	if err != nil {
		return fmt.Errorf("fail to parse url of remote config : %v", err)
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme == "etcd" {
		addr := fmt.Sprintf("http://%s", u.Host)
		log.Infof("k/v store is %s, %s , %s", scheme, addr, u.Path)
		viper.AddRemoteProvider(scheme, addr, u.Path)
	} else {
		// TODO
		return fmt.Errorf("remote config store %s is not support", scheme)
	}
	typ := "json"
	if idx := strings.LastIndex(u.Path, "."); idx >= 0 {
		typ = u.Path[1:]
	}
	viper.SetConfigType(typ)

	// 设置从k/v系统中读取配置信息
	if err := viper.ReadRemoteConfig(); err != nil {
		if _, ok := err.(viper.RemoteConfigError); ok {
			log.Warn("read remote config not found")
			return nil
		}
		return err
	}
	log.Infof("read remote config from %s successfully", scheme)
	return nil
}
