package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

// 通过反射复制struct，只能复制int string
func Assign(origin, target interface{}) (err error) {
	val_origin := reflect.ValueOf(origin).Elem()
	val_target := reflect.ValueOf(target).Elem()

	for i := 0; i < val_origin.NumField(); i++ {

		field_name := val_origin.Type().Field(i).Name
		origin_field := val_origin.Field(i)

		if !val_target.FieldByName(field_name).IsValid() {
			continue
		}

		target_field := val_target.FieldByName(field_name)

		if origin_field.Kind() == target_field.Kind() {
			target_field.Set(origin_field)
			continue
		}

		switch origin_field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			target_field.SetInt(origin_field.Int())
		case reflect.String:
			if target_field.Kind() == reflect.String {
				target_field.SetString(origin_field.String())
			} else {
				continue
			}
		case reflect.Map, reflect.Slice:
			continue
		}
	}
	return
}

// 生成uuid
func GenerateUuid() string {
	rand.Seed(time.Now().UnixNano())
	str := fmt.Sprintf("%d%d", int(time.Now().Unix())-rand.Intn(10000), rand.Intn(1000000))
	return MD5(str)
}

// 生成skey
func GenerateSkey() string {
	rand.Seed(time.Now().UnixNano())
	str := fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(1000000))
	return MD5(str)
}

// md5加密
func MD5(str string) string {
	r := md5.Sum([]byte(str))
	return hex.EncodeToString(r[:])
}
