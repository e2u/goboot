package cache

import (
	"time"

	redis "gopkg.in/redis.v4"
)

// Redis 用 Redis 作为同步锁使用
type Redis struct {
	*redis.Client
}

func NewRedis(addr string, db int) *Redis {
	return &Redis{
		Client: redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   db,
		}),
	}
}

// Present 如果指定的 key 存在,返回 true
func (r *Redis) Present(key string) bool {
	var v string
	r.Get(key).Scan(&v)
	if v != "" {
		return true
	}
	return false
}

// Exist ,SetNX 返回值相反的别名操作, 如果指定的 key 不存在,则根据 ttl 设置这个 key=value 返回 false,否则如果 key 已经存在返回 true
func (r *Redis) Exist(key string, value interface{}, ttl time.Duration) bool {
	return !r.SetNX(key, value, ttl).Val()
}

/*
GetSet 尝试从 redis 读取 key 的值,如果找到缓存值,则从 getValFunc 返回,否则用 setValFunc 方法中传入的值做 key 的缓存
范例:
	s3m := hutil.S3MobileInfo{}
	App.Redis.GetSet(key, func() (string, time.Duration,error) {
		if s3mi, err := hutil.GetS3MobileInfo(request.MobileNo); err == nil {
			return s3mi.String(), 0,nil
		}
		s3m = s3mi
		return "", 0,nil
	}, func(val string)error {
		return json.Unmarshal([]byte(val), &s3m)
	})


setValFunc()(string,time.Duration,err)
返回:
	string 要存入 redis 的 value
	time.Duration 這個緩存的 TTL
	err 存儲過程是否有錯誤


getValFunc(string) error
入參:
	string 從緩存中得到的 value
返回:
	error 反序列化的錯誤

*/
func (s *Redis) GetSet(key string, setValFunc func() (string, time.Duration, error), getValFunc func(string) error) error {
	if !s.Present(key) {
		val, ttl, err := setValFunc()
		if err != nil {
			return err
		}
		return s.Set(key, val, ttl).Err()
	}
	var v string
	if err := s.Get(key).Scan(&v); err != nil {
		return err
	}
	return getValFunc(v)
}
