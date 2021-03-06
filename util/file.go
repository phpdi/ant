package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//判定文件是否存在
func IsFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

//将数据刷入文件
func Data2File(filePath string, data interface{}) error {

	//文件不存在创建文件
	if !IsFileExist(filePath) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, content, 0644)

}

//将文件数据载入到变量中
func File2Data(filePath string, data interface{}) error {
	if !IsFileExist(filePath) {
		return nil
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, data)
}

////按照日期存储文件
//func SaveFileWithDate(dir string, fileHeader *multipart.FileHeader, randName bool) (filePath string /*上传的文件路径*/, err error) {
//
//	var (
//		fileNameArr []string
//		fileSuffix  string //文件后缀
//
//		newFile *os.File
//		file    *os.File
//		ok      bool
//	)
//
//	//取文件前缀
//	if fileHeader == nil {
//		return "", errors.New("文件不存在")
//	}
//
//	fileNameArr = strings.Split(fileHeader.Filename, ".")
//	if len(fileNameArr) < 1 {
//		return "", errors.New("文件名错误")
//	}
//
//	fileSuffix = "." + fileNameArr[len(fileNameArr)-1]
//
//	//计算目录
//	today := time.Now().Format("2006-01-02")
//
//	dir = strings.TrimRight(dir, "/") + "/" + today
//
//	if ok, err = PathExists(dir); err != nil {
//		return
//	}
//	//今日目录不存在创建
//	if !ok {
//		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
//			return
//		}
//	}
//
//	if randName {
//		filePath = dir + "/" + string(RandomCreateBytes(16)) + fileSuffix
//	} else {
//		filePath = dir + "/" + fileHeader.Filename
//	}
//
//	//检查文件是否存在
//	if ok, err = PathExists(filePath); err != nil {
//		return
//	}
//
//	if ok {
//		return "", errors.New("文件已存在")
//	}
//
//	//原文件
//	if file, err = os.Open(fileHeader.Filename); err != nil {
//		return err
//	}
//
//	defer file.Close()
//
//	//新文件
//	if newFile, err = os.Create(filePath); err != nil {
//		return "", err
//	}
//	defer newFile.Close()
//
//	if _, err = io.Copy(newFile, file); err != nil {
//		return "", err
//	}
//
//	return
//}

//计算文件路径
func FilePathWithDate(dir string, fileName string) (path string, err error) {
	var ok bool

	//计算目录
	today := time.Now().Format("2006-01-02")

	dir = strings.TrimRight(dir, "/") + "/" + today

	if ok, err = PathExists(dir); err != nil {
		return
	}
	//今日目录不存在创建
	if !ok {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return
		}
	}

	path = dir + "/" + fileName

	return
}

//检查文件是否符合要求
func CheckFile(fileHeader *multipart.FileHeader, size int64 /*支持的大小，单位kb*/, suffix []string /*支持的后缀*/) error {
	var (
		fileNameArr []string
		fileSuffix  string //文件后缀
	)

	if fileHeader == nil {
		return errors.New("文件不存在")
	}

	fileNameArr = strings.Split(fileHeader.Filename, ".")
	if len(fileNameArr) < 1 {
		return errors.New("文件名错误")
	}

	fileSuffix = "." + fileNameArr[len(fileNameArr)-1]

	//格式检查
	if len(suffix) > 0 {
		find := false
		for _, v := range suffix {
			if strings.ToLower(fileSuffix) == strings.ToLower(v) {
				find = true
				break
			}
		}

		if !find {
			return errors.New("不支持的文件格式")
		}
	}

	//大小检查
	if fileHeader.Size > size*1000 {
		return fmt.Errorf("文件过大:%d", fileHeader.Size)
	}

	return nil
}

//获取文件修改时间 返回unix时间戳
func GetFileModTime(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return 0, err
	}

	return fi.ModTime().Unix(), nil
}
func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)

	return hex.EncodeToString(_md5.Sum(nil))
}

//文件sha1
func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

//路径是否存在
func PathExists(path string) (bool, error) {
	var err error
	if _, err := os.Stat(path); err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

//获取文件大小
func GetFileSize(filename string) (bit int64, err error) {
	err = filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		bit = info.Size()
		return nil
	})

	return
}

func FileSha1(file io.Reader) (string, error) {
	var err error
	_sha1 := sha1.New()
	if _, err = io.Copy(_sha1, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(_sha1.Sum(nil)), nil
}

//根据日期转移文件
func MoveFileWithDate(filePath string, dir string, randName bool) (newPath string, err error) {
	var (
		fileSuffix string //文件后缀
		fileName   string //文件名
		file       *os.File
		newFile    *os.File
	)
	if arr := strings.Split(filePath, "/"); len(arr) > 0 {
		fileName = arr[len(arr)-1]
	} else {
		return "", errors.New("文件名解析错误")
	}

	if arr := strings.Split(fileName, "."); len(arr) > 0 {
		fileSuffix = "." + arr[len(arr)-1]
	} else {
		return "", errors.New("文件后缀名解析错误")
	}

	//计算目录
	today := time.Now().Format("2006-01-02")
	dir = strings.TrimRight(dir, "/") + "/" + today
	//今日目录不存在创建
	if !IsFileExist(dir) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return
		}
	}

	if randName {
		newPath = dir + "/" + string(RandomCreateBytes(16)) + fileSuffix
	} else {
		newPath = dir + "/" + fileName
	}

	//检查文件是否存在
	if IsFileExist(newPath) {
		return "", errors.New("文件已存在")
	}

	if file, err = os.Open(filePath); err != nil {
		return
	}
	defer file.Close()

	//新文件
	if newFile, err = os.Create(newPath); err != nil {
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)

	return

}

//扫描路径：dirPath
//扫描类型：scanType。0=全部，1=文件夹，2=文件
func ScanPath(dirPath string, scanType int) (fileList map[string]os.FileInfo, err error) {
	fileList = make(map[string]os.FileInfo)

	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}

		if scanType == 1 {
			if info.IsDir() {
				fileList[path] = info
			}
		} else if scanType == 2 {
			if !info.IsDir() {
				fileList[path] = info
			}
		} else {
			fileList[path] = info
		}

		return nil
	})

	return
}

// 扫描当前目录下文件，不递归扫描
//扫描类型：scanType。0=全部，1=文件夹，2=文件
func ScanDir(dirName string, scanType int) []string {

	dirNameSeparator := dirName
	if !strings.HasSuffix(dirNameSeparator, string(os.PathSeparator)) {
		dirNameSeparator += string(os.PathSeparator)
	}

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Println(err)
	}
	var fileList []string
	for _, file := range files {
		if scanType == 0 || (scanType == 1 && file.IsDir()) || (scanType == 2 && !file.IsDir()) {

			fileList = append(fileList, dirNameSeparator+file.Name())
		}
	}
	return fileList
}

//从go代码文件中解析出中文
func ParseChnFromGolang(filePath string) (words []string) {
	var tmp []string
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	re := regexp.MustCompile(`"(?U)(.)+"`)
	//re := regexp.MustCompile(`"\p{Han}[\p{Han}0-9a-zA-Z]+`)
	tmp = re.FindAllString(string(content), -1)

	rechn := regexp.MustCompile(`\p{Han}+`)

	for _, v := range tmp {
		if rechn.Match([]byte(v)) &&
			!strings.Contains(v, "=") &&
			!strings.Contains(v, "/") &&
			!strings.Contains(v, "->") {
			words = append(words, v)
		}
	}
	//
	//for _,v:=range words {
	//	fmt.Println(v)
	//}
	//
	return
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
