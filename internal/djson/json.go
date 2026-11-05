package djson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var 错误JSON无效 = errors.New("JSON无效")

func 解析(文本 string) (any, error) {
	解码器 := json.NewDecoder(strings.NewReader(文本))
	解码器.UseNumber()
	var 根 any
	if 错误 := 解码器.Decode(&根); 错误 != nil {
		return nil, 错误
	}
	return 根, nil
}

func 拆分路径(路径 string) []string {
	路径 = strings.TrimSpace(路径)
	if 路径 == "" {
		return nil
	}
	段列表 := strings.Split(路径, ".")
	结果 := make([]string, 0, len(段列表))
	for _, 段 := range 段列表 {
		段 = strings.TrimSpace(段)
		if 段 == "" {
			continue
		}
		结果 = append(结果, 段)
	}
	return 结果
}

func 取值(根 any, 路径 string) (any, bool) {
	段列表 := 拆分路径(路径)
	当前 := 根
	for _, 段 := range 段列表 {
		切换, 成功 := 当前.(map[string]any)
		if 成功 {
			下一级, 存在 := 切换[段]
			if !存在 {
				return nil, false
			}
			当前 = 下一级
			continue
		}
		数组, 成功 := 当前.([]any)
		if !成功 {
			return nil, false
		}
		下标, 错误 := strconv.Atoi(段)
		if 错误 != nil || 下标 < 0 || 下标 >= len(数组) {
			return nil, false
		}
		当前 = 数组[下标]
	}
	return 当前, true
}

func 紧凑编码(值 any) (string, error) {
	数据, 错误 := json.Marshal(值)
	if 错误 != nil {
		return "", 错误
	}
	return string(数据), nil
}

// D是否有效 判断文本是否为合法 JSON。
func 是否有效(文本 string) bool {
	_, 错误 := 解析(文本)
	return 错误 == nil
}

// D是否存在 判断路径是否存在。
func 是否存在(文本, 路径 string) bool {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		return false
	}
	_, 存在 := 取值(根, 路径)
	return 存在
}

// D取成员数 返回对象字段数或数组长度；路径空则对根取值；无效返回 0。
func 取成员数(文本, 路径 string) int {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		return 0
	}
	值, 存在 := 取值(根, 路径)
	if !存在 {
		return 0
	}
	switch 对象 := 值.(type) {
	case map[string]any:
		return len(对象)
	case []any:
		return len(对象)
	default:
		return 0
	}
}

// D取文本 取路径对应值并转为文本；找不到返回空串。
func 取文本(文本, 路径 string) string {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		return ""
	}
	值, 存在 := 取值(根, 路径)
	if !存在 || 值 == nil {
		return ""
	}
	switch v := 值.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		编码, 错误 := 紧凑编码(v)
		if 错误 != nil {
			return ""
		}
		return 编码
	}
}

// D取整数 取路径对应整数；找不到或非数字返回 0。
func 取整数(文本, 路径 string) int64 {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		return 0
	}
	值, 存在 := 取值(根, 路径)
	if !存在 || 值 == nil {
		return 0
	}
	switch v := 值.(type) {
	case json.Number:
		i, 错误 := v.Int64()
		if 错误 == nil {
			return i
		}
		f, 错误 := v.Float64()
		if 错误 != nil {
			return 0
		}
		return int64(f)
	case float64:
		return int64(v)
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		i, 错误 := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if 错误 != nil {
			f, 错误2 := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if 错误2 != nil {
				return 0
			}
			return int64(f)
		}
		return i
	default:
		return 0
	}
}

// D取小数 取路径对应浮点数；找不到或非数字返回 0。
func 取小数(文本, 路径 string) float64 {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		return 0
	}
	值, 存在 := 取值(根, 路径)
	if !存在 || 值 == nil {
		return 0
	}
	switch v := 值.(type) {
	case json.Number:
		f, 错误 := v.Float64()
		if 错误 != nil {
			return 0
		}
		return f
	case float64:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		f, 错误 := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if 错误 != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

// D取逻辑 取路径对应布尔；找不到返回 false。数字非 0 视为 true；文本 "true"/"1" 为 true。
func 取逻辑(文本, 路径 string) bool {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		return false
	}
	值, 存在 := 取值(根, 路径)
	if !存在 || 值 == nil {
		return false
	}
	switch v := 值.(type) {
	case bool:
		return v
	case json.Number:
		f, 错误 := v.Float64()
		if 错误 != nil {
			return false
		}
		return f != 0
	case float64:
		return v != 0
	case string:
		s := strings.TrimSpace(strings.ToLower(v))
		return s == "true" || s == "1"
	default:
		return false
	}
}

// D取JSON 取路径对应子树的紧凑 JSON；找不到返回空串。
func 取JSON(文本, 路径 string) string {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		return ""
	}
	值, 存在 := 取值(根, 路径)
	if !存在 {
		return ""
	}
	编码, 错误 := 紧凑编码(值)
	if 错误 != nil {
		return ""
	}
	return 编码
}

func 设置值(根 any, 路径 string, 新值 any) (any, error) {
	段列表 := 拆分路径(路径)
	if len(段列表) == 0 {
		return 新值, nil
	}
	if 根 == nil {
		根 = map[string]any{}
	}
	return 递归设置(根, 段列表, 新值)
}

func 递归设置(当前 any, 段列表 []string, 新值 any) (any, error) {
	段 := 段列表[0]
	最后 := len(段列表) == 1
	下标, 是数组下标 := 试解析下标(段)

	// 已是对象时，数字段也当对象键，避免把 {"0":...} 冲成数组。
	if _, 是对象 := 当前.(map[string]any); 是对象 && 是数组下标 {
		是数组下标 = false
	}

	if 是数组下标 {
		数组, 成功 := 当前.([]any)
		if !成功 {
			数组 = []any{}
		}
		for len(数组) <= 下标 {
			数组 = append(数组, nil)
		}
		if 最后 {
			数组[下标] = 新值
			return 数组, nil
		}
		下一级 := 数组[下标]
		下一级, 错误 := 递归设置(准备子节点(下一级, 段列表[1]), 段列表[1:], 新值)
		if 错误 != nil {
			return nil, 错误
		}
		数组[下标] = 下一级
		return 数组, nil
	}

	对象, 成功 := 当前.(map[string]any)
	if !成功 {
		对象 = map[string]any{}
	}
	if 最后 {
		对象[段] = 新值
		return 对象, nil
	}
	下一级 := 对象[段]
	下一级, 错误 := 递归设置(准备子节点(下一级, 段列表[1]), 段列表[1:], 新值)
	if 错误 != nil {
		return nil, 错误
	}
	对象[段] = 下一级
	return 对象, nil
}

func 试解析下标(段 string) (int, bool) {
	下标, 错误 := strconv.Atoi(段)
	if 错误 != nil || 下标 < 0 {
		return 0, false
	}
	return 下标, true
}

func 准备子节点(现有 any, 下一段 string) any {
	if 现有 != nil {
		if _, 是对象 := 现有.(map[string]any); 是对象 {
			return 现有
		}
		if _, 是数组 := 现有.([]any); 是数组 {
			return 现有
		}
	}
	if _, 是下标 := 试解析下标(下一段); 是下标 {
		return []any{}
	}
	return map[string]any{}
}

func 设置并编码(文本, 路径 string, 新值 any) (string, error) {
	根, 错误 := 解析(文本)
	if 错误 != nil {
		// 空文本视为空对象，便于从零构建
		if strings.TrimSpace(文本) == "" {
			根 = map[string]any{}
		} else {
			return "", 错误JSON无效
		}
	}
	新根, 错误 := 设置值(根, 路径, 新值)
	if 错误 != nil {
		return "", 错误
	}
	return 紧凑编码(新根)
}

// D设置文本 设置路径为字符串。
func 设置文本(文本, 路径, 值 string) (string, error) {
	return 设置并编码(文本, 路径, 值)
}

// D设置整数 设置路径为整数。
func 设置整数(文本, 路径 string, 值 int64) (string, error) {
	return 设置并编码(文本, 路径, 值)
}

// D设置小数 设置路径为小数。
func 设置小数(文本, 路径 string, 值 float64) (string, error) {
	return 设置并编码(文本, 路径, 值)
}

// D设置逻辑 设置路径为布尔。
func 设置逻辑(文本, 路径 string, 值 bool) (string, error) {
	return 设置并编码(文本, 路径, 值)
}

// D设置JSON 设置路径为 JSON 片段（须可解析）。
func 设置JSON(文本, 路径, JSON片段 string) (string, error) {
	值, 错误 := 解析(JSON片段)
	if 错误 != nil {
		return "", fmt.Errorf("JSON片段无效: %w", 错误)
	}
	return 设置并编码(文本, 路径, 值)
}

// D压缩 去掉空白，非法返回空串。
func 压缩(文本 string) string {
	var 缓冲 bytes.Buffer
	if 错误 := json.Compact(&缓冲, []byte(文本)); 错误 != nil {
		return ""
	}
	return 缓冲.String()
}

// D格式化 缩进 2 空格美化，非法返回空串。
func 格式化(文本 string) string {
	var 缓冲 bytes.Buffer
	if 错误 := json.Indent(&缓冲, []byte(文本), "", "  "); 错误 != nil {
		return ""
	}
	return 缓冲.String()
}
