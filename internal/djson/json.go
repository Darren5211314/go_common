package djson

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidJSON = errors.New("JSON无效")

func parse(text string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var root any
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}
	return root, nil
}

func splitPath(path string) []string {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	segs := strings.Split(path, ".")
	result := make([]string, 0, len(segs))
	for _, seg := range segs {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		result = append(result, seg)
	}
	return result
}

func load(root any, path string) (any, bool) {
	segs := splitPath(path)
	cur := root
	for _, seg := range segs {
		keyLen, ok := cur.(map[string]any)
		if ok {
			next, exists := keyLen[seg]
			if !exists {
				return nil, false
			}
			cur = next
			continue
		}
		arr, ok := cur.([]any)
		if !ok {
			return nil, false
		}
		index, err := strconv.Atoi(seg)
		if err != nil || index < 0 || index >= len(arr) {
			return nil, false
		}
		cur = arr[index]
	}
	return cur, true
}

func compactEncode(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Valid 判断文本是否为合法 JSON。
func isValidIP(text string) bool {
	_, err := parse(text)
	return err == nil
}

// Has 判断路径是否存在。
func existsAt(text, path string) bool {
	root, err := parse(text)
	if err != nil {
		return false
	}
	_, exists := load(root, path)
	return exists
}

// Count 返回对象字段数或数组长度；路径空则对根取值；无效返回 0。
func memberCount(text, path string) int {
	root, err := parse(text)
	if err != nil {
		return 0
	}
	value, exists := load(root, path)
	if !exists {
		return 0
	}
	switch obj := value.(type) {
	case map[string]any:
		return len(obj)
	case []any:
		return len(obj)
	default:
		return 0
	}
}

// GetString 取路径对应值并转为文本；找不到返回空串。
func getText(text, path string) string {
	root, err := parse(text)
	if err != nil {
		return ""
	}
	value, exists := load(root, path)
	if !exists || value == nil {
		return ""
	}
	switch v := value.(type) {
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
		encode, err := compactEncode(v)
		if err != nil {
			return ""
		}
		return encode
	}
}

// GetInt 取路径对应整数；找不到或非数字返回 0。
func getInt(text, path string) int64 {
	root, err := parse(text)
	if err != nil {
		return 0
	}
	value, exists := load(root, path)
	if !exists || value == nil {
		return 0
	}
	switch v := value.(type) {
	case json.Number:
		i, err := v.Int64()
		if err == nil {
			return i
		}
		f, err := v.Float64()
		if err != nil {
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
		i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			f, err2 := strconv.ParseFloat(strings.TrimSpace(v), 64)
			if err2 != nil {
				return 0
			}
			return int64(f)
		}
		return i
	default:
		return 0
	}
}

// GetFloat 取路径对应浮点数；找不到或非数字返回 0。
func getFloat(text, path string) float64 {
	root, err := parse(text)
	if err != nil {
		return 0
	}
	value, exists := load(root, path)
	if !exists || value == nil {
		return 0
	}
	switch v := value.(type) {
	case json.Number:
		f, err := v.Float64()
		if err != nil {
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
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}

// GetBool 取路径对应布尔；找不到返回 false。数字非 0 视为 true；文本 "true"/"1" 为 true。
func getBool(text, path string) bool {
	root, err := parse(text)
	if err != nil {
		return false
	}
	value, exists := load(root, path)
	if !exists || value == nil {
		return false
	}
	switch v := value.(type) {
	case bool:
		return v
	case json.Number:
		f, err := v.Float64()
		if err != nil {
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

// GetJSON 取路径对应子树的紧凑 JSON；找不到返回空串。
func getJSON(text, path string) string {
	root, err := parse(text)
	if err != nil {
		return ""
	}
	value, exists := load(root, path)
	if !exists {
		return ""
	}
	encode, err := compactEncode(value)
	if err != nil {
		return ""
	}
	return encode
}

func store(root any, path string, newVal any) (any, error) {
	segs := splitPath(path)
	if len(segs) == 0 {
		return newVal, nil
	}
	if root == nil {
		root = map[string]any{}
	}
	return setRecursive(root, segs, newVal)
}

func setRecursive(cur any, segs []string, newVal any) (any, error) {
	seg := segs[0]
	last := len(segs) == 1
	index, isArrayIndex := tryParseIndex(seg)

	// 已是对象时，数字段也当对象键，避免把 {"0":...} 冲成数组。
	if _, isObject := cur.(map[string]any); isObject && isArrayIndex {
		isArrayIndex = false
	}

	if isArrayIndex {
		arr, ok := cur.([]any)
		if !ok {
			arr = []any{}
		}
		for len(arr) <= index {
			arr = append(arr, nil)
		}
		if last {
			arr[index] = newVal
			return arr, nil
		}
		next := arr[index]
		next, err := setRecursive(ensureChild(next, segs[1]), segs[1:], newVal)
		if err != nil {
			return nil, err
		}
		arr[index] = next
		return arr, nil
	}

	obj, ok := cur.(map[string]any)
	if !ok {
		obj = map[string]any{}
	}
	if last {
		obj[seg] = newVal
		return obj, nil
	}
	next := obj[seg]
	next, err := setRecursive(ensureChild(next, segs[1]), segs[1:], newVal)
	if err != nil {
		return nil, err
	}
	obj[seg] = next
	return obj, nil
}

func tryParseIndex(seg string) (int, bool) {
	index, err := strconv.Atoi(seg)
	if err != nil || index < 0 {
		return 0, false
	}
	return index, true
}

func ensureChild(existing any, nextSeg string) any {
	if existing != nil {
		if _, isObject := existing.(map[string]any); isObject {
			return existing
		}
		if _, isArray := existing.([]any); isArray {
			return existing
		}
	}
	if _, isIndex := tryParseIndex(nextSeg); isIndex {
		return []any{}
	}
	return map[string]any{}
}

func setAndEncode(text, path string, newVal any) (string, error) {
	root, err := parse(text)
	if err != nil {
		// 空文本视为空对象，便于从零构建
		if strings.TrimSpace(text) == "" {
			root = map[string]any{}
		} else {
			return "", errInvalidJSON
		}
	}
	newRoot, err := store(root, path, newVal)
	if err != nil {
		return "", err
	}
	return compactEncode(newRoot)
}

// SetString 设置路径为字符串。
func setText(text, path, value string) (string, error) {
	return setAndEncode(text, path, value)
}

// SetInt 设置路径为整数。
func setInt(text, path string, value int64) (string, error) {
	return setAndEncode(text, path, value)
}

// SetFloat 设置路径为小数。
func setFloat(text, path string, value float64) (string, error) {
	return setAndEncode(text, path, value)
}

// SetBool 设置路径为布尔。
func setBool(text, path string, value bool) (string, error) {
	return setAndEncode(text, path, value)
}

// SetJSON 设置路径为 JSON 片段（须可解析）。
func setJSON(text, path, jsonFrag string) (string, error) {
	value, err := parse(jsonFrag)
	if err != nil {
		return "", fmt.Errorf("JSON片段无效: %w", err)
	}
	return setAndEncode(text, path, value)
}

// Compact 去掉空白，非法返回空串。
func compact(text string) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(text)); err != nil {
		return ""
	}
	return buf.String()
}

// Pretty 缩进 2 空格美化，非法返回空串。
func pretty(text string) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(text), "", "  "); err != nil {
		return ""
	}
	return buf.String()
}
