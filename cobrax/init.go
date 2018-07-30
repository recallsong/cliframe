package cobrax

import (
	"net/http"
	"runtime"

	"github.com/recallsong/go-utils/container/dic"
	vercmd "github.com/recallsong/go-utils/version/cobra-vercmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// 定义通用参数
type ComFlags struct {
	CfgFile   string
	RmtCfg    string
	PprofAddr string
	Debug     bool
}

var (
	Flags   = &ComFlags{}
	rootCmd = &cobra.Command{}
)

type RmtCfgReader func(rmtcfg string) error

type Options struct {
	CfgDir        string
	CfgFileName   string
	EnvPrefix     string
	RmtCfgReader  RmtCfgReader
	AppConfig     interface{}
	NoDefaultSets bool
	Init          func(cmd *cobra.Command)
	Run           func(cmd *cobra.Command, args []string)
}

func Execute(name string, opts *Options) {
	if len(name) <= 0 {
		panic("app name should not be empty")
	}
	rootCmd.Use = name
	vercmd.AddTo(rootCmd)
	if !opts.NoDefaultSets {
		setUpCommand(rootCmd, Flags, opts)
	}
	opts.Init(rootCmd)
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if !opts.NoDefaultSets {
			initCommand(rootCmd, Flags, opts)
			if opts.AppConfig != nil {
				GetConfig(opts.AppConfig)
			}
		}
		opts.Run(cmd, args)
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func setUpCommand(cmd *cobra.Command, flags *ComFlags, opts *Options) {
	if len(opts.CfgDir) <= 0 {
		opts.CfgDir = "."
	}
	if len(opts.CfgFileName) <= 0 {
		opts.CfgFileName = cmd.Use
	}
	peristFs := rootCmd.PersistentFlags()
	peristFs.StringVar(&flags.CfgFile, "config", "", "config file (default is "+opts.CfgDir+"/"+opts.CfgFileName+".json)")
	peristFs.BoolVar(&flags.Debug, "debug", false, "run in debug mode")
	peristFs.StringVar(&flags.PprofAddr, "pprof_addr", "", "address of listen for pprof http server")
	viper.BindPFlag("pprof_addr", peristFs.Lookup("pprof_addr"))
	viper.BindPFlag("debug", peristFs.Lookup("debug"))
	if opts.RmtCfgReader != nil {
		peristFs.StringVar(&flags.RmtCfg, "rmt_cfg", "", "url of remote config (example \"etcd://127.0.0.1:2379/path\")")
		viper.BindPFlag("rmt_cfg", peristFs.Lookup("rmt_cfg"))
	}
}

func initCommand(cmd *cobra.Command, flags *ComFlags, opts *Options) {
	readConfig(cmd, flags, opts)
	//日志初始化
	flags.Debug = viper.GetBool("debug")
	if flags.Debug {
		log.SetLevel(log.DebugLevel)
		log.Info("run in debug mode")
	}
	log_cfg := viper.GetStringMap("logs")
	if log_cfg != nil {
		if flags.Debug {
			out, err := dic.FromMap(log_cfg).GetDic("out")
			if err == nil && out != nil {
				out.Set("name", "stdout")
			}
		}
		initLogger(log_cfg)
	}
	// 启用pprof
	flags.PprofAddr = viper.GetString("pprof_addr")
	if flags.PprofAddr != "" {
		log.Info("starting pprof server at ", flags.PprofAddr)
		go func() {
			log.Error("fail to serve pprof server : ", http.ListenAndServe(flags.PprofAddr, nil))
		}()
	}
}

func readConfig(cmd *cobra.Command, flags *ComFlags, opts *Options) {
	// 设置可以从环境变量中读取信息
	if len(opts.EnvPrefix) > 0 {
		viper.SetEnvPrefix(opts.EnvPrefix)
		viper.AutomaticEnv()
	}

	// 设置配置文件路径
	if flags.CfgFile != "" {
		viper.SetConfigFile(flags.CfgFile)
	} else {
		viper.AddConfigPath(opts.CfgDir)
		viper.SetConfigName(opts.CfgFileName)
	}

	// 读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	} else {
		log.Info("using config file:", viper.ConfigFileUsed())
	}

	// 从远程读取配置信息
	if opts.RmtCfgReader != nil {
		flags.RmtCfg = viper.GetString("rmt_cfg")
		if len(flags.RmtCfg) > 0 {
			err := opts.RmtCfgReader(flags.RmtCfg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func GetConfig(out interface{}) {
	err := viper.Unmarshal(out)
	if err != nil {
		log.Fatal(err)
	}
}
