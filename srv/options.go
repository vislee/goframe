package srv

import (
	"fmt"
)

type urls struct {
	AlertUrl   string
	DisableUrl string
}

type logs struct {
	Dir    string
	Prefix string
	Level  int
}

type mmhost struct {
	Host string
	Port int32
}

type mysql struct {
	Host string
	Port int32
	User string
	Pass string
	Db   string
}

func (d *mysql) Dsn() string {
	netAddr := fmt.Sprintf("tcp(%s:%d)", d.Host, d.Port)
	return fmt.Sprintf("%s:%s@%s/%s?timeout=5s&strict=true&allowOldPasswords=1", d.User, d.Pass, netAddr, d.Db)
}

type uLevelMc struct {
	mmhost
	UlPrefix string
}

type redisHost struct {
	Host   string
	Port   int32
	Suffix string
}

type cycle struct {
	BakPath string
	Funcs   []string
	Args    []interface{}
}

type resMap struct {
	MapFreq []string
	MapStor []string
}

type listenHost struct {
	Host string
}

type Config struct {
	UsrMysql   mysql
	MsgRedis   redisHost
	Urls       urls
	Log        logs
	UserLevel  uLevelMc
	Cycle      map[string]cycle
	HttpListen map[string]listenHost
}
