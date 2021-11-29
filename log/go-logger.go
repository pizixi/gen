package log

import (
	"os"

	go_logger "github.com/pizixi/go-logger"
)

var Logger = go_logger.NewLogger()

func init() {
	Logger = newlogger()

}

func newlogger() *go_logger.Logger {
	path := "logs"

	// 如果不存在logs目录新建logs目录
	if _, err := os.Stat(path); err == nil {
		Logger.Noticef("path exists 1", path)
	} else {
		Logger.Noticef("path not exists ", path)
		err := os.MkdirAll(path, 0711)

		if err != nil {
			Logger.Noticef("Error creating directory")
			// return err
		}
	}
	// 初始化日志，以下控制台的配置和输出到文件的配置可以二选一，可以都要，具体官方文档上写的非常详细，三个例子可以参考

	// 这是控制台的配置======================================================
	Logger.Detach("console")

	consoleConfig := &go_logger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	Logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)
	// 这是控制台的配置======================================================

	// 这是输出到文件的配置==================================================
	fileConfig := &go_logger.FileConfig{
		Filename: "logs/all.log",
		LevelFileName: map[int]string{
			Logger.LoggerLevel("error"): "logs/error.log",
			Logger.LoggerLevel("info"):  "logs/info.log",
			Logger.LoggerLevel("debug"): "logs/debug.log",
			// 以下是自定义
			// Logger.LoggerLevel("notice"): "logs/notice.log",
		},
		MaxSize:    1024 * 1024,
		MaxLine:    100000,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
	// 这是输出到文件的配置==================================================
	// Logger.SetAsync()
	// 主程序结束前需要Flush
	// Logger.Flush()
	return Logger
}
