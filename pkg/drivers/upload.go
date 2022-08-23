package drivers

import (
	"FileTransfer/pkg/utils/unique"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/unknwon/com"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 创建文件的默认权限，比如 Upload.dir 若不存在，会使用此权限创建目录。
const defaultMode = 0777

// 常用错误类型
var (
	ErrNotAllowExt  = errors.New("不允许的文件上传类型")
	ErrNotAllowSize = errors.New("文件上传大小超过最大设定值")
	ErrNoUploadFile = errors.New("客户端没有上传文件")
)

// Upload 用于处理文件上传
type Upload struct {
	fs        fs.FS
	dir       string
	format    string
	maxSize   int64
	notExtS   []string
	filenames func(string) string
}

// UniqueFilename 生成唯一文件名
func UniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	return unique.String().String() + ext
}

// New 声明一个 Upload 对象
//
// dir 上传文件的保存目录，若目录不存在，则会尝试创建；
// format 子目录的格式，只能是时间格式；
// maxSize 允许上传文件的最大尺寸，单位为 byte；
// exts 允许的扩展名，若为空，将不允许任何文件上传。
func New(dir, format string, maxSize int64, notExtS ...string) (*Upload, error) {
	// 确保所有的后缀名都是以.作为开始符号的。
	es := make([]string, 0, len(notExtS))
	for _, ext := range notExtS {
		if ext[0] != '.' {
			ext = "." + ext
		}
		es = append(es, strings.ToLower(ext))
	}

	// 确保 dir 最后一个字符为目录分隔符。
	last := dir[len(dir)-1]
	if last != '/' && last != filepath.Separator {
		dir = dir + string(filepath.Separator)
	}

	// 若不存在目录，则尝试创建
	if err := os.MkdirAll(dir, defaultMode); err != nil {
		return nil, err
	}

	// 确保 dir 目录存在。
	// NOTE:此处的 dir 最后个字符为/，所以不用判断是否为目录。
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	return &Upload{
		fs:        os.DirFS(dir),
		dir:       dir,
		format:    format,
		maxSize:   maxSize,
		notExtS:   es,
		filenames: UniqueFilename,
	}, nil
}

func (u *Upload) Open(name string) (fs.File, error) { return u.fs.Open(name) }

// 判断扩展名是否符合要求
func (u *Upload) isNotAllowExt(ext string) bool {
	if len(ext) == 0 { // 没有扩展名，一律通过
		return true
	}

	for _, e := range u.notExtS {
		if e == ext {
			return false
		}
	}
	return true
}

func (u *Upload) getDestPath(filename string) string {
	return time.Now().Format(u.format) + u.filenames(filename)
}

// Dir 获取上传文件的保存目录
func (u *Upload) Dir() string { return u.dir }

// Do 执行上传的操作
//
// 若是多文件上传，其中某一个文件不符合要求，会中断后续操作，
// 但是已经处理成功的也会返回给用户，所以可能会出现两个返回参数都不为 nil 的情况。
//
// 返回的是相对于 Upload.Dir() 目录的文件名列表。
func (u *Upload) Do(field string, r *http.Request) ([]string, error) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, ErrNoUploadFile
	}

	heads := r.MultipartForm.File[field]
	if len(heads) == 0 {
		return nil, ErrNoUploadFile
	}

	ret := make([]string, 0, len(heads))

	for _, head := range heads {
		path, err := u.moveFile(head)
		if err != nil {
			return ret, err // 已经上传的内容 ret，需要正常返回给用户
		}

		ret = append(ret, path)
	}

	return ret, nil
}

// 将上传的文件移到 u.Dir 目录下
//
// 返回相对于 u.Dir 的地址
func (u *Upload) moveFile(head *multipart.FileHeader) (string, error) {
	if head.Size > u.maxSize {
		return "", ErrNotAllowSize
	}

	ext := strings.ToLower(filepath.Ext(head.Filename))
	if !u.isNotAllowExt(ext) {
		return "", ErrNotAllowExt
	}

	file, err := head.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	path := u.getDestPath(head.Filename)
	ret := path

	path = u.dir + path
	if err = os.MkdirAll(filepath.Dir(path), defaultMode); err != nil { // 若路径不存在，则创建
		return "", err
	}

	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err = io.Copy(f, file); err != nil {
		return "", err
	}

	return ret, nil
}
func DeleteFile(path string) bool {
	//删除上传文件
	if com.IsExist(path) {
		err := os.Remove(path)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

// GetFileMd5 获取文件的md5码
func GetFileMd5(path string) (string, error) {

	pFile, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("打开文件失败，filename=%v, err=%v", path, err)
	}
	defer pFile.Close()
	md5h := md5.New()
	io.Copy(md5h, pFile)

	return hex.EncodeToString(md5h.Sum(nil)), nil

}

// SetFilename 设置文件名的生成方式
func (u *Upload) SetFilename(f func(filename string) string) { u.filenames = f }
