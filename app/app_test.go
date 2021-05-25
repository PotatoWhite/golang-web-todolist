package app

import (
	"encoding/json"
	"fmt"
	"github.com/potatowhite/web/golang-todolist/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
)

func TestTodos(t *testing.T) {
	dbPath := "./unittest.db"

	os.Remove(dbPath)
	assert := assert.New(t)
	ah := MakeHandler(dbPath)
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	// ready to test
	res01, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"test todo 01"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res01.StatusCode)
	todo01 := new(model.Todo)
	err = json.NewDecoder(res01.Body).Decode(todo01)
	assert.NoError(err)
	assert.Equal(false, todo01.Completed)
	assert.Equal(1, todo01.Id)
	assert.Equal("test todo 01", todo01.Name)

	res02, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"test todo 02"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res02.StatusCode)
	todo02 := new(model.Todo)
	err = json.NewDecoder(res02.Body).Decode(todo02)
	assert.NoError(err)
	assert.Equal(false, todo01.Completed)
	assert.Equal(1, todo01.Id)

	assert.Equal(false, todo02.Completed)
	assert.Equal(2, todo02.Id)
	assert.Equal("test todo 02", todo02.Name)

	// test getall
	all, err := http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, all.StatusCode)

	todos := []model.Todo{}
	err = json.NewDecoder(all.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))

	for _, item := range todos {
		if item.Id == todo01.Id {
			assert.Equal(todo01.Name, item.Name)
			assert.Equal(todo01.Completed, item.Completed)
		} else if item.Id == todo02.Id {
			assert.Equal(todo02.Name, item.Name)
			assert.Equal(todo02.Completed, item.Completed)
		} else {
			assert.Error(fmt.Errorf("testID should be id"))
		}
	}

	// update test
	update01, err := http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(todo01.Id) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, update01.StatusCode)
	update02, err := http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(todo02.Id) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, update02.StatusCode)

	all, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, all.StatusCode)

	err = json.NewDecoder(all.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))

	for _, item := range todos {
		if item.Id == todo01.Id {
			assert.Equal(todo01.Name, item.Name)
			assert.Equal(true, item.Completed)
		} else if item.Id == todo02.Id {
			assert.Equal(todo02.Name, item.Name)
			assert.Equal(true, item.Completed)
		} else {
			assert.Error(fmt.Errorf("testID should be id"))
		}
	}

	// test delete 01
	request, err := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(todo01.Id), nil)
	response, err := http.DefaultClient.Do(request)
	assert.NoError(err)
	assert.Equal(http.StatusNoContent, response.StatusCode)

	all, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, all.StatusCode)

	err = json.NewDecoder(all.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(1, len(todos))
	assert.Equal(2, todos[0].Id)
}
