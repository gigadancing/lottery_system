package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// 参与抽奖的用户
var userLists []string

// 同步锁
var mu sync.Mutex

type lotteryController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

func main() {
	app := newApp()
	// userLists = make([]string, 0)
	userLists = []string{}
	mu = sync.Mutex{}
	app.Run(iris.Addr(":8080"))
}

func (c *lotteryController) Get() string {
	count := len(userLists)
	return fmt.Sprintf("total users: %d\n", count)
}

// POST http://localhost:8080/import
// 导入用户
func (c *lotteryController) PostImport() string {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	mu.Lock()
	defer mu.Unlock()
	count1 := len(userLists)
	for _, u := range users {
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			userLists = append(userLists, u)
		}
	}
	count2 := len(userLists)
	return fmt.Sprintf("total users: %d, success imported users: %d\n", count2, count2-count1)
}

// GET http://localhost:8080/lucky
// 抽奖
func (c *lotteryController) GetLucky() string {
	mu.Lock()
	defer mu.Unlock()
	count := len(userLists)
	if count > 1 {
		// 随机数种子
		seed := time.Now().UnixNano()
		// 产生[0, count)的随机数
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count))
		user := userLists[index]
		// 将中奖的用户删除
		userLists = append(userLists[0:index], userLists[index+1:]...)
		return fmt.Sprintf("the lucky: %s, remaining users: %d\n", user, count-1)
	} else if count == 1 {
		user := userLists[0]
		return fmt.Sprintf("the lucky: %s, remaining users: %d\n", user, count-1)
	} else {
		return fmt.Sprintf("no users \n")
	}
}
