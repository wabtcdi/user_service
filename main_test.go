package main

import (
	"errors"
	"testing"
)

type fakeStarter struct {
	start func() error
}

func (f *fakeStarter) Start() error {
	return f.start()
}

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		configName  string
		openFunc    func(string, string) (*sql.DB, error)
		starter     cmd.Starter
		expectError bool
	}{
		{
			name:       "successful run",
			configName: "local",
			openFunc: func(driverName, dataSourceName string) (*sql.DB, error) {
				return nil, nil
			},
			starter: &fakeStarter{start: func() error {
				return nil
			}},
			expectError: false,
		},
		{
			name:       "init error",
			configName: "local",
			openFunc: func(driverName, dataSourceName string) (*sql.DB, error) {
				return nil, errors.New("init error")
			},
			starter:     &fakeStarter{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.configName, tt.openFunc, tt.starter)
			if (err != nil) != tt.expectError {
				t.Errorf("run() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
