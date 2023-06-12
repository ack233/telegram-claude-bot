package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// viperConfig 全局配置变量
var (
	ViperConfig   *viper.Viper
	BotConfig     BotconfigStruct
	ForbiddenWord ForbiddenWordStruct
	Claudeconfig  ClaudeconfigStruct

	validate = validator.New()
)

// 初始化配置文件相关设置，在 main 包中调用进行初始化加载
func Init() {
	//安位置运算log.Lshortfile | log.LstdFlags: 10000 | 11 = 10011 (二进制)，即 19 (十进制)
	//log底层推导l.flag & log.LstdFlags != 0 时即为true
	//log.SetFlags(log.Lshortfile | log.LstdFlags) //设置日志行号
	log.SetFlags(log.LstdFlags) //设置日志行号

	//从环境变量获取配置文件路径
	configPath := getEnvConfigPath()

	var configfile = pflag.StringP("configfile", "c", configPath, "user -c set Your congfile")
	pflag.StringP("logfile", "f", "./bot.log", "user -f set Your logfile")
	pflag.StringP("loglevel", "l", "INFO", "user -l set loglevel")
	pflag.ErrHelp = errors.New("") //替换为一个空字符串。这意味着，当用户请求帮助信息时，程序将不会打印错误消息，只会显示帮助文本
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
		fmt.Fprint(os.Stderr, "优先级: 命令行 > 配置文件 > 命令行默认值")
	}

	pflag.Parse()
	if pflag.NArg() != 0 {
		fmt.Printf("参数 %v 无法解析,请参考以下语法:\n", pflag.Args())
		pflag.Usage()
		fmt.Print("\n")
		os.Exit(1)
	} //
	if !Exists(*configfile) {
		log.Fatalf("no such config directory: %s", *configfile)
		checkerr(errors.New("no such config directory: " + *configfile))
	}
	viper.SetConfigFile(*configfile) // path to look for the config file in
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		checkerr(err)
	}
	viper.BindPFlags(pflag.CommandLine)
	ViperConfig = viper.GetViper()

	if err = validateConfig(ViperConfig); err != nil {
		checkerr(err)
	}

	//err = ViperConfig.UnmarshalKey("botconfig", &BotConfig)
	//checkerr(err)
	//err = ViperConfig.UnmarshalKey("claudeconfig", &Claudeconfig)
	//checkerr(err)
	//err = ViperConfig.UnmarshalKey("forbiddenWord", &ForbiddenWord)
	//checkerr(err)
	err = ViperConfig.Unmarshal(&config{
		&BotConfig,
		&Claudeconfig,
		&ForbiddenWord,
	})
	checkerr(err)

	validateStruct(BotConfig, Claudeconfig)

}

func getEnvConfigPath() string {
	configPath := os.Getenv("CONFIG_FILE")
	if configPath == "" {
		configPath = "./config.yaml" // Replace with your actual default path
	}
	return configPath
}

func checkerr(err error) {
	_, file, line, _ := runtime.Caller(1)
	if err != nil {
		log.Fatalf("%v: %v,%v", file, line, err)
	}
}

func validateStruct(cfg ...interface{}) {
	for _, c := range cfg {
		if err := validate.Struct(c); err != nil {
			checkerr(err)
		}
	}
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil

	//if errors.Is(err, os.ErrNotExist) {
	//	return false
	//}
}

func validateConfig(v *viper.Viper) error {
	var (
		logdir = filepath.Dir(v.GetString("logfile"))
	)

	if !Exists(logdir) {
		return fmt.Errorf("no such directory: logdir: %s, please check configuration", logdir)
	}

	return nil
}
