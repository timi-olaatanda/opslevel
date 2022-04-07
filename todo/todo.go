package todo

import (
	"errors"
)

// TodoItem represents a task that needs to be completed. Smaller priority values
// indicates higher priority.
type TodoItem struct {
	Priority    int    `json:"priority"`
	Description string `json:"description"`
}

// todoItems is a contiguous group of task. Task are clustered using priority
// e.g. if Task with priority 1-3 exist, one task group is used to represent
type todoItems struct {
	todoos [][]*TodoItem
}

func (g *todoItems) highestPriority() int {
	return g.todoos[0][0].Priority
}

func (g *todoItems) lowestPriority() int {
	return g.todoos[len(g.todoos)-1][0].Priority
}

func (g *todoItems) addTodo(priority int, description string, extend bool) ([]*TodoItem, bool) {
	s, e, done := 0, len(g.todoos)-1, false
	//log.Printf("priority = %d, description = %s, extend = %v\r\n", priority, description, extend)
	for s <= e {
		m := s + ((e - s) / 2)
		if g.todoos[m][0].Priority == priority {
			//log.Println("got here")
			g.todoos[m], done = append(g.todoos[m], &TodoItem{priority, description}), true
			break
		} else if g.todoos[m][0].Priority < priority {
			s = e + 1
		} else {
			e = m - 1
		}
	}
	if !done && extend {
		if g.todoos[0][0].Priority == 1+priority {
			g.todoos, done = append([][]*TodoItem{{&TodoItem{priority, description}}}, g.todoos...), true
		} else if g.todoos[len(g.todoos)-1][0].Priority == priority-1 {
			g.todoos, done = append(g.todoos, []*TodoItem{{priority, description}}), true
		}
	}
	if done {
		return g.todoos[priority-g.highestPriority()], done
	}
	return nil, done
}

var allTodoos []*todoItems

// AddTask adds a todo item with the specified priority and description. Priority must be in [1, 99999999] range,
// description can only contain A-Z (case insensitive), space or full-stop (.).
func AddTask(priority int, description string) ([]*TodoItem, error) {
	if priority < 1 || description == "" {
		return nil, errors.New("priority must be at least one (1) and Description must be non-empty")
	}
	if len(allTodoos) < 1 {
		allTodoos = append(allTodoos, &todoItems{[][]*TodoItem{{&TodoItem{priority, description}}}})
		return allTodoos[0].todoos[0], nil
	}
	s, e, done, next := 0, len(allTodoos)-1, false, -1
	var todoos []*TodoItem

	for s <= e && !done {
		m := s + ((e - s) / 2)
		tg := allTodoos[m]
		if todoos, done = tg.addTodo(priority, description, false); done {
			break
		}

		if priority < tg.highestPriority() {
			e = m - 1
		} else if priority > tg.lowestPriority() {
			if next < 0 || allTodoos[next].lowestPriority() < allTodoos[m].lowestPriority() {
				next = m
			}
			s = m + 1
		}
	}

	if !done {
		if next >= 0 {
			todoos, done = allTodoos[next].addTodo(priority, description, true)
		}
		if !done {
			tg := &todoItems{[][]*TodoItem{{&TodoItem{priority, description}}}}
			if next >= 0 {
				allTodoos = append(allTodoos[:next+1], append([]*todoItems{tg}, allTodoos[next+1:]...)...)
			} else {
				allTodoos = append([]*todoItems{tg}, allTodoos...)
			}
			todoos = tg.todoos[0]
		}
	}

	return todoos, nil
}

// RemoveTasks removes the todoitem(s) with the specified priority. RemoveTask returns the
// removed task, if there is no todoitem with the specified priority, an error is returned.
// RemoveTasks runs in O(lg n) time.
func RemoveTasks(priority int) ([]*TodoItem, error) {
	if priority < 1 || len(allTodoos) < 1 {
		return nil, errors.New("there is no todo item with the specified priority")
	} else if allTodoos[len(allTodoos)-1].lowestPriority() < priority {
		return nil, errors.New("there is no todo item with the specified priority")
	}

	var res []*TodoItem

	s, e := 0, len(allTodoos)-1
	for s <= e && res == nil { // binary search to find the taskgroup to delete
		m := s + ((e - s) / 2)
		tg := allTodoos[m]
		if priority >= tg.highestPriority() && priority <= tg.lowestPriority() {
			st, en := 0, len(tg.todoos)-1
			for st <= en {
				mid := st + ((en - st) / 2)
				if tg.todoos[mid][0].Priority == priority {
					res = tg.todoos[mid]
					v := append(tg.todoos[:mid], tg.todoos[mid+1:]...)
					if len(v) < 1 {
						allTodoos = append(allTodoos[:m], allTodoos[m+1:]...)
					} else {
						tg.todoos = v
					}
					break
				} else if tg.todoos[mid][0].Priority > priority {
					en = mid - 1
				} else {
					st = mid + 1
				}
			}
		} else if priority < allTodoos[m].highestPriority() {
			e = m - 1
		} else {
			s = m + 1
		}
	}

	if res == nil {
		return nil, errors.New("there is no todo item(s) with the specified priority")
	}

	return res, nil
}

// GetMissPriorities returns all missing priorities with the boundary of highest-lowest priority.
func GetMissingPriorities() []int {
	var res []int
	for idx, tg := range allTodoos {
		if idx == 0 {
			for i := 1; i < tg.highestPriority(); i++ {
				res = append(res, i)
			}
		} else {
			for i := allTodoos[idx-1].lowestPriority() + 1; i < tg.highestPriority(); i++ {
				res = append(res, i)
			}
		}
	}
	return res
}

// GetAllTodoItems returns all todo item.
func GetAllTodoItems() []*TodoItem {
	var res []*TodoItem

	for _, todoItemGroup := range allTodoos {
		for _, todoItem := range todoItemGroup.todoos {
			res = append(res, todoItem...)
		}
	}

	return res
}
