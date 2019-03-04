package cache

import (
	"reflect"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	type args struct {
		defaultExpiration time.Duration
		cleanupInterval   time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Cache
	}{
		{
			name:"Create cache with not nullable default expiration and not nullable cleanup interval",
			args:args{
				defaultExpiration:1800,
				cleanupInterval:1000,
			},
			want:NewCache(1800, 1000),
		},
		{
			name:"Create cache with nullable default expiration and cleanup interval",
			args:args{
				defaultExpiration:0,
				cleanupInterval:0,
			},
			want:NewCache(0, 0),
		},
		{
			name:"Create cache with not nullable default expiration and nullable cleanup interval",
			args:args{
				defaultExpiration:1800,
				cleanupInterval:0,
			},
			want:NewCache(1800, 0),
		},
		{
			name:"Create cache with nullable default expiration and not nullable cleanup interval",
			args:args{
				defaultExpiration:1800,
				cleanupInterval:1000,
			},
			want:NewCache(1800, 1000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCache(tt.args.defaultExpiration, tt.args.cleanupInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Set(t *testing.T) {
	type args struct {
		key             string
		value           interface{}
		durationExpired time.Duration
	}
	tests := []struct {
		name string
		o    *Cache
		args args
	}{
		{
			name:"Set in cache with nullable default expiration and nullable cleanup interval and duration expired > 0",
			o:NewCache(0, 0),
			args:args{
				key:"1",
				value:"testSet1",
				durationExpired:10,
			},
		},
		{
			name:"Set in cache with nullable default expiration and nullable cleanup interval and duration expired == 0",
			o:NewCache(0, 0),
			args:args{
				key:"2",
				value:"testSet2",
				durationExpired:0,
			},
		},
		{
			name:"Set in cache with not nullable default expiration and nullable cleanup interval and duration expired > 0",
			o:NewCache(10, 0),
			args:args{
				key:"3",
				value:"testSet3",
				durationExpired:20,
			},
		},
		{
			name:"Set in cache with nullable default expiration and not nullable cleanup interval and duration expired > 0",
			o:NewCache(0, 10),
			args:args{
				key:"4",
				value:"testSet4",
				durationExpired:5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.Set(tt.args.key, tt.args.value, tt.args.durationExpired)
			tvalue, ok := tt.o.items[tt.args.key]
			if !ok {
				t.Errorf("Cache.Set() value %v not set in cache", tt.args.value)
			}
			if !reflect.DeepEqual(tvalue.Value, tt.args.value) {
				t.Errorf("Cache.Get() got1 = %v, want %v", tvalue.Value, tt.args.value)
			}
			var tExpiration int64
			if tt.args.durationExpired > 0 {
				tExpiration = time.Now().Add(tt.args.durationExpired).Unix()
			}
			if !reflect.DeepEqual(tvalue.Expiration, tExpiration) {
				t.Errorf("Cache.Get() got1 = %v, want %v", tvalue.Expiration, tExpiration)
			}
		})
	}
}

func TestCache_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		o     *Cache
		args  args
		want  interface{}
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.o.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCache_Del(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		o    *Cache
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.Del(tt.args.key)
		})
	}
}

func TestCache_Exist(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		o    *Cache
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Exist(tt.args.key); got != tt.want {
				t.Errorf("Cache.Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Count(t *testing.T) {
	tests := []struct {
		name string
		o    *Cache
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Count(); got != tt.want {
				t.Errorf("Cache.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_Rename(t *testing.T) {
	type args struct {
		oldKey string
		newKey string
	}
	tests := []struct {
		name    string
		o       *Cache
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Rename(tt.args.oldKey, tt.args.newKey); (err != nil) != tt.wantErr {
				t.Errorf("Cache.Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCache_IsCacheExpired(t *testing.T) {
	tests := []struct {
		name string
		o    *Cache
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.IsCacheExpired(); got != tt.want {
				t.Errorf("Cache.IsCacheExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache_FlushAll(t *testing.T) {
	tests := []struct {
		name string
		o    *Cache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.FlushAll()
		})
	}
}

func TestCache_runCleanExpiredItems(t *testing.T) {
	tests := []struct {
		name string
		o    *Cache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.runCleanExpiredItems()
		})
	}
}

func TestCache_clearItems(t *testing.T) {
	tests := []struct {
		name string
		o    *Cache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.o.clearItems()
		})
	}
}
