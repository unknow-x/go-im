// Package seeder
/**
  @author:kk
  @data:2021/11/20
  @note
**/
package main

import (
	"crypto/rand"
	"fmt"
	"im_app/config"
	"im_app/internal/app"
	user2 "im_app/internal/app/http/models/user"
	"im_app/internal/pkg/model"
	"im_app/pkg/helpler"
	"log"
	"math/big"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	config.Initialize()

	//设置池
	app.SetupPool()
	wg.Add(7)
	go install(6205, 10000)
	go install(10001, 20000)
	go install(20001, 30000)
	go install(30001, 40000)
	go install(40001, 50000)
	go install(50001, 60000)
	go install(60001, 70000)
	wg.Wait() //阻塞直到所有任务完成

	fmt.Println("over")
}

func install(start int, end int) {
	createTime := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

	for i := start; i < end; i++ {
		name := fmt.Sprintf("测试%d", i)
		age := randomAge()
		user := user2.Users{ID: int64(i),
			Name:          name,
			Avatar:        "https://cdn.learnku.com/uploads/avatars/27407_1531668878.png!/both/100x100",
			Password:      helpler.HashAndSalt("123456"),
			CreatedAt:     createTime,
			Sex:           1,
			Status:        0,
			ClientType:    1,
			Age:           age,
			LastLoginTime: createTime,
		}
		model.DB.Create(&user)
	}
}

func randomAge() int {
	max := big.NewInt(100)
	i, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Fatal("rand:", err)
	}
	return i.BitLen()
}
