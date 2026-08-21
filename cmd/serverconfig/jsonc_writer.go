package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// MarshalJSONC 把 v 序列化为两空格缩进的 JSONC 文本。不带注释时输出与
// json.MarshalIndent(v, "", "  ") 逐字节一致；comments 按路径注入中文注释：
//
//   - "game.homeMap" 匹配对象成员（点号连接各层键名）；
//   - "items[].name" 匹配数组元素内的字段（[] 代表数组元素层）。
//
// commentAll 为 false 时元素注释只渲染在各数组的首个元素上（避免
// std_items 这类大文件注释爆炸），为 true 时逐元素渲染。
//
// 注释位置规则：叶子字段（数字/字符串/布尔/null）用行尾注释，
// 对象/非空数组成员用键上方整行注释。
func MarshalJSONC(v any, comments map[string]string, commentAll bool) (string, error) {
	w := &jsoncWriter{comments: comments, commentAll: commentAll}
	if err := w.writeValue(reflect.ValueOf(v), "", 0, true); err != nil {
		return "", err
	}
	return w.sb.String(), nil
}

type jsoncWriter struct {
	sb         strings.Builder
	comments   map[string]string
	commentAll bool
}

const (
	placeAbove    = true
	placeTrailing = false
)

func (w *jsoncWriter) writeValue(v reflect.Value, path string, depth int, active bool) error {
	v = unwrap(v)
	if !v.IsValid() {
		w.sb.WriteString("null")
		return nil
	}
	switch v.Kind() {
	case reflect.Struct:
		return w.writeStruct(v, path, depth, active)
	case reflect.Map:
		return w.writeMap(v, path, depth, active)
	case reflect.Slice, reflect.Array:
		return w.writeSlice(v, path, depth, active)
	default:
		return w.writeLeaf(v)
	}
}

func (w *jsoncWriter) writeLeaf(v reflect.Value) error {
	b, err := json.Marshal(v.Interface())
	if err != nil {
		return err
	}
	w.sb.Write(b)
	return nil
}

type member struct {
	name  string
	value reflect.Value
}

func (w *jsoncWriter) writeStruct(v reflect.Value, path string, depth int, active bool) error {
	var fields []member
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		name, omitempty := parseJSONTag(sf.Tag.Get("json"))
		if name == "-" {
			continue
		}
		if name == "" {
			name = sf.Name
		}
		fv := v.Field(i)
		if omitempty && isEmptyValue(fv) {
			continue
		}
		fields = append(fields, member{name: name, value: fv})
	}
	if len(fields) == 0 {
		w.sb.WriteString("{}")
		return nil
	}
	w.sb.WriteString("{\n")
	indent := strings.Repeat("  ", depth+1)
	for i, f := range fields {
		memberPath := joinPath(path, f.name)
		above := placementOf(f.value) == placeAbove
		if active && above {
			w.writeAboveComment(memberPath, indent)
		}
		w.sb.WriteString(indent)
		keyJSON, err := json.Marshal(f.name)
		if err != nil {
			return err
		}
		w.sb.Write(keyJSON)
		w.sb.WriteString(": ")
		if err := w.writeValue(f.value, memberPath, depth+1, active); err != nil {
			return err
		}
		if i < len(fields)-1 {
			w.sb.WriteString(",")
		}
		if active && !above {
			w.writeTrailingComment(memberPath)
		}
		w.sb.WriteString("\n")
	}
	w.sb.WriteString(strings.Repeat("  ", depth))
	w.sb.WriteString("}")
	return nil
}

func (w *jsoncWriter) writeMap(v reflect.Value, path string, depth int, active bool) error {
	if v.IsNil() {
		w.sb.WriteString("null")
		return nil
	}
	keys := v.MapKeys()
	names := make([]string, 0, len(keys))
	for _, k := range keys {
		if k.Kind() != reflect.String {
			return fmt.Errorf("jsonc: 不支持的 map 键类型 %s", k.Type())
		}
		names = append(names, k.String())
	}
	if len(names) == 0 {
		w.sb.WriteString("{}")
		return nil
	}
	sort.Strings(names) // 与 encoding/json 一致：字典序
	w.sb.WriteString("{\n")
	indent := strings.Repeat("  ", depth+1)
	for i, name := range names {
		memberPath := joinPath(path, name)
		mv := v.MapIndex(reflect.ValueOf(name))
		above := placementOf(mv) == placeAbove
		if active && above {
			w.writeAboveComment(memberPath, indent)
		}
		w.sb.WriteString(indent)
		keyJSON, err := json.Marshal(name)
		if err != nil {
			return err
		}
		w.sb.Write(keyJSON)
		w.sb.WriteString(": ")
		if err := w.writeValue(mv, memberPath, depth+1, active); err != nil {
			return err
		}
		if i < len(names)-1 {
			w.sb.WriteString(",")
		}
		if active && !above {
			w.writeTrailingComment(memberPath)
		}
		w.sb.WriteString("\n")
	}
	w.sb.WriteString(strings.Repeat("  ", depth))
	w.sb.WriteString("}")
	return nil
}

func (w *jsoncWriter) writeSlice(v reflect.Value, path string, depth int, active bool) error {
	if v.Kind() == reflect.Slice && v.IsNil() {
		w.sb.WriteString("null")
		return nil
	}
	n := v.Len()
	if n == 0 {
		w.sb.WriteString("[]")
		return nil
	}
	w.sb.WriteString("[\n")
	indent := strings.Repeat("  ", depth+1)
	elemPath := path + "[]"
	for i := 0; i < n; i++ {
		elemActive := active && (i == 0 || w.commentAll)
		above := placementOf(v.Index(i)) == placeAbove
		if elemActive && above {
			w.writeAboveComment(elemPath, indent)
		}
		w.sb.WriteString(indent)
		if err := w.writeValue(v.Index(i), elemPath, depth+1, elemActive); err != nil {
			return err
		}
		if i < n-1 {
			w.sb.WriteString(",")
		}
		if elemActive && !above {
			w.writeTrailingComment(elemPath)
		}
		w.sb.WriteString("\n")
	}
	w.sb.WriteString(strings.Repeat("  ", depth))
	w.sb.WriteString("]")
	return nil
}

func (w *jsoncWriter) writeAboveComment(path, indent string) {
	c := w.comments[path]
	if c == "" {
		return
	}
	for _, line := range strings.Split(c, "\n") {
		w.sb.WriteString(indent)
		w.sb.WriteString("// ")
		w.sb.WriteString(line)
		w.sb.WriteString("\n")
	}
}

func (w *jsoncWriter) writeTrailingComment(path string) {
	c := w.comments[path]
	if c == "" {
		return
	}
	w.sb.WriteString(" // ")
	w.sb.WriteString(strings.ReplaceAll(c, "\n", " "))
}

// placementOf 决定注释位置：对象与非空数组/映射占多行，注释放键上方；
// 叶子值与单行输出（null/[]/{}）注释放行尾。
func placementOf(v reflect.Value) bool {
	u := unwrap(v)
	if !u.IsValid() {
		return placeTrailing // null
	}
	switch u.Kind() {
	case reflect.Struct:
		return placeAbove
	case reflect.Map, reflect.Slice, reflect.Array:
		if u.Len() == 0 {
			return placeTrailing // null 或 []/{}
		}
		return placeAbove
	default:
		return placeTrailing
	}
}

// unwrap 解引用指针/接口；nil 返回无效 Value（序列化为 null）。
func unwrap(v reflect.Value) reflect.Value {
	for v.IsValid() && (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}

func parseJSONTag(tag string) (name string, omitempty bool) {
	if tag == "" {
		return "", false
	}
	parts := strings.Split(tag, ",")
	name = parts[0]
	for _, opt := range parts[1:] {
		if opt == "omitempty" {
			omitempty = true
		}
	}
	return
}

// isEmptyValue 与 encoding/json 的 omitempty 判空逻辑一致。
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func joinPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}
