package logging

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/howeyc/fsnotify"
)

type ReopenWriter struct {
	sync.RWMutex
	filename  string
	fd        *os.File
	exit      chan bool
	closeOnce sync.Once
}

func NewReopenWriter(filename string) (*ReopenWriter, error) {
	w := &ReopenWriter{
		filename: filename,
		exit:     make(chan bool),
	}

	err := w.open()
	if err != nil {
		return nil, err
	}

	w.hupWatch()
	w.fsWatch()

	return w, nil
}

func (w *ReopenWriter) open() error {
	w.Lock()
	defer w.Unlock()

	var newf *os.File
	var err error

	newf, err = os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	oldf := w.fd
	w.fd = newf
	if oldf != nil {
		oldf.Close()
	}

	return nil
}

func (w *ReopenWriter) Write(p []byte) (int, error) {
	w.RLock()
	n, err := w.fd.Write(p)
	w.RUnlock()
	return n, err
}

func (w *ReopenWriter) Close() {
	w.closeOnce.Do(func() {
		close(w.exit)
		w.fd.Close()
	})
}

func (w *ReopenWriter) hupWatch() {
	signalChan := make(chan os.Signal, 16)
	signal.Notify(signalChan, syscall.SIGHUP)

	go func() {
		for {
			select {
			case <-signalChan:
				err := w.open()
				logrus.Infof("HUP received, reopen log %s", w.filename)
				if err != nil {
					logrus.Errorf("Reopen log %s failed: %#s", w.filename, err.Error())
				}
			case <-w.exit:
				return
			}
		}
	}()
}

func (w *ReopenWriter) fsWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Warningf("fsnotify.NewWatcher(): %s", err)
		return
	}

	subscribe := func() {
		if err := watcher.WatchFlags(w.filename, fsnotify.FSN_CREATE|fsnotify.FSN_DELETE|fsnotify.FSN_RENAME); err != nil {
			logrus.Warningf("fsnotify.Watcher.Watch(%s): %s", w.filename, err)
		}
	}

	subscribe()

	go func() {
		defer watcher.Close()

		for {
			select {
			case <-watcher.Event:
				w.open()
				subscribe()

				logrus.Infof("Reopen log %#v by fsnotify event", w.filename)
				if err != nil {
					logrus.Errorf("Reopen log %#v failed: %#s", w.filename, err.Error())
				}

			case <-w.exit:
				return
			}
		}
	}()
}
