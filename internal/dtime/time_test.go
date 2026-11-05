package dtime

import (
	"github.com/Darren5211314/go_common/internal/dlog"
	"strings"
	"testing"
	"time"
)

func Test通配符转布局(t *testing.T) {
	用例 := []struct {
		通配符, 期望 string
	}{
		{"yyyy-mm-dd hh:mm:ss", "2006-01-02 15:04:05"},
		{"yyyy/mm/dd", "2006/01/02"},
		{"yyyy-MM-dd HH:mm:ss", "2006-01-02 15:04:05"},
		{"yyyy年mm月dd日", "2006年01月02日"},
	}
	for _, 项 := range 用例 {
		t.Run(项.通配符, func(t *testing.T) {
			得到 := 通配符转布局(项.通配符)
			dlog.Info("通配符转布局", "通配符", 项.通配符, "结果", 得到, "期望", 项.期望)
			if 得到 != 项.期望 {
				t.Errorf("通配符转布局 不符")
			}
		})
	}
}

func TestD取现行时间戳(t *testing.T) {
	前秒 := time.Now().Unix()
	十位 := 取现行时间戳(true)
	后秒 := time.Now().Unix()
	dlog.Info("取现行时间戳", "取十位", true, "结果", 十位, "范围起", 前秒, "范围止", 后秒)
	if 十位 < 前秒 || 十位 > 后秒 {
		t.Errorf("十位时间戳 不符")
	}

	前毫秒 := time.Now().UnixMilli()
	十三位 := 取现行时间戳(false)
	后毫秒 := time.Now().UnixMilli()
	dlog.Info("取现行时间戳", "取十位", false, "结果", 十三位, "范围起", 前毫秒, "范围止", 后毫秒)
	if 十三位 < 前毫秒 || 十三位 > 后毫秒 {
		t.Errorf("十三位时间戳 不符")
	}
}

func TestD时间戳转时间(t *testing.T) {
	固定 := time.Date(2026, 11, 5, 20, 53, 0, 0, time.Local)
	期望 := "2026-11-05 20:53:00"
	格式 := "yyyy-mm-dd hh:mm:ss"

	秒结果 := 时间戳转时间(固定.Unix(), 格式)
	dlog.Info("时间戳转时间", "戳类型", "秒", "戳", 固定.Unix(), "格式", 格式, "结果", 秒结果, "期望", 期望)
	if 秒结果 != 期望 {
		t.Errorf("秒戳 不符")
	}

	毫秒结果 := 时间戳转时间(固定.UnixMilli(), 格式)
	dlog.Info("时间戳转时间", "戳类型", "毫秒", "戳", 固定.UnixMilli(), "格式", 格式, "结果", 毫秒结果, "期望", 期望)
	if 毫秒结果 != 期望 {
		t.Errorf("毫秒戳 不符")
	}
}

func TestD时间转时间戳(t *testing.T) {
	期望秒 := time.Date(2026, 11, 5, 20, 53, 0, 0, time.Local).Unix()
	期望毫秒 := 期望秒 * 1000

	得到 := 时间转时间戳("2026-11-05 20:53:00", true)
	dlog.Info("时间转时间戳", "输入", "2026-11-05 20:53:00", "取十位", true, "结果", 得到, "期望", 期望秒)
	if 得到 != 期望秒 {
		t.Errorf("横杠取十位 不符")
	}

	得到 = 时间转时间戳("2026/11/05 20:53:00", false)
	dlog.Info("时间转时间戳", "输入", "2026/11/05 20:53:00", "取十位", false, "结果", 得到, "期望", 期望毫秒)
	if 得到 != 期望毫秒 {
		t.Errorf("斜杠取十三位 不符")
	}

	期望日 := time.Date(2026, 11, 5, 0, 0, 0, 0, time.Local).Unix()
	得到 = 时间转时间戳("2026-11-05", true)
	dlog.Info("时间转时间戳", "输入", "2026-11-05", "取十位", true, "结果", 得到, "期望", 期望日)
	if 得到 != 期望日 {
		t.Errorf("仅日期 不符")
	}

	得到 = 时间转时间戳("坏数据", false)
	dlog.Info("时间转时间戳", "输入", "坏数据", "取十位", false, "结果", 得到, "期望", 0)
	if 得到 != 0 {
		t.Errorf("解析失败应返回 0")
	}
}

func TestD时间格式化(t *testing.T) {
	得到 := 时间格式化("2026/11/05", "yyyy-mm-dd")
	dlog.Info("时间格式化", "输入", "2026/11/05", "目标格式", "yyyy-mm-dd", "结果", 得到, "期望", "2026-11-05")
	if 得到 != "2026-11-05" {
		t.Errorf("斜杠转横杠 不符")
	}

	完整 := 时间格式化("2026-11-05 20:53:00", "yyyy年mm月dd日")
	dlog.Info("时间格式化", "输入", "2026-11-05 20:53:00", "目标格式", "yyyy年mm月dd日", "结果", 完整, "期望", "2026年11月05日")
	if 完整 != "2026年11月05日" {
		t.Errorf("中文日期 不符")
	}

	现行 := 时间格式化("", "yyyy-mm-dd")
	今天 := time.Now().Format("2006-01-02")
	dlog.Info("时间格式化现行", "输入", "", "目标格式", "yyyy-mm-dd", "结果", 现行, "期望", 今天)
	if 现行 != 今天 {
		t.Errorf("空文本现行时间 不符")
	}

	现行完整 := 时间格式化("", "yyyy-mm-dd hh:mm:ss")
	dlog.Info("时间格式化现行完整", "结果", 现行完整, "应含", 今天)
	if !strings.Contains(现行完整, 今天) {
		t.Errorf("空文本完整格式应包含今天日期")
	}

	失败 := 时间格式化("坏", "yyyy-mm-dd")
	dlog.Info("时间格式化失败", "输入", "坏", "结果", 失败, "期望", "")
	if 失败 != "" {
		t.Errorf("解析失败应为空")
	}
}

func Test自产格式回环解析(t *testing.T) {
	用例 := []struct {
		输入   string
		取十位  bool
		期望秒  int64
	}{
		{"2024-03-05 06:07", true, time.Date(2024, 3, 5, 6, 7, 0, 0, time.Local).Unix()},
		{"2024/03/05 06:07", true, time.Date(2024, 3, 5, 6, 7, 0, 0, time.Local).Unix()},
		{"2024年03月05日 06时07分08秒", true, time.Date(2024, 3, 5, 6, 7, 8, 0, time.Local).Unix()},
		{"2024年03月05日", true, time.Date(2024, 3, 5, 0, 0, 0, 0, time.Local).Unix()},
	}
	for _, 项 := range 用例 {
		得到 := 时间转时间戳(项.输入, 项.取十位)
		dlog.Info("自产回环", "输入", 项.输入, "结果", 得到, "期望", 项.期望秒)
		if 得到 != 项.期望秒 {
			t.Errorf("%s: 得到 %d 期望 %d", 项.输入, 得到, 项.期望秒)
		}
	}
	仅时分 := 时间格式化("2024-03-05 06:07:08", "hh:mm")
	if 仅时分 != "06:07" {
		t.Fatalf("格式化到 hh:mm 不符: %q", 仅时分)
	}
	回戳 := 时间转时间戳("06:07", true)
	if 回戳 == 0 {
		t.Fatal("仅时分应能解析")
	}
}

func TestD延时(t *testing.T) {
	开始0 := time.Now()
	延时(0)
	if time.Since(开始0) > 20*time.Millisecond {
		t.Errorf("延时(0) 应几乎立即返回")
	}

	开始 := time.Now()
	延时(50)
	耗时 := time.Since(开始)
	dlog.Info("延时", "毫秒", 50, "耗时", 耗时)
	if 耗时 < 40*time.Millisecond {
		t.Errorf("延时(50) 耗时过短: %v", 耗时)
	}
}
