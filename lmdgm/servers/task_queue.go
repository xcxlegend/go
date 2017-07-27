package servers

var taskQueue *TaskQueue

func init() {
	// taskQueue = new(TaskQueue)
	// taskQueue.Init()
}

func GetTaskQuestInst() *TaskQueue {
	return taskQueue
}

type TaskQueue struct {
	// Schedules map[string]*Schedule
}

func (this *TaskQueue) Init() {
	// this.Schedules = make(map[string]*Schedule)
	// go this.manange()
}

func (this *TaskQueue) manange() {
	for {
		select {}
	}
}

func (this *TaskQueue) Start(task *Task) {

}

type Task interface {
	Code() string
	Start()
}

// type SyncTask struct {
// }

// func (this *SyncTask) Code() string {
// 	return "main"
// }

// func (this *SyncTask) Start() {

// }

// type Schedule struct {
// }
