package rss

//
//import (
//	"github.com/microcosm-cc/bluemonday"
//	"reflect"
//	"rss-feed/pkg/http"
//	"testing"
//)
//
//func TestHtmlSanitizer_Name(t *testing.T) {
//	tests := []struct {
//		name string
//		want string
//	}{
//		{
//			name: "Name must be html-sanitizer",
//			want: "html-sanitizer",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := NewHtmlSanitizer()
//
//			if got := h.Name(); got != tt.want {
//				t.Errorf("Name() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestHtmlSanitizer_Process(t *testing.T) {
//	type fields struct {
//		policy *bluemonday.Policy
//	}
//	type args struct {
//		items *[]http.Item
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *[]http.Item
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &HtmlSanitizer{
//				policy: tt.fields.policy,
//			}
//			got, err := h.Process(tt.args.items)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Process() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNewHtmlSanitizer(t *testing.T) {
//	tests := []struct {
//		name string
//		want *HtmlSanitizer
//	}{
//		{
//			want: &HtmlSanitizer{
//				policy: bluemonday.NewPolicy(),
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewHtmlSanitizer(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewHtmlSanitizer() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
///// MaxLengthProcessor
//
//func TestMaxLengthProcessor_Name(t *testing.T) {
//	tests := []struct {
//		name string
//		want string
//	}{
//		{
//			name: "Max Length Must Be max-length",
//			want: "max-length",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := NewMaxLengthProcessor()
//			if got := h.Name(); got != tt.want {
//				t.Errorf("Name() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMaxLengthProcessor_Process(t *testing.T) {
//	type fields struct {
//		maxLength int
//		postfix   string
//	}
//	type args struct {
//		items *[]http.Item
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *[]http.Item
//		wantErr bool
//	}{
//		{
//			name: "no truncation needed",
//			fields: fields{
//				maxLength: 128,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{{Description: "short text"}},
//			},
//			want:    &[]http.Item{{Description: "short text"}},
//			wantErr: false,
//		},
//		{
//			name: "truncation needed",
//			fields: fields{
//				maxLength: 20,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{{Description: "this is a very long text that needs to be truncated"}},
//			},
//			want:    &[]http.Item{{Description: "this is a very long ..."}},
//			wantErr: false,
//		},
//		{
//			name: "empty postfix",
//			fields: fields{
//				maxLength: 10,
//				postfix:   "",
//			},
//			args: args{
//				items: &[]http.Item{{Description: "this is a very long text"}},
//			},
//			want:    &[]http.Item{{Description: "this is a "}},
//			wantErr: false,
//		},
//		{
//			name: "unicode support",
//			fields: fields{
//				maxLength: 7,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{{Description: "привет, мир!"}},
//			},
//			want:    &[]http.Item{{Description: "привет,..."}},
//			wantErr: false,
//		},
//		{
//			name: "multiple items",
//			fields: fields{
//				maxLength: 10,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{
//					{Description: "short text"},
//					{Description: "this is a very long text"},
//				},
//			},
//			want: &[]http.Item{
//				{Description: "short text"},
//				{Description: "this is a ..."},
//			},
//			wantErr: false,
//		},
//		{
//			name: "empty list",
//			fields: fields{
//				maxLength: 128,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{},
//			},
//			want:    &[]http.Item{},
//			wantErr: false,
//		},
//		{
//			name: "maxLength zero",
//			fields: fields{
//				maxLength: 0,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{{Description: "this is a text"}},
//			},
//			want:    &[]http.Item{{Description: "this is a text"}},
//			wantErr: false,
//		},
//		{
//			name: "negative maxLength",
//			fields: fields{
//				maxLength: -10,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{{Description: "this is a text"}},
//			},
//			want:    &[]http.Item{{Description: "this is a text"}},
//			wantErr: false,
//		},
//		{
//			name: "text equals maxLength",
//			fields: fields{
//				maxLength: 12,
//				postfix:   "...",
//			},
//			args: args{
//				items: &[]http.Item{{Description: "exact length"}},
//			},
//			want:    &[]http.Item{{Description: "exact length"}},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &MaxLengthProcessor{
//				maxLength: tt.fields.maxLength,
//				postfix:   tt.fields.postfix,
//			}
//			got, err := h.Process(tt.args.items)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Process() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestMaxLengthProcessor_Truncate(t *testing.T) {
//	type fields struct {
//		maxLength int
//		postfix   string
//	}
//	type args struct {
//		str string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   string
//	}{
//		{
//			name: "no truncation needed",
//			fields: fields{
//				maxLength: 128,
//				postfix:   "^^^",
//			},
//			args: args{
//				str: "short text",
//			},
//			want: "short text",
//		},
//		{
//			name: "Truncation needed",
//			fields: fields{
//				maxLength: 10,
//				postfix:   "^^^",
//			},
//			args: args{
//				str: "short text, short text, short text",
//			},
//			want: "short text^^^",
//		},
//		{
//			name: "Empty postfix",
//			fields: fields{
//				maxLength: 10,
//				postfix:   "",
//			},
//			args: args{
//				str: "short text, short text, short text",
//			},
//			want: "short text",
//		},
//		{
//			name: "unicode dotted string ",
//			fields: fields{
//				maxLength: 10,
//				postfix:   "...",
//			},
//			args: args{
//				str: "Короткий текст, Короткий текст, Короткий текст, Короткий текст",
//			},
//			want: "Короткий т...",
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			h := &MaxLengthProcessor{
//				maxLength: tt.fields.maxLength,
//				postfix:   tt.fields.postfix,
//			}
//			if got := h.Truncate(tt.args.str); got != tt.want {
//				t.Errorf("Truncate() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestNewMaxLengthProcessor(t *testing.T) {
//	type args struct {
//		opts []MaxLengthOption
//	}
//	tests := []struct {
//		name string
//		args args
//		want *MaxLengthProcessor
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := NewMaxLengthProcessor(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewMaxLengthProcessor() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
