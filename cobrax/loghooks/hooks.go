package loghooks

import (
	_ "github.com/gogap/logrus_mate/hooks/airbrake"
	_ "github.com/gogap/logrus_mate/hooks/graylog"
	_ "github.com/gogap/logrus_mate/hooks/lfshook"
	_ "github.com/gogap/logrus_mate/hooks/mail"
	_ "github.com/gogap/logrus_mate/hooks/syslog"
)
