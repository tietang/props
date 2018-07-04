package utils

import (
    "bytes"
    "strconv"
)

func Panic(err error) {
    if err != nil {
        panic(err)
    }
}

//byte转16进制字符串
func ByteToHex(data []byte) string {
    buffer := new(bytes.Buffer)
    for _, b := range data {

        s := strconv.FormatInt(int64(b&0xff), 16)
        if len(s) == 1 {
            buffer.WriteString("0")
        }
        buffer.WriteString(s)
    }

    return buffer.String()
}

//16进制字符串转[]byte
func HexToBye(hex string) []byte {
    length := len(hex) / 2
    slice := make([]byte, length)
    rs := []rune(hex)

    for i := 0; i < length; i++ {
        s := string(rs[i*2: i*2+2])
        value, _ := strconv.ParseInt(s, 16, 10)
        slice[i] = byte(value & 0xFF)
    }
    return slice
}
