package processor

import (
	"reflect"
	"rss-feed/internal/domain/rss"
	"testing"
)

func TestSizeOfProcessor_Name(t *testing.T) {
	type fields struct {
		maxSize int
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "max size = 10",
			fields: fields{
				maxSize: 10,
			},
			want: "size-of-10",
		},
		{
			name: "max size - negative",
			fields: fields{
				maxSize: -8,
			},
			want: "size-of-8",
		},
		{
			name: "max size greater than 64",
			fields: fields{
				maxSize: 128,
			},
			want: "size-of-64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSizeOfProcessor(tt.fields.maxSize)
			if got := s.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSizeOfProcessor_Process(t *testing.T) {
	type fields struct {
		maxSize int
	}

	type args struct {
		items []*rss.Item
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*rss.Item
		wantErr bool
	}{
		{
			name: "max size = 1",
			fields: fields{
				maxSize: 1,
			},
			args: args{
				items: []*rss.Item{
					{}, {}, {}, {}, {},
				},
			},
			want: []*rss.Item{
				{},
			},
			wantErr: false,
		},
		{
			name: "max size = -1",
			fields: fields{
				maxSize: -1,
			},
			args: args{
				items: []*rss.Item{
					{}, {}, {}, {}, {}, {}, {}, {},
				},
			},
			want: []*rss.Item{
				{}, {}, {}, {}, {}, {}, {}, {},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSizeOfProcessor(tt.fields.maxSize)

			got, err := s.Process(tt.args.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() got = %v, want %v", got, tt.want)
			}
		})
	}
}
