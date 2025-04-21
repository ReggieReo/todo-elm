package todolist

type status int

const (
	Todo status = iota
	InProgress
	Done
)

func (s status) getNext() status {
	if s == Done {
		return Todo
	}
	return s + 1
}

func (s status) getPrev() status {
	if s == Todo {
		return Done
	}
	return s - 1
}

const margin = 4

type Task struct {
	status      status
	title       string
	description string
}

func NewTask(status status, title, description string) Task {
	return Task{status: status, title: title, description: description}
}

func (t *Task) Next() {
	if t.status == Done {
		t.status = Todo
	} else {
		t.status++
	}
}

// implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
