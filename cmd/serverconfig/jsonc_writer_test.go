package main

import (
	"encoding/json"
	"strings"
	"testing"
)

type writerTestInner struct {
	Name  string   `json:"name"`
	Props []string `json:"props,omitempty"`
	Val   int      `json:"val"`
}

type writerTestOuter struct {
	Title   string            `json:"title"`
	Count   int64             `json:"count"`
	Ratio   float64           `json:"ratio"`
	Enabled bool              `json:"enabled"`
	Inner   writerTestInner   `json:"inner"`
	Items   []writerTestInner `json:"items"`
	Nums    []int             `json:"nums"`
	Nulls   []string          `json:"nulls"`
	Empty   []string          `json:"empty"`
	Tags    map[string]string `json:"tags"`
	Scores  map[string]int64  `json:"scores"`
	Mixed   map[string]any    `json:"mixed"`
	Text    string            `json:"text"`
	Omit    int               `json:"omit,omitempty"`

	private int
}

func writerTestFixture() writerTestOuter {
	return writerTestOuter{
		Title:   "服务器<配置> & \"测试\"",
		Count:   1234567890123,
		Ratio:   1.5,
		Enabled: true,
		Inner:   writerTestInner{Name: "内层", Props: []string{"a", "b"}, Val: 7},
		Items: []writerTestInner{
			{Name: "物品一", Val: 1},
			{Name: "物品二", Props: []string{"x"}, Val: 2},
		},
		Nums:    []int{3, 1, 2},
		Nulls:   nil,
		Empty:   []string{},
		Tags:    map[string]string{"z": "尾", "a": "首", "中": "中文键"},
		Scores:  map[string]int64{"b": 2, "a": 1},
		Mixed:   map[string]any{"_source": "来源", "list": []writerTestInner{{Name: "嵌入", Val: 9}}, "none": nil},
		Text:    "第一行\n第二行\t制表符 <html> & 引号\"结束",
		private: 42,
	}
}

// 无注释时输出必须与 json.MarshalIndent 逐字节一致，
// 保证替换 MarshalIndent 后数据部分不漂移。
func TestMarshalJSONCEquivalentToMarshalIndent(t *testing.T) {
	fixture := writerTestFixture()
	got, err := MarshalJSONC(fixture, nil, false)
	if err != nil {
		t.Fatalf("MarshalJSONC: %v", err)
	}
	wantBytes, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent: %v", err)
	}
	if got != string(wantBytes) {
		t.Errorf("输出与 MarshalIndent 不一致\n--- got ---\n%s\n--- want ---\n%s", got, wantBytes)
	}
}

// map[string]interface{} 包装（convert_*.go 的主流形态）同样逐字节一致。
func TestMarshalJSONCMapWrapperEquivalent(t *testing.T) {
	fixture := map[string]interface{}{
		"_source":      "asset/server/Envir/MonGen.txt",
		"_description": "刷怪配置",
		"items":        []writerTestInner{{Name: "a", Val: 1}, {Name: "b", Val: 2}},
		"count":        2,
	}
	got, err := MarshalJSONC(fixture, nil, false)
	if err != nil {
		t.Fatalf("MarshalJSONC: %v", err)
	}
	wantBytes, err := json.MarshalIndent(fixture, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent: %v", err)
	}
	if got != string(wantBytes) {
		t.Errorf("输出与 MarshalIndent 不一致\n--- got ---\n%s\n--- want ---\n%s", got, wantBytes)
	}
}

type smallItem struct {
	Name string `json:"name"`
	Val  int    `json:"val"`
}

type smallDoc struct {
	Title string      `json:"title"`
	Items []smallItem `json:"items"`
	Nums  []int       `json:"nums"`
}

func smallFixture() smallDoc {
	return smallDoc{
		Title: "t",
		Items: []smallItem{{Name: "a", Val: 1}, {Name: "b", Val: 2}},
		Nums:  []int{1, 2},
	}
}

var smallComments = map[string]string{
	"title":        "标题",
	"items":        "物品列表",
	"items[].name": "名称",
	"items[].val":  "数值",
	"nums":         "数字",
	"nums[]":       "单个数字",
}

// commentAll=false：元素注释只渲染在首元素。
func TestMarshalJSONCCommentsFirstElement(t *testing.T) {
	got, err := MarshalJSONC(smallFixture(), smallComments, false)
	if err != nil {
		t.Fatalf("MarshalJSONC: %v", err)
	}
	want := `{
  "title": "t", // 标题
  // 物品列表
  "items": [
    {
      "name": "a", // 名称
      "val": 1 // 数值
    },
    {
      "name": "b",
      "val": 2
    }
  ],
  // 数字
  "nums": [
    1, // 单个数字
    2
  ]
}`
	if got != want {
		t.Errorf("输出不符\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

// commentAll=true：元素注释逐元素渲染。
func TestMarshalJSONCCommentsAllElements(t *testing.T) {
	got, err := MarshalJSONC(smallFixture(), smallComments, true)
	if err != nil {
		t.Fatalf("MarshalJSONC: %v", err)
	}
	if strings.Count(got, "// 名称") != 2 || strings.Count(got, "// 单个数字") != 2 {
		t.Errorf("commentAll=true 时元素注释应逐元素出现\n%s", got)
	}
}

// 多行注释在键上方逐行展开；行尾注释不换行。
func TestMarshalJSONCMultilineComment(t *testing.T) {
	comments := map[string]string{
		"items":       "第一行说明\n第二行说明",
		"items[].val": "数值\n不应换行",
	}
	got, err := MarshalJSONC(smallFixture(), comments, false)
	if err != nil {
		t.Fatalf("MarshalJSONC: %v", err)
	}
	if !strings.Contains(got, "  // 第一行说明\n  // 第二行说明\n  \"items\": [") {
		t.Errorf("多行上方注释展开错误\n%s", got)
	}
	if !strings.Contains(got, `"val": 1 // 数值 不应换行`) {
		t.Errorf("行尾注释应压成单行\n%s", got)
	}
}

// nil 切片输出 null，行尾注释可用；空容器输出 {} / []。
func TestMarshalJSONCNullAndEmpty(t *testing.T) {
	type doc struct {
		NilSlice []string          `json:"nilSlice"`
		NilMap   map[string]string `json:"nilMap"`
		EmptyMap map[string]string `json:"emptyMap"`
	}
	v := doc{NilSlice: nil, NilMap: nil, EmptyMap: map[string]string{}}
	comments := map[string]string{"nilSlice": "可空列表", "nilMap": "可空映射"}
	got, err := MarshalJSONC(v, comments, false)
	if err != nil {
		t.Fatalf("MarshalJSONC: %v", err)
	}
	want := `{
  "nilSlice": null, // 可空列表
  "nilMap": null, // 可空映射
  "emptyMap": {}
}`
	if got != want {
		t.Errorf("输出不符\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}
