package controller

import (
	"testing"
	"time"
)

func TestController_isInQuietHours(t *testing.T) {
	type fields struct {
		config Config
		status status
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "testEmptyQuietHours",
			fields: fields{
				config: Config{
					QuietHours: [2]int{},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 0, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "testEmptyHoursIn",
			fields: fields{
				config: Config{
					QuietHours: [2]int{10, 12},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 11, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "testEmptyHoursNotIn",
			fields: fields{
				config: Config{
					QuietHours: [2]int{10, 12},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 13, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "testEmptyHoursInInclusive",
			fields: fields{
				config: Config{
					QuietHours: [2]int{10, 12},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 12, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "testEmptyHoursInInclusive",
			fields: fields{
				config: Config{
					QuietHours: [2]int{10, 12},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 10, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "testEmptyHoursCrossDay",
			fields: fields{
				config: Config{
					QuietHours: [2]int{23, 10},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 4, 0, 0, 0, time.UTC),
			},
			want: true,
		},
		{
			name: "testEmptyHoursCrossDayOut",
			fields: fields{
				config: Config{
					QuietHours: [2]int{23, 10},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 11, 0, 0, 0, time.UTC),
			},
			want: false,
		},
		{
			name: "testEmptyHoursCrossDayOut",
			fields: fields{
				config: Config{
					QuietHours: [2]int{23, 10},
				},
			},
			args: args{
				t: time.Date(2021, 0, 0, 22, 0, 0, 0, time.UTC),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Controller{
				config: tt.fields.config,
				status: tt.fields.status,
			}
			if got := c.isInQuietHours(tt.args.t); got != tt.want {
				t.Errorf("isInQuietHours() = %v, want %v", got, tt.want)
			}
		})
	}
}
