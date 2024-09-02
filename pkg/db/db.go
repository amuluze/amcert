// Package db
// Date       : 2024/9/2 10:35
// Author     : Amu
// Description:
package db

import "path/filepath"

func Init(path string) error {
	storagePath := filepath.Join(path, "storage.db")
	err := Default.SetPath(storagePath)
	if err != nil {
		return err
	}
	return nil
}

func SaveKey(k, v string) error {
	return Default.Put(k, []byte(v))
}

func GetValue(key string) (string, error) {
	return Default.GetString(key)
}

func GetPrefixKeys(key string) ([]string, error) {
	return Default.Keys(key, false)
}

func DeleteKey(key string) error {
	return Default.Delete(key)
}

func PutJson(key string, value interface{}) error {
	return Default.PutJson(key, value)
}

func GetJson(key string) (interface{}, error) {
	return Default.GetString(key)
}
