package pbsdk

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cast"
	"gopkg.in/resty.v1"
)

type Client struct {
	EndPoint  string
	AuthToken string
}

func NewClient() *Client {
	pb := &Client{}
	pb.EndPoint = os.Getenv("pb_endpoint")
	pb.AuthToken = os.Getenv("pb_authkey")
	if len(pb.EndPoint) == 0 || len(pb.AuthToken) == 0 {
		return nil
	}
	return pb
}

// FetchList 查询多条记录
func (r *Client) FetchList(collection string,
	page, perPage int,
	sorts []string,
	fields []string,
	filters map[string]string) ([]byte, *Response) {

	// 配置URL请求
	target := fmt.Sprintf(`%s/api/collections/%s/records`, r.EndPoint, collection)

	client := resty.New()
	// 配置鉴权Token
	request := client.R().SetHeader("Content-Type", "application/json")
	if len(r.AuthToken) > 0 {
		request = request.SetHeader("Authorization", r.AuthToken)
	}
	// 配置分页
	if page > 0 {
		request.SetQueryParam("page", cast.ToString(page))
	}
	if perPage > 0 {
		request.SetQueryParam("perPage", cast.ToString(perPage))
	}
	// 配置排序
	if len(sorts) > 0 {
		request.SetQueryParam("sort", strings.Join(sorts, ","))
	}
	// 配置返回字段
	if len(fields) > 0 {
		request.SetQueryParam("fields", strings.Join(fields, ","))
	}
	// 配置过滤器
	if len(filters) > 0 {
		items := []string{}
		for k, v := range filters {
			items = append(items, fmt.Sprintf(`%s=%s`, k, v))
		}
		filter := strings.Join(items, " && ")
		request.SetQueryParam("filter", filter)
	}

	resp, err := request.Get(target)
	if err != nil {
		return nil, ResponseFromError(err)
	}
	return resp.Body(), ResponseFromResty(resp)
}

// FetchOne 查询一条记录
func (r *Client) FetchOne(collection, id string, fields string) ([]byte, *Response) {
	target := fmt.Sprintf(`%s/api/collections/%s/records/%s`, r.EndPoint, collection, id)
	client := resty.New()
	request := client.R().SetHeader("Content-Type", "application/json")
	if len(r.AuthToken) > 0 {
		request = request.SetHeader("Authorization", r.AuthToken)
	}

	if len(fields) > 0 {
		request.SetQueryParam("fields", fields)
	}

	resp, err := request.Get(target)
	if err != nil {
		return nil, ResponseFromError(err)
	}
	return resp.Body(), ResponseFromResty(resp)
}

// CreateOne 创建一条记录
func (r *Client) CreateOne(collection string, body interface{}) ([]byte, *Response) {
	target := fmt.Sprintf(`%s/api/collections/%s/records`, r.EndPoint, collection)
	client := resty.New()
	request := client.R().SetHeader("Content-Type", "application/json")
	if len(r.AuthToken) > 0 {
		request = request.SetHeader("Authorization", r.AuthToken)
	}
	resp, err := request.SetBody(body).Post(target)
	if err != nil {
		return nil, ResponseFromError(err)
	}
	return resp.Body(), ResponseFromResty(resp)
}

// UpdateOne 更新一条记录
func (r *Client) UpdateOne(collection, id string, body interface{}) ([]byte, *Response) {
	target := fmt.Sprintf(`%s/api/collections/%s/records/%s`, r.EndPoint, collection, id)
	client := resty.New()
	request := client.R().SetHeader("Content-Type", "application/json")
	if len(r.AuthToken) > 0 {
		request = request.SetHeader("Authorization", r.AuthToken)
	}
	resp, err := request.SetBody(body).Patch(target)
	if err != nil {
		return nil, ResponseFromError(err)
	}
	return resp.Body(), ResponseFromResty(resp)
}

// DeleteOne 删除一条记录
func (r *Client) DeleteOne(collection string, id string) *Response {
	target := fmt.Sprintf(`%s/api/collections/%s/records/%s`, r.EndPoint, collection, id)
	client := resty.New()
	request := client.R().SetHeader("Content-Type", "application/json")
	if len(r.AuthToken) > 0 {
		request = request.SetHeader("Authorization", r.AuthToken)
	}
	resp, err := request.Delete(target)
	if err != nil {
		return ResponseFromError(err)
	}
	return ResponseFromResty(resp)
}
