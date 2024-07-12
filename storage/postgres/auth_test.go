package postgres

import (
	pb "auth-service/generated/auth_service"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	user := NewUserRepo(db)

	resp, err := user.CreateUser(&pb.RegisterRequest{
		Username: "diyorbeknematov",
		Password: "qwerty2004",
		Email:    "diyorbeknematov@gmail.com",
		FullName: "Diyorbek Ne'matov",
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := &pb.RegisterResponse{
		Message: "User created successfully",
	}

	assert.NoError(t, err)

	assert.Equal(t, expectedResponse.Message, resp.Message)
}

func TestLogin(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	user := NewUserRepo(db)

	resp, err := user.GetByEmail("diyorbeknematov@gmail.com")
	if err != nil  {
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatal("user not found")
		}
		t.Fatal(err)
	}

	expectedResponse := &pb.LoginResponse{
		UserId:   "fc27aae7-e777-45f1-9431-f00c31dfdea0",
		Username: "diyorbeknematov",
		Email:    "diyorbeknematov@gmail.com",
	}

	if !reflect.DeepEqual(resp, expectedResponse) {
		t.Errorf("have %v , wont %v", resp, &expectedResponse)
	}
}

func TestLogoutUser(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	user := NewUserRepo(db)

	id := "ee681887-fb2f-4722-ab48-0d83dece505c"
	resp, err := user.LogoutUser(id)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := &pb.LogoutResponse{
		Message: "user deleted successully",
	}

	if !reflect.DeepEqual(resp, expectedResponse) {
		t.Errorf("have %v, wont %v", resp, expectedResponse)
	}
}

func TestGetUserProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	user := NewUserRepo(db)

	username := "diyorbeknematov"
	resp, err := user.GetUserProfile(username)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := &pb.GetUserProfileResponse{
		Fullname: "Diyorbek Ne'matov",
		Username: "diyorbeknematov",
		DateOfBirth: "2004-11-20T00:00:00Z",
	}

	if !reflect.DeepEqual(resp, expectedResponse) {
		t.Errorf("have %v, wont %v", resp, expectedResponse)
	}
}

func TestUpdateUserProfile(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	user := NewUserRepo(db)
	resp, err := user.UpdateUserProfile(&pb.UpdateUserProfileRequest{
		UserId:      "fc27aae7-e777-45f1-9431-f00c31dfdea0",
		Username:    "diyorbeknematov",
		FullName:    "Diyorbek Ne'matov",
		DateOfBirth: "2004-11-20",
		PhoneNumber: "+998939955726",
		Address:     "",
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := &pb.UpdateUserProfileResponse{
		Message: "User updated successfully",
	}

	if !reflect.DeepEqual(resp, expectedResponse) {
		t.Errorf("have %v wont %v", resp, expectedResponse)
	}
}