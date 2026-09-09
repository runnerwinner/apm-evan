package dogapm

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

func TestMysqlWrapper(t *testing.T) {
	testDriver := &Driver{
		Driver: mysql.MySQLDriver{},
		hooks: Hooks{
			Before: func(ctx context.Context, query string, args ...any) (context.Context, error) {
				t.Logf("Before query: %s, args: %v", query, args)
				return ctx, nil
			},
			After: func(ctx context.Context, query string, args ...any) (context.Context, error) {
				t.Logf("After query: %s, args: %v", query, args)
				return ctx, nil
			},
			OnError: func(ctx context.Context, err error, query string, args ...any) error {
				t.Logf("OnError query: %s, args: %v", query, args)
				return err
			},
		},
	}

	sql.Register("test-driver", testDriver)

	db, err := sql.Open("test-driver", "root:password@tcp(localhost:3307)/ordersvc")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test a simple query
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test_table (id INT PRIMARY KEY, name VARCHAR(50))")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec("INSERT INTO test_table (id, name) VALUES (?, ?)", 1, "test")
	if err != nil {
		t.Fatalf("Failed to insert data: %v", err)
	}

	rows, err := db.Query("SELECT id, name FROM test_table WHERE id = ?", 1)
	if err != nil {
		t.Fatalf("Failed to query data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			t.Fatalf("Failed to scan data: %v", err)
		}
		t.Logf("Got data: %d %s", id, name)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("Rows iteration failed: %v", err)
	}
}

func TestTraceDriver(t *testing.T) {
	// Register the wrapped driver
	Infra.Init(
		InfraEnableApm("127.0.0.1:54317"),
		InfraDbOption("root:password@tcp(localhost:3307)/ordersvc"),		
	)
	var slept int
	if err := Infra.Db.QueryRowContext(context.Background(), "SELECT SLEEP(5)").Scan(&slept); err != nil {
		t.Fatalf("Failed to run slow query: %v", err)
	}
	if slept != 0 {
		t.Fatalf("Unexpected sleep result: %d", slept)
	}
	EndPoint.Close()
}

func TestLongTx(t *testing.T) {
	Infra.Init(
		InfraEnableApm("127.0.0.1:54317"),
		InfraDbOption("root:password@tcp(localhost:3307)/ordersvc"),
	)

	ctx,span := Tracer.Start(context.Background(), "longTx")
	tx, err := Infra.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	time.Sleep(5*time.Second)
	tx.Rollback()
	span.End()
	EndPoint.Close()
}