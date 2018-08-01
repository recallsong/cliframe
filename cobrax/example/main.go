package main

import (
	"github.com/recallsong/cliframe/cobrax"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Toggle    bool
	Name      string
	OtherConf string `mapstructure:"other_conf"`
}

func main() {
	cfg := &AppConfig{}
	var subCmd = &cobra.Command{
		Use: "subcmd",
		Run: func(cmd *cobra.Command, args []string) {
			cobrax.InitCommand(cfg)
			log.Info("toggle: ", cfg.Toggle)
			log.Info("name: ", cfg.Name)
			log.Info("other: ", cfg.OtherConf)
		},
	}
	cobrax.Execute("example", &cobrax.Options{
		CfgDir:      ".",
		CfgFileName: "example",
		AppConfig:   cfg,
		// RmtCfgReader: rmtcfg.Read, Add remote config support
		Init: func(cmd *cobra.Command) {
			fs := cmd.Flags()
			fs.BoolVar(&cfg.Toggle, "toggle", true, "watch store, default is true")
			fs.StringVarP(&cfg.Name, "name", "n", "Ruiguo", "it's my name")
			fs.StringVar(&cfg.OtherConf, "other_conf", "", "other config")
			viper.BindPFlags(fs) // bind to pflag
			cmd.AddCommand(subCmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Info("toggle: ", cfg.Toggle)
			log.Info("name: ", cfg.Name)
			log.Info("other: ", cfg.OtherConf)
		},
	})
}
