package main

import (
	"testing"
	"errors"
)

/*
If the email is empty, return errors.New("email can't be empty")
If the password is empty, password can't be empty
If the age is less than 18, return age must be at least AGE years old, where AGE is the actual age.
*/
func userIsEligible(email, password string, age int) error {
	if email == "" {
		return errors.New("email can't be empty")
	}
	if password == "" {
		return errors.New("password can't be empty")
	}
	if age < 18 {
		return errors.New("age must be at least 18 years old")
	}
	return nil
}

func TestUserIsEligible(t *testing.T) {

	// create test data
	var tests = []struct {
		email       string
		password    string
		age         int
		expectedErr error
	}{
		{
			email:       "test@example.com",
			password:    "12345",
			age:         18,
			expectedErr: nil,
		},
		{
			email:       "",
			password:    "12345",
			age:         18,
			expectedErr: errors.New("email can't be empty"),
		},
		{
			email:       "test@example.com",
			password:    "",
			age:         18,
			expectedErr: errors.New("password can't be empty"),
			// expectedErr: nil,
		},
		{
			email:       "test@example.com",
			password:    "12345",
			age:         16,
			expectedErr: errors.New("age must be at least 18 years old"),
		},
	}

	// run all tests
	for _, tt := range tests {
		err := userIsEligible(tt.email, tt.password, tt.age)
		errString := ""
		expectedErrString := ""
		if err != nil {
			errString = err.Error()
		}
		if tt.expectedErr != nil {
			expectedErrString = tt.expectedErr.Error()
		}
		if errString != expectedErrString {
			t.Errorf("got %s, want %s", errString, expectedErrString)
		}
	}
}