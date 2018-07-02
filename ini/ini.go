package ini

import (
    "io"
    "os"
    log "github.com/sirupsen/logrus"
    "github.com/go-ini/ini"
    "strings"
    "github.com/tietang/props/kvs"
)

type IniProperties struct {
    kvs.MapProperties
    //values map[string]string
    IniFile *ini.File
}

func NewIniProperties() *IniProperties {
    p := &IniProperties{
        //values: make(map[string]string),
    }
    p.Values = make(map[string]string)
    return p
}

// Read creates a new property set and fills it with the contents of a file.
// See Load for the supported file format.
func ReadIni(r io.Reader) (*IniProperties, error) {
    p := NewIniProperties()
    err := p.Load(r)
    if err != nil {
        log.Error(err)
        return nil, err
    }
    return p, nil
}

func ReadIniFile(f string) (*IniProperties, error) {

    file, err := os.Open(f)
    defer file.Close()

    if err != nil {
        d, _ := os.Getwd()
        log.WithField("error", err.Error()).Fatal("read file: ", d, "  ", f)
        return nil, err
    }
    return ReadIni(file)
}

func (p *IniProperties) Load(r io.Reader) error {

    props, err := ini.Load(r)
    p.IniFile = props
    if err != nil {
        log.Warn(err)
        return err
    }
    sections := props.Sections()
    for _, section := range sections {
        name := strings.TrimSpace(section.Name())
        for _, kv := range section.Keys() {
            kvName := strings.TrimSpace(kv.Name())
            key := strings.Join([]string{name, kvName}, ".")
            value := strings.TrimSpace(kv.String())
            p.Values[key] = value
        }
    }

    return nil
}

// Write saves the property set to a file. The output will be in "key=value"
// format, with appropriate characters escaped. See Load for more details on
// the file format.
//
// Note: if the property set was loaded from a file, the formatting and
// comments from the original file will not be retained in the output file.
func (p *IniProperties) Write(w io.Writer) error {
    _, err := p.IniFile.WriteTo(w)
    return err
}
