package xfs

type FileChangedCallback func(filename string)

func WatchChange(filename string, cb FileChangedCallback) {
}

func UnWatchChange(cb FileChangedCallback) {

}
