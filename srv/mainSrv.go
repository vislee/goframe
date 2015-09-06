//Package srv the services
package srv

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vislee/goframe/util"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type MainSrv struct {
	baseSrv
	QuitCh chan bool
	s      util.GoSrv
	usrDb  *sql.DB
	http   *httpSrv

	//other service
	//pass

}

//NewMainSrv return the MainSrv instance
func NewMainSrv(cfg *Config) *MainSrv {
	if bTyp != TypeFreq && bTyp != TypeStor {
		return nil
	}

	ms := MainSrv{baseSrv: baseSrv{cfgInfo: cfg}}
	ms.QuitCh = make(chan bool)

	return &ms
}

func (ms *MainSrv) Init() error {
	var err error

	ms.logs, err = util.NewLogs(ms.cfgInfo.Log.Dir, ms.cfgInfo.Log.Prefix, ms.cfgInfo.Log.Level)
	if err != nil {
		return err
	}

	usrDsn := ms.cfgInfo.UsrMysql.Dsn()
	ms.logs.Debug("dsn: %s", usrDsn)

	ms.usrDb, err = sql.Open("mysql", usrDsn)
	if err != nil {
		ms.logs.Error("init mysql error. error: %s", err.Error())
		return err
	}
	ms.http = NewHttpSrv(ms.cfgInfo, ms.logs)
	//other service
	//pass

	return nil
}

//Main  service main function
func (ms *MainSrv) Main() {
	ms.logs.Debug("MainSrv Main ...")
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := ms.http.InitListener()
	if err == nil {
		ms.s.Wrap(ms.http.httpMain)
	}
	//other service
	//pass

	// <-ms.QuitCh
}

//Stop  service stop
func (ms *MainSrv) Stop() {
	ms.logs.Debug("MainSrv Stop ...")
	ms.UsrDb.Close()
	ms.http.Close()
	//other service
	//pass

	close(ms.QuitCh)
	ms.s.Wait()
}
