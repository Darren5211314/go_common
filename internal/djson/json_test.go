package djson

import (
	"strings"
	"testing"

	"github.com/Darren5211314/go_common/internal/dlog"
)

const 样本 = `{
  "code": 0,
  "msg": "ok",
  "data": {
    "user": {"id": 1001, "name": "张三", "vip": true},
    "tags": ["go", "json"],
    "score": 98.5
  }
}`

func Test取值与判断(t *testing.T) {
	dlog.Info("是否有效", "结果", 是否有效(样本))
	if !是否有效(样本) {
		t.Errorf("是否有效 不符")
	}
	if 是否有效("{bad") {
		t.Errorf("非法 JSON 不应有效")
	}

	存在 := 是否存在(样本, "data.user.name")
	dlog.Info("是否存在", "路径", "data.user.name", "结果", 存在)
	if !存在 || 是否存在(样本, "data.nope") {
		t.Errorf("是否存在 不符")
	}

	成员数 := 取成员数(样本, "data.tags")
	dlog.Info("取成员数", "tags", 成员数, "user", 取成员数(样本, "data.user"), "根", 取成员数(样本, ""))
	if 成员数 != 2 || 取成员数(样本, "data.user") != 3 || 取成员数(样本, "") != 3 {
		t.Errorf("取成员数 不符")
	}

	名 := 取文本(样本, "data.user.name")
	标签 := 取文本(样本, "data.tags.0")
	码文本 := 取文本(样本, "code")
	dlog.Info("取文本", "name", 名, "tag0", 标签, "code", 码文本)
	if 名 != "张三" || 标签 != "go" || 码文本 != "0" {
		t.Errorf("取文本 不符")
	}

	id := 取整数(样本, "data.user.id")
	分 := 取小数(样本, "data.score")
	vip := 取逻辑(样本, "data.user.vip")
	dlog.Info("取类型值", "id", id, "score", 分, "vip", vip)
	if id != 1001 || 分 != 98.5 || !vip {
		t.Errorf("取整数/小数/逻辑 不符")
	}

	子树 := 取JSON(样本, "data.user")
	dlog.Info("取JSON", "结果", 子树)
	if !strings.Contains(子树, `"name":"张三"`) || !strings.Contains(子树, `"id":1001`) {
		t.Errorf("取JSON 不符")
	}
}

func Test设置与整形(t *testing.T) {
	s, err := 设置文本(样本, "data.user.name", "李四")
	dlog.Info("设置文本", "结果", s, "错误", err)
	if err != nil || 取文本(s, "data.user.name") != "李四" {
		t.Errorf("设置文本 不符")
	}

	s, err = 设置整数(样本, "code", 1)
	if err != nil || 取整数(s, "code") != 1 {
		t.Errorf("设置整数 不符")
	}

	s, err = 设置小数(样本, "data.score", 100)
	if err != nil || 取小数(s, "data.score") != 100 {
		t.Errorf("设置小数 不符")
	}

	s, err = 设置逻辑(样本, "data.user.vip", false)
	if err != nil || 取逻辑(s, "data.user.vip") {
		t.Errorf("设置逻辑 不符")
	}

	s, err = 设置JSON(样本, "data.tags", `["sm3","sm4"]`)
	dlog.Info("设置JSON", "结果", s, "错误", err)
	if err != nil || 取文本(s, "data.tags.0") != "sm3" || 取成员数(s, "data.tags") != 2 {
		t.Errorf("设置JSON 不符")
	}

	新建, err := 设置文本("", "a.b", "x")
	dlog.Info("空根设置", "结果", 新建, "错误", err)
	if err != nil || 取文本(新建, "a.b") != "x" {
		t.Errorf("空根设置 不符")
	}

	紧凑 := 压缩(样本)
	美化 := 格式化(紧凑)
	dlog.Info("压缩格式化", "压缩", 紧凑, "美化含换行", strings.Contains(美化, "\n"))
	if 紧凑 == "" || !strings.Contains(美化, "\n") {
		t.Errorf("压缩/格式化 不符")
	}

	对象根 := `{"0":"a","1":"b"}`
	改, err := 设置文本(对象根, "0", "x")
	dlog.Info("对象数字键", "结果", 改, "错误", err)
	if err != nil || 取文本(改, "0") != "x" || 取文本(改, "1") != "b" {
		t.Fatalf("对象数字键不应冲成数组: %s", 改)
	}
	if 是否有效(改) && strings.HasPrefix(strings.TrimSpace(改), "[") {
		t.Fatalf("结果不应变成数组: %s", 改)
	}

	数组根 := `["a","b"]`
	改数组, err := 设置文本(数组根, "0", "z")
	if err != nil || 取文本(改数组, "0") != "z" || 取文本(改数组, "1") != "b" {
		t.Fatalf("数组下标行为应保持: %s %v", 改数组, err)
	}
}
