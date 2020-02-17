package setting

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	PageSize        int
	JwtSecret       string
	RuntimeRootPath string
	PrefixUrl       string
	LogSavePath     string
	ExportSavePath  string
	ImageSavePath   string
	ImageMaxSize    int
	ImageAllowExts  []string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Type         string
	User         string
	Password     string
	Host         string
	Name         string
	TablePrefix  string
	MaxIdleConns int
	MaxOpenConns int
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var AppSetting = &App{}
var ServerSetting = &Server{}
var DatabaseSetting = &Database{}
var RedisSetting = &Redis{}

func Setup(iniPath string) {

	if iniPath == "" {
		iniPath = "conf/app.ini"
	}

	cfg, err := ini.Load(iniPath)
	if err != nil {
		log.Fatalf("Fail to load '%s': %v", iniPath, err)
	}

	err = cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("MapTo AppSetting err: %v", err)
	}

	err = cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("MapTo ServerSetting err: %v", err)
	}

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	err = cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("MapTo DatabaseSetting err: %v", err)
	}

	err = cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("MapTo ServerSetting err: %v", err)
	}

	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

func PrintSetting() {
	bufAppSetting, _ := json.MarshalIndent(AppSetting, "", " ")
	log.Printf("AppSetting: \n%s\n", string(bufAppSetting))

	bufServerSetting, _ := json.MarshalIndent(ServerSetting, "", " ")
	log.Printf("ServerSetting: \n%s\n", string(bufServerSetting))

	bufDatabaseSetting, _ := json.MarshalIndent(DatabaseSetting, "", " ")
	log.Printf("DatabaseSetting: \n%s\n", string(bufDatabaseSetting))

	bufRedisSetting, _ := json.MarshalIndent(RedisSetting, "", " ")
	log.Printf("RedisSetting: \n%s\n", string(bufRedisSetting))
}
