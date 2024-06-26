package filex

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// FileLineCounter 文件统计行数
func FileLineCounter(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return LineCounter(file)
}

func LineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], []byte("\n"))
		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

// ReadFilePagination 分页读取文件
func ReadFilePagination(filePath string, page, limit int) (res []string, err error) {
	if page <= 0 {
		page = 1
	}
	lineTotal, err := FileLineCounter(filePath)
	if err != nil {
		return
	}
	// 读取全部内容并且返回数组
	begin := lineTotal - page*limit
	end := begin + limit
	if begin < 0 {
		begin = 0
	}
	if end < 0 {
		end = 0
	}
	return ReadFileByLine(filePath, begin, end)
}

// DirSize 获取目录大小
func DirSize(path string) (uint64, error) {
	var size uint64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info == nil {
			return errors.New("dir does not exist")
		}
		if !info.IsDir() {
			size += uint64(info.Size())
		}
		return err
	})
	return size, err
}

func ReadFileByLine(filePath string, begin, end int) (res []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	var count int
	for {
		lineData, err := rd.ReadString('\n')
		if err != nil || io.EOF == err {
			return res, err
		}
		count++
		if count > begin {
			res = append(res, lineData)
		}
		if count == end {
			return res, err
		}
	}
}
