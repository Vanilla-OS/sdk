package goodies

import (
	"container/heap"
	"fmt"
)

// ErrorHandler is the protocol for handling errors, use this to implement
// custom error handling logic.
type ErrorHandler interface {
	// HandleError triggers the error handling logic
	HandleError(args ...interface{}) error
}

// NoErrorHandler is a default error handler that does nothing, useful when
// you don't need to handle errors.
type NoErrorHandler struct{}

func (n *NoErrorHandler) HandleError(args ...interface{}) error {
	return nil
}

// ErrorHandlerFn is an error handler that implements the ErrorHandler
// protocol, use this to create custom error handling logic.
type ErrorHandlerFn func(args ...interface{}) error

func (f ErrorHandlerFn) HandleError(args ...interface{}) error {
	return f(args...)
}

// CleanupQueue is a priority queue that runs cleanup tasks in order of
// priority and uses the defined error handler to handle errors.
type CleanupQueue struct {
	tasks cleanupHeap
}

// NewCleanupQueue creates a new cleanup queue.
func NewCleanupQueue() *CleanupQueue {
	q := &CleanupQueue{}
	heap.Init(&q.tasks)
	return q
}

// CleanupTask is a task that defines a cleanup task to be run in the cleanup
// queue. It has a priority, a task to run, a list of arguments, an error
// handler to handle errors, and a flag to ignore error handler failures.
type CleanupTask struct {
	Priority                  int
	Task                      func(args ...interface{}) error
	Args                      []interface{}
	ErrorHandler              ErrorHandler
	IgnoreErrorHandlerFailure bool // Flag to ignore error handler failure
	index                     int
}

// Add adds a new task to the cleanup queue.
func (q *CleanupQueue) Add(task func(args ...interface{}) error, args []interface{}, priority int, errorHandler ErrorHandler, ignoreErrorHandlerFailure bool) {
	cleanupTask := &CleanupTask{
		Task:                      task,
		Args:                      args,
		Priority:                  priority,
		ErrorHandler:              errorHandler,
		IgnoreErrorHandlerFailure: ignoreErrorHandlerFailure,
	}
	heap.Push(&q.tasks, cleanupTask)
}

// Run runs the cleanup queue and executes all the tasks in order of priority.
// It returns an error if any of the tasks encounter an error, including the
// error handler function, unless the task is marked to ignore error handler
// failures.
func (q *CleanupQueue) Run() error {
	for q.tasks.Len() > 0 {
		task := heap.Pop(&q.tasks).(*CleanupTask)
		err := task.Task(task.Args...)
		if err != nil {
			errHandle := task.ErrorHandler.HandleError(task.Args...)
			if errHandle != nil {
				if !task.IgnoreErrorHandlerFailure {
					return fmt.Errorf("error handling failed: %v", errHandle)
				}
			}
		}
	}
	return nil
}

type cleanupHeap []*CleanupTask

func (h cleanupHeap) Len() int           { return len(h) }
func (h cleanupHeap) Less(i, j int) bool { return h[i].Priority < h[j].Priority }
func (h cleanupHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *cleanupHeap) Push(x interface{}) {
	n := len(*h)
	task := x.(*CleanupTask)
	task.index = n
	*h = append(*h, task)
}

func (h *cleanupHeap) Pop() interface{} {
	old := *h
	n := len(old)
	task := old[n-1]
	task.index = -1
	*h = old[0 : n-1]
	return task
}
