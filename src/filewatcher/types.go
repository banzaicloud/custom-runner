package filewatcher

import (
	"context"
	"path/filepath"

	"example.com/gocr/src/events"
	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	// eventPipe events.Pipe
	watcher *fsnotify.Watcher
	ctx     context.Context
	cancel  context.CancelFunc
	files   map[string]bool
}

// func New(eventPipe events.Pipe) *FileWatcher {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	return &FileWatcher{eventPipe: eventPipe, ctx: ctx, cancel: cancel, files: make(map[string]bool)}
// }

func New() *FileWatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &FileWatcher{ctx: ctx, cancel: cancel, files: make(map[string]bool)}
}

func (f *FileWatcher) Start() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	f.watcher = watcher
	go f.listen()
	return nil
}

func (f *FileWatcher) Stop() error {
	f.cancel()
	defer f.watcher.Close()
	return nil
}

func (f *FileWatcher) listen() {
	for {
		select {
		case <-f.ctx.Done():
			return
		case event, ok := <-f.watcher.Events:
			if !ok {
				f.Stop()
			}
			// fmt.Println(event, ok)
			if e := f.eventForFile(event); e != nil {
				//f.eventPipe <- e
				events.Add(e)
			}
		case err, ok := <-f.watcher.Errors:
			if !ok {
				f.Stop()
			}
			// fmt.Println(err, ok)
			//f.eventPipe <- events.OnError(err)
			events.Add(events.OnError(err))
		}
	}
}

func (f *FileWatcher) Add(file string) error {
	path := filepath.Dir(file)
	f.files[file] = true
	return f.watcher.Add(path)
}

func (f *FileWatcher) eventForFile(event fsnotify.Event) events.IEvent {
	// file := filepath.Dir(event.Name) + "/" + filepath.Base(event.Name)
	file := event.Name
	if _, ok := f.files[file]; ok {
		switch event.Op {
		case fsnotify.Create:
			return events.OnFileCreate(event.Name)
		case fsnotify.Write:
			return events.OnFileWrite(event.Name)
		case fsnotify.Remove:
			return events.OnFileRemove(event.Name)
		case fsnotify.Rename:
			return events.OnFileRename(event.Name)
		case fsnotify.Chmod:
			return events.OnFileChmod(event.Name)
		default:
			return nil
		}
	}
	return nil
}