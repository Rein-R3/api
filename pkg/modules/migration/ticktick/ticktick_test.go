// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-2021 Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licensee as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licensee for more details.
//
// You should have received a copy of the GNU Affero General Public Licensee
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ticktick

import (
	"testing"
	"time"

	"code.vikunja.io/api/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertTicktickTasksToVikunja(t *testing.T) {
	time1, err := time.Parse(time.RFC3339Nano, "2022-11-18T03:00:00.4770000Z")
	require.NoError(t, err)
	time2, err := time.Parse(time.RFC3339Nano, "2022-12-18T03:00:00.4770000Z")
	require.NoError(t, err)
	time3, err := time.Parse(time.RFC3339Nano, "2022-12-10T03:00:00.4770000Z")
	require.NoError(t, err)
	duration, err := time.ParseDuration("24h")
	require.NoError(t, err)

	tickTickTasks := []*tickTickTask{
		{
			TaskID:    1,
			ParentID:  0,
			ListName:  "List 1",
			Title:     "Test task 1",
			Tags:      []string{"label1", "label2"},
			Content:   "Lorem Ipsum Dolor sit amet",
			StartDate: time1,
			DueDate:   time2,
			Reminder:  duration,
			Repeat:    "FREQ=WEEKLY;INTERVAL=1;UNTIL=20190117T210000Z",
			Status:    "0",
			Order:     -1099511627776,
		},
		{
			TaskID:        2,
			ParentID:      1,
			ListName:      "List 1",
			Title:         "Test task 2",
			Status:        "1",
			CompletedTime: time3,
			Order:         -1099511626,
		},
		{
			TaskID:    3,
			ParentID:  0,
			ListName:  "List 1",
			Title:     "Test task 3",
			Tags:      []string{"label1", "label2", "other label"},
			StartDate: time1,
			DueDate:   time2,
			Reminder:  duration,
			Status:    "0",
			Order:     -109951627776,
		},
		{
			TaskID:   4,
			ParentID: 0,
			ListName: "List 2",
			Title:    "Test task 4",
			Status:   "0",
			Order:    -109951627777,
		},
	}

	vikunjaTasks := convertTickTickToVikunja(tickTickTasks)

	assert.Len(t, vikunjaTasks, 1)
	assert.Len(t, vikunjaTasks[0].Lists, 2)

	assert.Len(t, vikunjaTasks[0].Lists[0].Tasks, 3)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Title, tickTickTasks[0].ListName)

	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].Title, tickTickTasks[0].Title)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].Description, tickTickTasks[0].Content)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].StartDate, tickTickTasks[0].StartDate)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].EndDate, tickTickTasks[0].DueDate)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].DueDate, tickTickTasks[0].DueDate)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].Labels, []*models.Label{
		{Title: "label1"},
		{Title: "label2"},
	})
	//assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].Reminders, tickTickTasks[0].) // TODO
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].Position, tickTickTasks[0].Order)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].Done, false)

	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[1].Title, tickTickTasks[1].Title)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[1].Position, tickTickTasks[1].Order)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[1].Done, true)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[1].DoneAt, tickTickTasks[1].CompletedTime)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[1].RelatedTasks, models.RelatedTaskMap{
		models.RelationKindParenttask: []*models.Task{
			{
				ID: tickTickTasks[1].ParentID,
			},
		},
	})

	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].Title, tickTickTasks[2].Title)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].Description, tickTickTasks[2].Content)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].StartDate, tickTickTasks[2].StartDate)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].EndDate, tickTickTasks[2].DueDate)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].DueDate, tickTickTasks[2].DueDate)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].Labels, []*models.Label{
		{Title: "label1"},
		{Title: "label2"},
		{Title: "other label"},
	})
	//assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[0].Reminders, tickTickTasks[0].) // TODO
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].Position, tickTickTasks[2].Order)
	assert.Equal(t, vikunjaTasks[0].Lists[0].Tasks[2].Done, false)

	assert.Len(t, vikunjaTasks[0].Lists[1].Tasks, 1)
	assert.Equal(t, vikunjaTasks[0].Lists[1].Title, tickTickTasks[3].ListName)

	assert.Equal(t, vikunjaTasks[0].Lists[1].Tasks[0].Title, tickTickTasks[3].Title)
	assert.Equal(t, vikunjaTasks[0].Lists[1].Tasks[0].Position, tickTickTasks[3].Order)
}
