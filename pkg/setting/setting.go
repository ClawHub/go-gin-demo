package setting

import (
	"go-gin-demo/pkg/logging"
	"go.uber.org/zap"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	ExportSavePath string
	QrCodeSavePath string
	FontSavePath   string
}

var AppSetting = &App{}

//服务相关
type Server struct {
	ProjectName  string
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

var cfg *ini.File

func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		logging.AppLogger.Fatal("setting.Setup, fail to parse 'conf/app.ini' ", zap.Error(err))
	}

	MapTo("app", AppSetting)
	MapTo("server", ServerSetting)

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
}

func MapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		logging.AppLogger.Fatal("Cfg.MapTo Setting err' ", zap.Error(err))
	}
}
