package setting

import model "golang-todo-app-atp/models"

var DatabaseSchema = []interface{}{
	&model.Task{},
	&model.TaskImage{},
}
