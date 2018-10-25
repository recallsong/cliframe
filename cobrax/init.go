package cobrax

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

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
	// 配置文件默认目录
	CfgDir string
	// 不带后缀的配置文件名
	CfgFileName string
	// 环境变量前缀
	EnvPrefix string
	// 远程配置阅读器
	rmtCfgReader RmtCfgReader
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
	CfgDir = opts.CfgDir
	CfgFileName = opts.CfgFileName
	EnvPrefix = opts.EnvPrefix
	rmtCfgReader = opts.RmtCfgReader
	if !opts.NoDefaultSets {
		setUpCommand(rootCmd, Flags, opts)
	}
	if opts.Init != nil {
		opts.Init(rootCmd)
	}
	if opts.Run != nil {
		rootCmd.Run = func(cmd *cobra.Command, args []string) {
			if !opts.NoDefaultSets {
				InitCommand(opts.AppConfig)
			}
			opts.Run(cmd, args)
		}
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
	peristFs := cmd.PersistentFlags()
	peristFs.StringVar(&flags.CfgFile, "config", "", "Config file (default is "+opts.CfgDir+"/"+opts.CfgFileName+".json or .yml)")
	peristFs.BoolVar(&flags.Debug, "debug", false, "Run in debug mode")
	peristFs.StringVar(&flags.PprofAddr, "pprof_addr", "", "Address of listen for pprof http server")
	viper.BindPFlag("pprof_addr", peristFs.Lookup("pprof_addr"))
	viper.BindPFlag("debug", peristFs.Lookup("debug"))
	if opts.RmtCfgReader != nil {
		peristFs.StringVar(&flags.RmtCfg, "rmt_cfg", "", "URL of remote config (example \"etcd://127.0.0.1:2379/path\")")
		viper.BindPFlag("rmt_cfg", peristFs.Lookup("rmt_cfg"))
	}
}

func InitCommand(cfgOutput interface{}) {
	readConfig(cfgOutput)
	//日志初始化
	Flags.Debug = viper.GetBool("debug")
	if Flags.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log_cfg := viper.GetStringMap("logs")
	if log_cfg != nil {
		if Flags.Debug {
			cfg := dic.FromMap(log_cfg)
			out, err := cfg.GetDic("out")
			if err == nil && out != nil {
				out.Set("name", "stdout")
			}
			cfg.Set("level", "DEBUG")
		}
		initLogger(log_cfg)
	}
	// 启用pprof
	Flags.PprofAddr = viper.GetString("pprof_addr")
	if Flags.PprofAddr != "" {
		log.Info("starting pprof server at ", Flags.PprofAddr)
		go func() {
			log.Error("fail to serve pprof server : ", http.ListenAndServe(Flags.PprofAddr, nil))
		}()
	}
}

func readConfig(cfgOutput interface{}) {
	// 设置可以从环境变量中读取信息
	if len(EnvPrefix) > 0 {
		viper.SetEnvPrefix(EnvPrefix)
		viper.AutomaticEnv()
		if cfgOutput != nil {
			typ := reflect.TypeOf(cfgOutput)
			for typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			for i, num := 0, typ.NumField(); i < num; i++ {
				field := typ.Field(i)
				key := field.Tag.Get("mapstructure")
				if len(key) == 0 {
					key = field.Name
				}
				key = strings.ToUpper(key)
				viper.BindEnv(key)
			}
		}
	}

	// 设置配置文件路径
	if Flags.CfgFile != "" {
		viper.SetConfigFile(Flags.CfgFile)
	} else if CfgDir != "" && CfgFileName != "" {
		viper.AddConfigPath(CfgDir)
		viper.SetConfigName(CfgFileName)
	}

	// 读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal(err)
		}
	} else {
		log.Debug("using config file:", viper.ConfigFileUsed())
	}

	// 从远程读取配置信息
	if rmtCfgReader != nil {
		Flags.RmtCfg = viper.GetString("rmt_cfg")
		if len(Flags.RmtCfg) > 0 {
			err := rmtCfgReader(Flags.RmtCfg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// 输出配置信息
	if cfgOutput != nil {
		GetConfig(cfgOutput)
	}
}

func GetConfig(out interface{}) {
	err := viper.Unmarshal(out)
	if err != nil {
		log.Fatal(err)
	}
}
