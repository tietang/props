package main

import (
	"fmt"
	"github.com/tietang/props"
	"time"
)

type ServerProperties struct {
	_prefix string `prefix:"http.server"`
	Port    int
	Timeout int `val:"1"`
	Enabled bool
	Foo     int `val:"1"`
	Time    time.Duration `val:"1s"`
	Float   float32 `val:"0.000001"`
}

func main() {
	//t := ServerProperties{}
	//s := reflect.ValueOf(&t).Elem()
	//typeOfT := s.Type()
	//
	//for i := 0; i < s.NumField(); i++ {
	//	f := s.Field(i)
	//	fmt.Printf("%d: %s %s = %v\n", i,
	//		typeOfT.Field(i).Name, f.Type(), f.Interface())
	//}
	//
	//s.Field(0).SetInt(25)
	//s.Field(1).SetString("nicky")
	//fmt.Println(t)

	//
	p := props.NewMapProperties()
	p.Set("http.server.port", "8080")
	//p.Set("http.server.timeout", "1234")
	//p.Set("http.server.enabled", "true")
	//p.Set("http.server.time", "10s")
	//p.Set("http.server.float", "23.45")
	s := &ServerProperties{}
	p.Unmarshal(s)
	fmt.Println(s)
}

//}
//
////用map填充结构
//func FillStruct(data map[string]interface{}, obj interface{}) error {
//	for k, v := range data {
//		err := SetField(obj, k, v)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
////用map的值替换结构的值
//func SetField(obj interface{}, name string, value interface{}) error {
//	structValue := reflect.ValueOf(obj).Elem()        //结构体属性值
//	structFieldValue := structValue.FieldByName(name) //结构体单个属性值
//	if !structFieldValue.IsValid() {
//		return fmt.Errorf("No such field: %s in obj", name)
//	}
//	if !structFieldValue.CanSet() {
//		return fmt.Errorf("Cannot set %s field value", name)
//	}
//	structFieldType := structFieldValue.Type() //结构体的类型
//	val := reflect.ValueOf(value)              //map值的反射值
//	var err error
//	if structFieldType != val.Type() {
//		val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name()) //类型转换
//		if err != nil {
//			return err
//		}
//	}
//	structFieldValue.Set(val)
//	return nil
//}
//
////类型转换
//func TypeConversion(value string, ntype string) (reflect.Value, error) {
//	if ntype == "string" {
//		return reflect.ValueOf(value), nil
//	} else if ntype == "time.Time" {
//		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
//		return reflect.ValueOf(t), err
//	} else if ntype == "Time" {
//		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
//		return reflect.ValueOf(t), err
//	} else if ntype == "int" {
//		i, err := strconv.Atoi(value)
//		return reflect.ValueOf(i), err
//	} else if ntype == "int8" {
//		i, err := strconv.ParseInt(value, 10, 64)
//		return reflect.ValueOf(int8(i)), err
//	} else if ntype == "int32" {
//		i, err := strconv.ParseInt(value, 10, 64)
//		return reflect.ValueOf(int64(i)), err
//	} else if ntype == "int64" {
//		i, err := strconv.ParseInt(value, 10, 64)
//		return reflect.ValueOf(i), err
//	} else if ntype == "float32" {
//		i, err := strconv.ParseFloat(value, 64)
//		return reflect.ValueOf(float32(i)), err
//	} else if ntype == "float64" {
//		i, err := strconv.ParseFloat(value, 64)
//		return reflect.ValueOf(i), err
//	}
//	//else if .......增加其他一些类型的转换
//	return reflect.ValueOf(value), errors.New("未知的类型："
//	ntype)
//}
