// Package db
// Date       : 2024/9/2 10:35
// Author     : Amu
// Description:
package db

func Initialize(storagePath string) error {
	// TODO: 初始化数据库, 如果 storagePath 不存在则创建
	err := Default.SetPath(storagePath)
	if err != nil {
		return err
	}
	return nil
}

func DeleteKey(key string) error {
	return Default.Delete(key)
}

func GetPrefixKeys(key string) ([]string, error) {
	return Default.Keys(key, false)
}

func PutString(k, v string) error {
	return Default.Put(k, []byte(v))
}

func GetString(key string) (string, error) {
	return Default.GetString(key)
}

func PutJson(key string, value interface{}) error {
	return Default.PutJson(key, value)
}

func GetJson(key string, out interface{}) error {
	return Default.GetJson(key, out)
}
