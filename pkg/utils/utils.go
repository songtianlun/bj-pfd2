package utils

import (
    "bytes"
    "crypto/sha1"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

func StrToUInt64(s string) (i uint64) {
    i, err := strconv.ParseUint(s, 10, 64)
    if err != nil {
        i = 0
    }
    return
}

func StrToFloat64(s string) (i float64) {
    i, err := strconv.ParseFloat(s, 64)
    if err != nil {
        i = 0
    }
    return
}

func Float64ToStr(f float64) string {
    return strconv.FormatFloat(f, 'f', -1, 64)
}

func ByteToString(b []byte) string {
    return string(b)
}

func ByteToMap(b []byte) map[string]interface{} {
    var m map[string]interface{}
    err := json.Unmarshal(b, &m)
    if err != nil {
        return nil
    }
    return m
}

func InterfaceToSlice(i interface{}) []interface{} {
    var s []interface{}
    s = append(s, i)
    return s
}

func P(a ...interface{}) {
    fmt.Println(a...)
}

func PrettyPrint(v interface{}) {
    bs, _ := json.Marshal(v)
    var out bytes.Buffer
    json.Indent(&out, bs, "", "\t")
    fmt.Printf("%v\n", out.String())
}

func PrettyJsonString(s string) string {
    var out bytes.Buffer
    json.Indent(&out, []byte(s), "", "\t")
    return out.String()
}

func IntToString(i int) string {
    return strconv.Itoa(i)
}

func Int32ToString(i int32) string {
    return strconv.Itoa(int(i))
}

func Int64ToString(i int64) string {
    return strconv.FormatInt(i, 10)
}

func Float64ToString(f float64) string {
    return strconv.FormatFloat(f, 'f', -1, 64)
}

// ErrorMessage Convenience function to redirect to the error message page
func ErrorMessage(writer http.ResponseWriter, request *http.Request, msg string) {
    url := []string{"/err?msg=", msg}
    http.Redirect(writer, request, strings.Join(url, ""), 302)
}

func EnDateWithYM(year int64, month int64) (code string) {
    if year <= 1000 || year >= 10000 {
        year = int64(time.Now().Year())
    }
    if month <= 0 || month > 12 {
        month = int64(time.Now().Month())
    }

    code += Int64ToString(year)
    code += "-"
    if month < 10 {
        code += "0"
    }

    code += Int64ToString(month)
    return
}

func EnDateWithYMD(year int64, month int64, day int64) (code string) {
    if year <= 1000 || year >= 10000 {
        year = int64(time.Now().Year())
    }
    if month <= 0 || month > 12 {
        month = int64(time.Now().Month())
    }
    if day <= 0 || day > 31 {
        day = int64(time.Now().Day())
    }

    code += Int64ToString(year)
    code += "-"
    if month < 10 {
        code += "0"
    }
    code += Int64ToString(month)
    code += "-"
    if day < 10 {
        code += "0"
    }
    code += Int64ToString(day)

    return
}

func PrintRMB(rmb float64) string {
    return fmt.Sprintf("%.2f ¥", rmb)
}

func Float64ToIntStrRMB(f float64) string {
    return fmt.Sprintf("%.0f", f)
}

func Encrypt(plaintext string) (cryptext string) {
    cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
    return
}

func GetTypeString(i interface{}) string {
    return fmt.Sprintf("%T", i)
}
