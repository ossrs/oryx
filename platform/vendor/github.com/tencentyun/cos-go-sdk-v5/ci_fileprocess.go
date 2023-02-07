package cos

import (
	"context"
	"encoding/xml"
	"net/http"
)

type FileHashCodeConfig struct {
	Type        string `xml:",omitempty"`
	AddToHeader bool   `xml:",omitempty"`
}

type FileHashCodeResult struct {
	MD5          string `xml:",omitempty"`
	SHA1         string `xml:",omitempty"`
	SHA256       string `xml:",omitempty"`
	FileSize     int    `xml:",omitempty"`
	LastModified string `xml:",omitempty"`
	Etag         string `xml:",omitempty"`
}

type FileUncompressConfig struct {
	Prefix         string `xml:",omitempty"`
	PrefixReplaced string `xml:",omitempty"`
}

type FileUncompressResult struct {
	Region    string `xml:",omitempty"`
	Bucket    string `xml:",omitempty"`
	FileCount string `xml:",omitempty"`
}

type FileCompressConfig struct {
	Flatten string `xml:",omitempty"`
	Format  string `xml:",omitempty"`
	UrlList string `xml:",omitempty"`
	Prefix  string `xml:",omitempty"`
	Key     string `xml:",omitempty"`
}

type FileCompressResult struct {
	Region string `xml:",omitempty"`
	Bucket string `xml:",omitempty"`
	Object string `xml:",omitempty"`
}

type FileProcessInput FileCompressResult
type FileProcessOutput FileCompressResult

type FileProcessJobOperation struct {
	FileHashCodeConfig   *FileHashCodeConfig   `xml:",omitempty"`
	FileHashCodeResult   *FileHashCodeResult   `xml:",omitempty"`
	FileUncompressConfig *FileUncompressConfig `xml:",omitempty"`
	FileUncompressResult *FileUncompressResult `xml:",omitempty"`
	FileCompressConfig   *FileCompressConfig   `xml:",omitempty"`
	FileCompressResult   *FileCompressResult   `xml:",omitempty"`
	Output               *FileProcessOutput    `xml:",omitempty"`
	UserData             string                `xml:",omitempty"`
}

type FileProcessJobOptions struct {
	XMLName          xml.Name                 `xml:"Request"`
	Tag              string                   `xml:",omitempty"`
	Input            *FileProcessInput        `xml:",omitempty"`
	Operation        *FileProcessJobOperation `xml:",omitempty"`
	QueueId          string                   `xml:",omitempty"`
	CallBackFormat   string                   `xml:",omitempty"`
	CallBackType     string                   `xml:",omitempty"`
	CallBack         string                   `xml:",omitempty"`
	CallBackMqConfig string                   `xml:",omitempty"`
}

type FileProcessJobResult struct {
	XMLName    xml.Name               `xml:"Response"`
	JobsDetail *FileProcessJobsDetail `xml:",omitempty"`
}

type FileProcessJobsDetail struct {
	Code         string                   `xml:",omitempty"`
	Message      string                   `xml:",omitempty"`
	JobId        string                   `xml:",omitempty"`
	Tag          string                   `xml:",omitempty"`
	State        string                   `xml:",omitempty"`
	CreationTime string                   `xml:",omitempty"`
	StartTime    string                   `xml:",omitempty"`
	EndTime      string                   `xml:",omitempty"`
	QueueId      string                   `xml:",omitempty"`
	Input        *FileProcessInput        `xml:",omitempty"`
	Operation    *FileProcessJobOperation `xml:",omitempty"`
}

// 提交哈希值计算任务 https://cloud.tencent.com/document/product/436/83108
// 提交文件解压任务 https://cloud.tencent.com/document/product/436/83110
// 提交多文件打包压缩任务 https://cloud.tencent.com/document/product/436/83112
func (s *CIService) CreateFileProcessJob(ctx context.Context, opt *FileProcessJobOptions) (*FileProcessJobResult, *Response, error) {
	var res FileProcessJobResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.CIURL,
		uri:     "/file_jobs",
		method:  http.MethodPost,
		body:    opt,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// 查询哈希值计算结果 https://cloud.tencent.com/document/product/436/83109
// 查询文件解压结果 https://cloud.tencent.com/document/product/436/83111
// 查询多文件打包压缩结果 https://cloud.tencent.com/document/product/436/83113
func (s *CIService) DescribeFileProcessJob(ctx context.Context, jobid string) (*FileProcessJobResult, *Response, error) {
	var res FileProcessJobResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.CIURL,
		uri:     "/file_jobs/" + jobid,
		method:  http.MethodGet,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// GetFileHashOptions is the option of GetFileHash
type GetFileHashOptions struct {
	CIProcess   string `url:"ci-process,omitempty"`
	Type        string `url:"type,omitempty"`
	AddToHeader bool   `url:"addtoheader,omitempty"`
}

// GetFileHashResult is the result of GetFileHash
type GetFileHashResult struct {
	XMLName            xml.Name            `xml:"Response"`
	FileHashCodeResult *FileHashCodeResult `xml:",omitempty"`
	Input              *FileProcessInput   `xml:",omitempty"`
}

// 哈希值计算同步请求 https://cloud.tencent.com/document/product/436/83107
func (s *CIService) GetFileHash(ctx context.Context, name string, opt *GetFileHashOptions) (*GetFileHashResult, *Response, error) {
	var res GetFileHashResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.BucketURL,
		uri:      "/" + encodeURIComponent(name),
		method:   http.MethodGet,
		optQuery: opt,
		result:   &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}
