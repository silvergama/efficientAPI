package domain

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

var created_at = time.Now()

func TestMessageRepo_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)

	test := []struct {
		name    string
		s       messageRepoInterface
		msgId   int64
		mock    func()
		want    *Message
		wantErr bool
	}{
		{
			// When everything works as expected
			name:  "OK",
			s:     s,
			msgId: 1,
			mock: func() {
				// We added  one row
				rows := sqlmock.NewRows([]string{
					"Id",
					"Title",
					"Body",
					"CreatedAt",
				}).AddRow(
					1,
					"title",
					"body",
					created_at,
				)
				mock.ExpectPrepare("SELECT (.+) FROM messages").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			want: &Message{
				ID:        1,
				Title:     "title",
				Body:      "body",
				CreatedAt: created_at,
			},
		},
		{
			// When the role tried to access is not found
			name:  "Not Found",
			s:     s,
			msgId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"Id",
					"Title",
					"Body",
					"CreatedAt",
				}) // observe that we didn't add any role  here
				mock.ExpectPrepare("SELECT (.+) FROM messages").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			// When invalid statement is provided, ie the SQL syntax is wrong(in this case, we provided a wrong database)
			name:  "Invalid Prepare",
			s:     s,
			msgId: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{
					"Id",
					"Title",
					"Body",
					"CreatedAt",
				}).AddRow(1, "title", "body", created_at)
				mock.ExpectPrepare("SELECT (.+) FROM wrong_table").ExpectQuery().WithArgs(1).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Get(tt.msgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get(%d) error new = %v, wantErr %v", tt.msgId, err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestMessageRepo_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when opening a stub database", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)
	tm := time.Now()

	tests := []struct {
		name    string
		s       messageRepoInterface
		request *Message
		mock    func()
		want    *Message
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			request: &Message{
				Title:     "title",
				Body:      "body",
				CreatedAt: tm,
			},
			mock: func() {
				mock.ExpectPrepare("INSERT INTO messages").ExpectExec().WithArgs("title", "body", tm).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want: &Message{
				ID:        1,
				Title:     "title",
				Body:      "body",
				CreatedAt: tm,
			},
		},
		{
			name: "Empty title",
			s:    s,
			request: &Message{
				Title:     "",
				Body:      "body",
				CreatedAt: tm,
			},
			mock: func() {
				mock.ExpectPrepare("INSERT INTO messges").ExpectExec().WithArgs("title", "body", tm).WillReturnError(errors.New("empty title"))
			},
			wantErr: true,
		},
		{
			name: "Empty body",
			s:    s,
			request: &Message{
				Title:     "title",
				Body:      "",
				CreatedAt: tm,
			},
			mock: func() {
				mock.ExpectPrepare("INSERT INTO messages").ExpectQuery().WithArgs("title", "body", tm).WillReturnError(errors.New("empty body"))
			},
			wantErr: true,
		},
		{
			name: "Invalid SQL query",
			s:    s,
			request: &Message{
				Title:     "title",
				Body:      "body",
				CreatedAt: tm,
			},
			mock: func() {
				mock.ExpectPrepare("INSERT INTO wrong_table").ExpectQuery().WithArgs("title", "body", tm).WillReturnError(errors.New("invalid sql query"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Create(tt.request)
			if (err != nil) != tt.wantErr {
				fmt.Println("this is error message: ", err.Message())
				t.Errorf("Create() error %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageRepo_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %v was not expected when opening a stab database", err)
	}
	defer db.Close()
	s := NewMessageRepository(db)

	tests := []struct {
		name    string
		s       messageRepoInterface
		request *Message
		mock    func()
		want    *Message
		wantErr bool
	}{
		{
			name: "OK",
			s:    s,
			request: &Message{
				ID:    1,
				Title: "update title",
				Body:  "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "update body", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			want: &Message{
				ID:    1,
				Title: "update title",
				Body:  "update body",
			},
		},
		{
			name: "Invalid SQL Query",
			s:    s,
			request: &Message{
				ID:    1,
				Title: "update title",
				Body:  "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATER messages").ExpectExec().WithArgs("update title", "update body", 1).WillReturnError(errors.New("error in sql query statements"))
			},
			wantErr: true,
		},
		{
			name: "Invalid query ID",
			s:    s,
			request: &Message{
				ID:    0,
				Title: "update title",
				Body:  "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "update body", 0).WillReturnError(errors.New("invalid update id"))
			},
			wantErr: true,
		},
		{
			name: "Empty title",
			s:    s,
			request: &Message{
				ID:    1,
				Title: "",
				Body:  "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("", "update body", 1).WillReturnError(errors.New("Please enter a valid title"))
			},
			wantErr: true,
		},
		{
			name: "Empty body",
			s:    s,
			request: &Message{
				ID:    1,
				Title: "update title",
				Body:  "",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "", 1).WillReturnError(errors.New("Please enter a valid body"))
			},
			wantErr: true,
		},
		{
			name: "Update failed",
			s:    s,
			request: &Message{
				ID:    1,
				Title: "update title",
				Body:  "update body",
			},
			mock: func() {
				mock.ExpectPrepare("UPDATE messages").ExpectExec().WithArgs("update title", "update body", 1).WillReturnResult(sqlmock.NewErrorResult(errors.New("failed update")))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := tt.s.Update(tt.request)
			if (err != nil) != tt.wantErr {
				fmt.Println("this is an error message: ", err.Message())
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}

}
