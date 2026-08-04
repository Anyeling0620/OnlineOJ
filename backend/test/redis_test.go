package test

//import (
//	"context"
//	"github.com/go-redis/redis/v8"
//	"testing"
//	"time"
//)
//
//var ctx = context.Background()
//
//var rdb = redis.NewClient(&redis.Options{
//	Addr:     "",
//	Password: "",
//	DB:       0,
//})
//
//func TestRedisSet(t *testing.T) {
//	err := rdb.Set(ctx, "name", "mmp", time.Second*10).Err()
//	if err != nil {
//		t.Fatal(err)
//	}
//}
//
//func TestRedisGet(t *testing.T) {
//	res, err := rdb.Get(ctx, "name").Result()
//	if err != nil {
//		t.Fatal(err)
//	}
//	println(res)
//}
