package ncm

import (
	"reflect"
	"testing"
)

func TestNcm_Process(t *testing.T) {
	type fields struct {
		Path    string
		OutPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "测试文件转换",
			fields: fields{
				Path:    "/Users/Apple/Desktop/12.ncm",
				OutPath: "/Users/Apple/Desktop/",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := &Ncm{
				Path:    tt.fields.Path,
				OutPath: tt.fields.OutPath,
			}
			if err := receiver.Process(); (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNcm_outPathName(t *testing.T) {
	type fields struct {
		Path    string
		OutPath string
	}
	type args struct {
		meta *MetaInfo
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "测试文件名称",
			fields: fields{
				Path:    "/User/Apple/Desktop/12.ncm",
				OutPath: "/User/Apple/Desktop/",
			},
			args: args{
				meta: &MetaInfo{
					MusicName: "musicName",
					Format:    "mp3",
				},
				name: "12.ncm",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := &Ncm{
				Path:    tt.fields.Path,
				OutPath: tt.fields.OutPath,
			}
			if got := receiver.outPathName(tt.args.meta, tt.args.name); got != tt.want {
				t.Errorf("outPathName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNcm_ParseMateInfo(t *testing.T) {
	type fields struct {
		Path    string
		OutPath string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *MetaInfo
		wantErr bool
	}{
		{
			name: "测试文件解析",
			fields: fields{
				Path:    "/Users/Apple/Application/github/zt/file/12.ncm",
				OutPath: "/Users/Apple/Desktop/",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := &Ncm{
				Path:    tt.fields.Path,
				OutPath: tt.fields.OutPath,
			}
			got, err := receiver.ParseMateInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMateInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMateInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
