package cron

type Route struct {
}

func (r *Route) Run(cron *Cron) {
	////每3秒执行1次
	// cron.Schedule("*/15 * * * * *").RunFunc(Test1) //每隔3秒运行
	////不设置Schedule会复用上一条的Schedule设置
	//cron.RunFunc(Test2)
	//使用job方式运行，只需实现Run方法
	cron.Schedule("*/10 * * * * *").RunJob(Job{}) //每隔一分钟运行
}
