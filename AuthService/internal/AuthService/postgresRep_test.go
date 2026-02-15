package AuthService

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestPostgresRep_CreateUser(t *testing.T) {
	fixedTime := time.Date(2026, 2, 4, 10, 0, 0, 0, time.UTC)

	db, mock, err := sqlmock.New()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if errClose := db.Close(); errClose != nil {
			logrus.Error("can't close db")
		}
	}()

	gormDB, errGorm := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if errGorm != nil {
		logrus.Fatal(err)
	}

	r := NewAuthPostgresRepo(gormDB)

	type args struct {
		listID int
		item   User
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		expectedUser User
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listID: 1,
				item: User{
					UserName:     "test",
					Email:        "test@test.com",
					PasswordHash: "qwerty",
					Role:         "user",
					CreatedAt:    fixedTime,
					UpdatedAt:    fixedTime,
				},
			},
			expectedUser: User{
				ID:           2,
				UserName:     "test",
				Email:        "test@test.com",
				PasswordHash: "qwerty",
				Role:         "user",
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			mockBehavior: func(args args, expectedId int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "user_name", "email", "password_hash", "role", "created_at", "updated_at"}).
					AddRow(expectedId, args.item.UserName, args.item.Email, args.item.PasswordHash, args.item.Role,
						args.item.CreatedAt, args.item.UpdatedAt)
				mock.ExpectQuery(`^INSERT INTO "users"`).
					WithArgs(args.item.UserName, args.item.Email, args.item.PasswordHash, args.item.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				listID: 1,
				item: User{
					UserName:     "",
					Email:        "",
					PasswordHash: "",
					Role:         "",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
				mock.ExpectQuery(`^INSERT INTO "users"`).
					WithArgs(args.item.UserName, args.item.Email, args.item.PasswordHash, args.item.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.expectedUser.ID)

			got, errCreate := r.CreateUser(testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, errCreate)
			} else {
				assert.NoError(t, errCreate)
				assert.Equal(t, testCase.expectedUser, got)
			}
		})
	}
}

func TestPostgresRep_SaveRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if errClose := db.Close(); errClose != nil {
			logrus.Error("can't close db")
		}
	}()

	gormDB, errGorm := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if errGorm != nil {
		logrus.Fatal(errGorm)
	}

	r := NewAuthPostgresRepo(gormDB)

	type args struct {
		listID int
		item   RefreshToken
	}

	type mockBehavior func(args args, id int)

	testTable := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		id           int
		wantErr      bool
	}{
		{
			name: "OK",
			id:   2,
			args: args{
				listID: 1,
				item: RefreshToken{
					UserID:    1,
					TokenHash: "qwerty",
					ExpiresAt: time.Now().Add(5 * time.Second),
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(`^INSERT INTO "refresh_tokens"`).
					WithArgs(1, "qwerty", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			id:   2,
			args: args{
				listID: 1,
				item:   RefreshToken{},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
				mock.ExpectQuery(`^INSERT INTO "refresh_tokens"`).
					WithArgs(args.item.UserID, args.item.TokenHash, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			errCreate := r.SaveRefreshToken(testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, errCreate)
			} else {
				assert.NoError(t, errCreate)
			}
		})
	}
}

func TestPostgresRep_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if errClose := db.Close(); errClose != nil {
			logrus.Error("can't close db")
		}
	}()

	gormDB, errGorm := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if errGorm != nil {
		logrus.Fatal(errGorm)
	}

	r := NewAuthPostgresRepo(gormDB)

	type args struct {
		userEmail string
	}

	testTable := []struct {
		name    string
		mock    func()
		input   args
		want    User
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_name", "email", "password_hash", "role", "created_at", "updated_at"}).
					AddRow(1, "test", "test@test.com", "qwerty", "user", time.Now(), time.Now())

				mock.ExpectQuery(`.*WHERE.*email.*`).
					WithArgs("test@test.com", 1).WillReturnRows(rows)
			},
			input: args{
				userEmail: "test@test.com",
			},
			want: User{
				ID:           1,
				UserName:     "test",
				Email:        "test@test.com",
				PasswordHash: "qwerty",
				Role:         "user",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(`.*WHERE.*email.*`).
					WithArgs("notfound@test.com", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			input: args{
				userEmail: "notfound@test.com",
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, errGet := r.GetUser(testCase.input.userEmail)
			if testCase.wantErr {
				assert.Error(t, errGet)
			} else {
				assert.NoError(t, errGet)
				assert.Equal(t, testCase.want.ID, got.ID)
				assert.Equal(t, testCase.want.UserName, got.UserName)
				assert.Equal(t, testCase.want.Email, got.Email)
				assert.Equal(t, testCase.want.PasswordHash, got.PasswordHash)
				assert.Equal(t, testCase.want.Role, got.Role)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostgresRep_GetUserByRefreshToken(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		if errClose := db.Close(); errClose != nil {
			logrus.Error("can't close db")
		}
	}()

	gormDB, errGorm := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if errGorm != nil {
		logrus.Fatal(errGorm)
	}

	r := NewAuthPostgresRepo(gormDB)

	type args struct {
		refreshToken string
	}

	testTable := []struct {
		name    string
		mock    func()
		input   args
		want    User
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_name", "email", "password_hash", "role", "created_at", "updated_at"}).
					AddRow(1, "test", "test@test.com", "qwerty", "user", time.Now(), time.Now())

				mock.ExpectQuery(`.*WHERE id = \(select user_id from refresh_tokens where token_hash = .*`).
					WithArgs("refreshToken", 1).WillReturnRows(rows)
			},
			input: args{
				refreshToken: "refreshToken",
			},
			want: User{
				ID:           1,
				UserName:     "test",
				Email:        "test@test.com",
				PasswordHash: "qwerty",
				Role:         "user",
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(`.*WHERE id = \(select user_id from refresh_tokens where token_hash = .*`).
					WithArgs("notFoundRefreshToken", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			input: args{
				refreshToken: "notFoundRefreshToken",
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, errGet := r.GetUserByRefreshToken(testCase.input.refreshToken)
			if testCase.wantErr {
				assert.Error(t, errGet)
			} else {
				assert.NoError(t, errGet)
				assert.Equal(t, testCase.want.ID, got.ID)
				assert.Equal(t, testCase.want.UserName, got.UserName)
				assert.Equal(t, testCase.want.Email, got.Email)
				assert.Equal(t, testCase.want.PasswordHash, got.PasswordHash)
				assert.Equal(t, testCase.want.Role, got.Role)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
