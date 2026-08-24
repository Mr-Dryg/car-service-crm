package config

import "testing"

func TestLoadDatabaseURL(t *testing.T) {
	type TestCase struct {
		name, host, port, user, password, db_name, sslmode, want string
		err                                                      error
	}

	tests := []TestCase{
		{
			name:     "default",
			host:     "localhost",
			port:     "5432",
			user:     "db_user",
			password: "db_password",
			db_name:  "car_service",
			sslmode:  "disable",
			want:     "postgres://db_user:db_password@localhost:5432/car_service?sslmode=disable",
			err:      nil,
		},
		{
			name:     "missing host",
			host:     "",
			port:     "5432",
			user:     "db_user",
			password: "db_password",
			db_name:  "car_service",
			sslmode:  "disable",
			want:     "",
			err:      ErrIncompleteDatabaseEnv,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("DB_HOST", test.host)
			t.Setenv("DB_PORT", test.port)
			t.Setenv("DB_USER", test.user)
			t.Setenv("DB_PASSWORD", test.password)
			t.Setenv("DB_NAME", test.db_name)
			t.Setenv("DB_SSLMODE", test.sslmode)

			get, err := Load()

			if err != nil && err != test.err {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && err != test.err {
				t.Errorf("expect err, but err = nil")
			} else if test.want != get.DatabaseURL {
				t.Errorf("expect: %v, get: %v", test.want, get.DatabaseURL)
			}
		})
	}
}
