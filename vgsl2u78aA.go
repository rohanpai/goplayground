package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const MaxSerialNo = math.MaxUint64 - math.MaxUint32

func SerialNoInc(sn string) (string, error) {
	n, err := strconv.ParseUint(sn, 10, 64)
	if err == nil && n >= MaxSerialNo {
		err = errors.New("serial number invalid after increment")
	}
	if err != nil {
		return "", fmt.Errorf("SerialNoInc(%#v): %s", sn, err.Error())
	}
	s := strconv.FormatUint(n+1, 10)
	z := len(sn) - len(s)
	if z < 0 {
		z = 0
	}
	return sn[:z] + s, nil
}

func pubSerialStrInc(ss, sep string) (string, error) {
	i := strings.Index(ss, sep)
	if i < 0 {
		return "", fmt.Errorf("pubSerialStrInc(%#v, %#v): separator %#v not found", ss, sep, sep)
	}
	i += len(sep)
	if i >= len(ss) {
		return "", fmt.Errorf("pubSerialStrInc(%#v, %#v): serial number not found", ss, sep)
	}
	sn, err := SerialNoInc(ss[i:])
	if err != nil {
		return "", fmt.Errorf("pubSerialStrInc(%#v, %#v): %s", ss, sep, err.Error())
	}
	return ss[:i] + sn, nil
}

func main() {
	serialStrs := []struct {
		ss, sep string
	}{
		{"xs-00009", "-"},
		{"yxs<>000099", "<>"},
		{"yxzs-99", "-"},
		{"yxzs-" + strconv.FormatUint(MaxSerialNo-1, 10), "-"},
		{"yxzs-0000" + strconv.FormatUint(MaxSerialNo-1, 10), "-"},
		{"#01", "#"},
		{"01", ""},
		{"err-09", "="},
		{"err-", "-"},
		{"err-09x", "-"},
		{"err-" + strconv.FormatUint(MaxSerialNo, 10), "-"},
		{"err-" + strconv.FormatUint(math.MaxUint64, 10), "-"},
	}
	for _, serialStr := range serialStrs {
		si, err := pubSerialStrInc(serialStr.ss, serialStr.sep)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(serialStr.ss, "++  ==>", si)
	}
}
