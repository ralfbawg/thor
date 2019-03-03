package task

type AppInfo struct {
	AppName        string
	AppId          string
	AppKey         string
	LimitSpeed     int
	LimitCount     int
	MaxMessageSize int
}

type task struct {
	hub *TaskHub
}


func VerifyApp(appId string, appKey string) (appName string,token string, err error) {

}
