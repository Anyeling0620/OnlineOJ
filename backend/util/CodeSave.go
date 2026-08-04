package util

import "os"

func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.MkdirAll(dirName, 077)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	_, err = f.Write(code)
	if err != nil {
		return "", err
	}
	defer func() {
		err = f.Close()
	}()
	return path, nil
}
