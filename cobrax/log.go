package cobrax

import (
	"github.com/gogap/logrus_mate"
	_ "github.com/gogap/logrus_mate/writers/rotatelogs"
	"github.com/recallsong/go-utils/encoding/jsonx"
	"github.com/sirupsen/logrus"
)

func initLogger(cfg map[string]interface{}) {
	logrus_mate.Hijack(
		logrus.StandardLogger(),
		logrus_mate.ConfigString(jsonx.Marshal(cfg)),
	)
}
