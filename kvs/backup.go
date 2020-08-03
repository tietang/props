package kvs

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

const (
	BackupFile = ".conf/all.properties"
)

type Backup interface {
	Restore()
	Backup()
}

var _ Backup = new(CompositeConfigSource)

type DiskBackup struct {
	ccs            CompositeConfigSource
	BackupFileName string
}

func (d *DiskBackup) Restore() {
	props := NewPropertiesConfigSource(d.BackupFileName)
	d.ccs.Add(props)
}

func (d *DiskBackup) Backup() {

	dir, _ := os.Getwd()
	now := time.Now()
	time := now.Format("20060102150405")
	d.BackupFileName = filepath.Join(dir, BackupFile)
	if PathExists(d.BackupFileName) {
		newFileName := filepath.Join(dir, BackupFile+"."+time)
		err := os.Rename(d.BackupFileName, newFileName)
		if err != nil {
			log.Error(err)
			log.Error("重命名本地文件失败:"+d.BackupFileName, " to ", newFileName)
			return
		}
	}
	os.MkdirAll(filepath.Dir(d.BackupFileName), os.ModePerm)
	f, err := os.Create(d.BackupFileName)
	if err != nil {
		log.Error(err)
		log.Error("备份到本地文件失败:" + d.BackupFileName)
		return
	}
	defer f.Close()
	for _, k := range d.ccs.Keys() {
		v := d.ccs.GetDefault(k, "")
		f.WriteString(k)
		f.WriteString("=")
		f.WriteString(v)
		f.WriteString("\n")
		err = f.Sync()
		if err != nil {
			log.Error(err)
		}
	}

}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		log.Warn(err)
		return false
	}
	return false
}
