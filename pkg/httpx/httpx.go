package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// GetRequest 封装了HTTP GET请求的逻辑。
func GetRequest(ctx context.Context, client *http.Client, urlStr string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return resp, nil
}

// PostRequest 封装了HTTP POST请求的逻辑。
func PostRequest(ctx context.Context, client *http.Client, urlStr string, body io.Reader, headers map[string]string) (
	*http.Response, error) {
	// 创建一个新的POST请求
	req, err := http.NewRequestWithContext(ctx, "POST", urlStr, body)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 检查HTTP响应状态码
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close() // 确保在错误情况下也关闭响应体
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return resp, nil
}
