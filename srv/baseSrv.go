package srv

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/vislee/goframe/util"
	"hash/crc32"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

type baseSrv struct {
	cfgInfo *Config
	logs    *util.GfLogs
}

func (bs *baseSrv) GetRedisConn(host string, port int32) (redis.Conn, error) {
	addr := net.JoinHostPort(host, strconv.Itoa(int(port)))
	return redis.DialTimeout("tcp", addr, 2*time.Second, 2*time.Second, 1*time.Second)
}

func (bs *baseSrv) StrCrc32(s string, mod uint32) uint32 {
	return crc32.ChecksumIEEE([]byte(s)) % mod
}
