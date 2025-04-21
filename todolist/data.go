package todolist

import "github.com/charmbracelet/bubbles/list"

// Provides the mock data to fill the kanban board

func (b *Board) initLists() {
	b.cols = []column{
		newColumn(Todo),
		newColumn(InProgress),
		newColumn(Done),
	}
	// Init To Do
	b.cols[Todo].list.Title = "To Do"
	b.cols[Todo].list.SetItems([]list.Item{
		Task{status: Todo, title: "buy milk", description: "strawberry milk"},
		Task{status: Todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
		Task{status: Todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
	})
	// Init in progress
	b.cols[InProgress].list.Title = "In Progress"
	b.cols[InProgress].list.SetItems([]list.Item{
		Task{status: InProgress, title: "write code", description: "don't worry, it's Go"},
	})
	// Init Done
	b.cols[Done].list.Title = "Done"
	b.cols[Done].list.SetItems([]list.Item{
		Task{status: Done, title: "stay cool", description: "as a cucumber"},
	})
}
