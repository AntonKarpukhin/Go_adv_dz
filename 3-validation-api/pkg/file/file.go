package file

import (
	"fmt"
	"os"
)

type JsFile struct {
	fileName string
}

func NewJsFile(fileName string) *JsFile {
	return &JsFile{
		fileName: fileName,
	}
}

func (JsFile *JsFile) Read() ([]byte, error) {
	data, err := os.ReadFile(JsFile.fileName)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	return data, nil
}

func (JsFile *JsFile) Write(data []byte) {
	file, err := os.Create(JsFile.fileName)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Println("Запись успешна")
}
