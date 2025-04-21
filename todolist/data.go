package todolist

import (
	"log"

	persistence "github.com/ReggieReo/todo-elm/persistance"
	"github.com/charmbracelet/bubbles/list"
)

// Provides the mock data to fill the kanban board

// initLists initializes the kanban board columns and loads tasks from the database
func (b *Board) initLists() {
	b.cols = []column{
		newColumn(todo),
		newColumn(inProgress),
		newColumn(done),
	}

	// Set column titles
	b.cols[todo].list.Title = "To Do"
	b.cols[inProgress].list.Title = "In Progress"
	b.cols[done].list.Title = "Done"

	// Load tasks from the database
	b.loadTasks()
}

// loadTasks loads tasks from the database
func (b *Board) loadTasks() {
	// Load todo tasks
	todoTasks, err := b.store.LoadTasks(b.username, persistence.Todo)
	if err != nil {
		log.Printf("Error loading todo tasks: %v", err)
		// Fall back to default tasks if error occurs
		b.loadDefaultTasks()
		return
	}

	// Load in-progress tasks
	inProgressTasks, err := b.store.LoadTasks(b.username, persistence.InProgress)
	if err != nil {
		log.Printf("Error loading in-progress tasks: %v", err)
		b.loadDefaultTasks()
		return
	}

	// Load done tasks
	doneTasks, err := b.store.LoadTasks(b.username, persistence.Done)
	if err != nil {
		log.Printf("Error loading done tasks: %v", err)
		b.loadDefaultTasks()
		return
	}

	// Convert persistence tasks to todolist tasks and add to columns
	var todoItems []list.Item
	for _, t := range todoTasks {
		todoItems = append(todoItems, Task{
			status:      status(t.Status),
			title:       t.Title,
			description: t.Description,
		})
	}
	b.cols[todo].list.SetItems(todoItems)

	var inProgressItems []list.Item
	for _, t := range inProgressTasks {
		inProgressItems = append(inProgressItems, Task{
			status:      status(t.Status),
			title:       t.Title,
			description: t.Description,
		})
	}
	b.cols[inProgress].list.SetItems(inProgressItems)

	var doneItems []list.Item
	for _, t := range doneTasks {
		doneItems = append(doneItems, Task{
			status:      status(t.Status),
			title:       t.Title,
			description: t.Description,
		})
	}
	b.cols[done].list.SetItems(doneItems)
}

// loadDefaultTasks loads default demo tasks if no tasks are found in the database
func (b *Board) loadDefaultTasks() {
	// Init To Do
	b.cols[todo].list.SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "strawberry milk"},
		Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	b.cols[inProgress].list.SetItems([]list.Item{
		Task{status: inProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init done
	b.cols[done].list.SetItems([]list.Item{
		Task{status: done, title: "stay cool", description: "as a cucumber"},
	})
}

// saveTasks saves all tasks to the database
func (b *Board) saveTasks() error {
	// Save todo tasks
	var todoTasks []persistence.Task
	for _, item := range b.cols[todo].list.Items() {
		task := item.(Task)
		todoTasks = append(todoTasks, persistence.Task{
			Status:      persistence.TaskStatus(task.status),
			Title:       task.title,
			Description: task.description,
		})
	}
	if err := b.store.SaveTasks(b.username, persistence.Todo, todoTasks); err != nil {
		return err
	}

	// Save in-progress tasks
	var inProgressTasks []persistence.Task
	for _, item := range b.cols[inProgress].list.Items() {
		task := item.(Task)
		inProgressTasks = append(inProgressTasks, persistence.Task{
			Status:      persistence.TaskStatus(task.status),
			Title:       task.title,
			Description: task.description,
		})
	}
	if err := b.store.SaveTasks(b.username, persistence.InProgress, inProgressTasks); err != nil {
		return err
	}

	// Save done tasks
	var doneTasks []persistence.Task
	for _, item := range b.cols[done].list.Items() {
		task := item.(Task)
		doneTasks = append(doneTasks, persistence.Task{
			Status:      persistence.TaskStatus(task.status),
			Title:       task.title,
			Description: task.description,
		})
	}
	if err := b.store.SaveTasks(b.username, persistence.Done, doneTasks); err != nil {
		return err
	}

	return nil
}
