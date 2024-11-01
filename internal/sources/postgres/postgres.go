// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package postgres

import (
	"context"
	"fmt"

	"github.com/googleapis/genai-toolbox/internal/sources"
	"github.com/jackc/pgx/v5/pgxpool"
)

const SourceKind string = "postgres"

// validate interface
var _ sources.SourceConfig = Config{}

type Config struct {
	Name     string `yaml:"name"`
	Kind     string `yaml:"kind"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (r Config) SourceConfigKind() string {
	return SourceKind
}

func (r Config) Initialize() (sources.Source, error) {
	pool, err := initPostgresConnectionPool(r.Host, r.Port, r.User, r.Password, r.Database)
	if err != nil {
		return nil, fmt.Errorf("Unable to create pool: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Unable to connect successfully: %w", err)
	}

	s := Source{
		Name: r.Name,
		Kind: SourceKind,
		Pool: pool,
	}
	return s, nil
}

var _ sources.Source = Source{}

type Source struct {
	Name string `yaml:"name"`
	Kind string `yaml:"kind"`
	Pool *pgxpool.Pool
}

func (s Source) SourceKind() string {
	return SourceKind
}

func (s *Source) PostgresPool() *pgxpool.Pool {
	return s.Pool
}

func initPostgresConnectionPool(host, port, user, pass, dbname string) (*pgxpool.Pool, error) {
	// urlExample := "postgres:dd//username:password@localhost:5432/database_name"
	i := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, dbname)
	pool, err := pgxpool.New(context.Background(), i)
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %w", err)
	}

	return pool, nil
}
