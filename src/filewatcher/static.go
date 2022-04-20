package filewatcher

var DefaultFileWatcher = New()

func Add(file string) error {
	return DefaultFileWatcher.Add(file)
}

func Start() error {
	return DefaultFileWatcher.Start()
}

func Stop() error {
	return DefaultFileWatcher.Stop()
}
