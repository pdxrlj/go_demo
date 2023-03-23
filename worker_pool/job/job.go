package job

type Job interface {
	RunTask()
}

type Chan chan Job
