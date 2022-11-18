package kvs

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasttemplate"
	"io"
	"regexp"
	"strings"
)

type EvalValue interface {
	EvalValue(value string) (string, error)
}
type Eval interface {
	EvalValue
	EvalAll()
}

const (
	__START_TAG   = "${"
	__END_TAG     = "}"
	DEFAULT_VALUE = ""
)

var __reg = regexp.MustCompile("\\$\\{(.*)}")
var _ Eval = new(DefaultEval)

func NewEval(conf ConfigSource) Eval {
	return &DefaultEval{
		conf:     conf,
		StartTag: __START_TAG,
		EndTag:   __END_TAG,
	}
}

type DefaultEval struct {
	StartTag string
	EndTag   string
	conf     ConfigSource
}

func (e *DefaultEval) EvalAll() {
	for _, key := range e.conf.Keys() {
		val, err := e.evalKey(key)
		if err == nil {
			e.conf.Set(key, val)
		} else {
			log.Warn("eval key: ", key, err)
		}
	}
}
func (e *DefaultEval) evalKey(key string) (val string, err error) {
	val, err = e.conf.Get(key)
	if err != nil {
		return val, err
	}
	return e.EvalValue(val)
}

func (e *DefaultEval) EvalValue(val string) (string, error) {
	if __reg.MatchString(val) {
		if strings.Contains(val, e.StartTag) {
			eval := fasttemplate.New(val, e.StartTag, e.EndTag)
			str := eval.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
				s, err := e.conf.Get(tag)
				if err == nil {
					return w.Write([]byte(s))
				} else {
					return w.Write([]byte(""))
				}
			})
			var err error
			if str != val {
				str, err = e.EvalValue(str)
				//fmt.Println(str, err)
			}
			return str, err
		}
	}
	return val, nil
}

func (e *DefaultEval) calculateEvalValue(value string) (string, error) {

	sub := __reg.FindStringSubmatch(value)
	if len(sub) == 0 {
		return value, nil
	}
	defaultValue := ""
	for _, k := range sub {

		keys := strings.Split(k, ":")
		if len(keys) > 1 {
			k = keys[0]
			defaultValue = keys[1]
		}
		v, err := e.conf.Get(k)
		if err == nil {
			return v, nil
		}
	}

	return defaultValue, errors.New("not exists")
}
