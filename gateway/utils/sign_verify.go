/***************************************************
 ** @Desc : This file for ...
 ** @Time : 2019/10/26 11:08
 ** @Author : yuebin
 ** @File : sign_verify
 ** @Last Modified by : yuebin
 ** @Last Modified time: 2019/10/26 11:08
 ** @Software: GoLand
****************************************************/
package utils

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"time"
)

func GetMD5Sign(params map[string]string, keys []string, paySecret string) string {
	str := ""
	for i := 0; i < len(keys); i++ {
		k := keys[i]
		if len(params[k]) == 0 {
			continue
		}
		str += k + "=" + params[k] + "&"
	}
	str += "paySecret=" + paySecret
	sign := GetMD5Upper(str)
	return sign
}

/*
* 验签
 */
func Md5Verify(params map[string]string, paySecret string) bool {
	sign := params["sign"]
	if sign == "" {
		return false
	}

	delete(params, "sign")
	keys := SortMap(params)
	tmpSign := GetMD5Sign(params, keys, paySecret)
	if tmpSign != sign {
		return false
	} else {
		return true
	}
}

/*
* Md5VerifyWithTimestamp 增强版验签（包含时间戳验证，防重放攻击）
* params: 请求参数
* paySecret: 支付密钥
* timestampKey: 时间戳参数名（默认为 "timestamp"）
* timeout: 允许的时间偏差（秒），默认300秒（5分钟）
 */
func Md5VerifyWithTimestamp(params map[string]string, paySecret string, timestampKey string, timeout int64) (bool, string) {
	// 1. 先验证签名
	sign := params["sign"]
	if sign == "" {
		return false, "签名不能为空"
	}

	// 2. 验证时间戳是否存在
	if timestampKey == "" {
		timestampKey = "timestamp"
	}

	timestampStr := params[timestampKey]
	if timestampStr == "" {
		return false, fmt.Sprintf("时间戳参数 %s 不能为空", timestampKey)
	}

	// 3. 解析时间戳
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		logs.Error("解析时间戳失败: %s, error: %v", timestampStr, err)
		return false, "时间戳格式错误"
	}

	// 4. 验证时间戳是否在有效期内（防重放攻击）
	if timeout <= 0 {
		timeout = 300 // 默认5分钟
	}

	currentTime := time.Now().Unix()
	timeDiff := currentTime - timestamp

	if timeDiff < 0 {
		// 请求时间在未来，可能是时钟不同步或恶意请求
		logs.Warn("请求时间在未来，当前时间: %d, 请求时间: %d, 差值: %d", currentTime, timestamp, timeDiff)
		return false, "请求时间异常（时间在未来）"
	}

	if timeDiff > timeout {
		// 请求超时，可能是重放攻击
		logs.Warn("请求时间戳超时，当前时间: %d, 请求时间: %d, 差值: %d秒, 允许: %d秒",
			currentTime, timestamp, timeDiff, timeout)
		return false, fmt.Sprintf("请求超时，时间差: %d秒", timeDiff)
	}

	// 5. 验证MD5签名
	tmpParams := make(map[string]string)
	for k, v := range params {
		tmpParams[k] = v
	}
	delete(tmpParams, "sign")

	keys := SortMap(tmpParams)
	tmpSign := GetMD5Sign(tmpParams, keys, paySecret)

	if tmpSign != sign {
		logs.Warn("签名验证失败，期望: %s, 实际: %s", tmpSign, sign)
		return false, "签名验证失败"
	}

	return true, "验证成功"
}

/*
* VerifySignAndTimestamp 便捷方法，使用默认参数验证签名和时间戳
* 时间戳参数名: timestamp
* 超时时间: 300秒（5分钟）
 */
func VerifySignAndTimestamp(params map[string]string, paySecret string) (bool, string) {
	return Md5VerifyWithTimestamp(params, paySecret, "timestamp", 300)
}

/*
* GenerateSignWithTimestamp 生成包含时间戳的签名
* 用于生成请求签名时自动添加时间戳
 */
func GenerateSignWithTimestamp(params map[string]string, paySecret string) (sign string, timestamp int64) {
	// 1. 添加当前时间戳
	timestamp = time.Now().Unix()
	params["timestamp"] = strconv.FormatInt(timestamp, 10)

	// 2. 生成签名
	keys := SortMap(params)
	sign = GetMD5Sign(params, keys, paySecret)

	return sign, timestamp
}
