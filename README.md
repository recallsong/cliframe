# cliframe
Command-Line Interface Framework. 

这是一个方便Go开发者快速开发CLI的程序的基础框架。

基于如下库进行封装
*   [Logrus](https://github.com/sirupsen/logrus)
*   [Logrus Mate](https://github.com/gogap/logrus_mate)
*   [Cobra](https://github.com/spf13/cobra)
*   [Viper](https://github.com/spf13/viper)

具有特色
*   从命令行读取flag参数，支持短参数名和长参数名
*   从 JSON, TOML, YAML, HCL, Java properties 等文件读取配置信息
*   从远程K/V存储服务器（etcd）上读取读取配置信息
*   从文件、远程K/V配置日志行为
*   默认具有version命令，debug、config等通用参数
*   提供Dockerfile、Makefile、配置文件的例子， [cobrax/example](https://github.com/recallsong/cliframe/tree/master/cobrax/example)

# Download

        go get github.com/recallsong/cliframe

# Quick Start
## Example
Cypy [cobrax/example](https://github.com/recallsong/cliframe/tree/master/cobrax/example) and change it as your application
```go
import (
    "github.com/recallsong/cliframe/cobrax"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    // "github.com/recallsong/cliframe/cliframe/rmtcfg" Add remote config support 
    // _ "github.com/recallsong/cliframe/cliframe/loghooks" Add more log outputs support
)

type AppConfig struct {
    Toggle    bool
    Name      string
    OtherConf string `mapstructure:"other_conf"`
}

func main() {
    cfg := &AppConfig{}
    cobrax.Execute("example", &cobrax.Options{
        CfgDir:      ".",
        CfgFileName: "example",
        // RmtCfgReader: rmtcfg.Read, Add remote config support
        AppConfig:   cfg,
        Init: func(cmd *cobra.Command) {
            fs := cmd.Flags()
            fs.BoolVar(&cfg.Toggle, "toggle", true, "watch store, default is true")
            fs.StringVarP(&cfg.Name, "name", "n", "Ruiguo", "it's my name")
            fs.StringVar(&cfg.OtherConf, "other_conf", "", "other config")
            viper.BindPFlags(fs) // bind to pflag
        },
        Run: func(cmd *cobra.Command, args []string) {
            log.Info("toggle: ", cfg.Toggle)
            log.Info("name: ", cfg.Name)
            log.Info("other: ", cfg.OtherConf)
        },
    })
}

```
## Run

    make run

Or

    make build
    ./example
    
Or run in docker container

    make docker-run

Output:

    INFO[0000] run in debug mode
    INFO[0000] toggle: true
    INFO[0000] name: Ruiguo
    INFO[0000] other: this parameter from flag

## Other Make Targets

- make cross-build
```
    GO_OS=linux GO_ARCH=amd64 make cross-build
```
- make docker-build
- make clean
- make run-help
- make run-version

# License
[MIT](https://github.com/recallsong/cliframe/blob/master/LICENSE)
