package cron

import (
	"gen/log"
)

func Test1() {
	log.Logger.Infof("Run by func Test1!!!!!")

}

func Test2() {
	log.Logger.Infof("Run by func Test2!!!!!")
}

type Job struct {
}

func (j Job) Run() {
	log.Logger.Infof("这是一个定时任务------->Run by job!!!!!")
	log.Logger.Infof("开始发送钉钉消息")
}
