package dogapm

import (
	"context"
	"database/sql/driver"
	"errors"
)

type Hooks struct {
	Before func(ctx context.Context, query string, args ...any) (context.Context, error)
	After  func(ctx context.Context, query string, args ...any) (context.Context, error)
	OnError  func(ctx context.Context, err error, query string, args ...any) (error)
}

type Driver struct {
	driver.Driver
	hooks Hooks
}

type Conn struct {
	driver.Conn
	hooks Hooks
}

type Stmt struct {
	driver.Stmt
	hooks Hooks
	query string
}

func callOnError(h Hooks, ctx context.Context, err error, query string, args ...any) error {
	if err == nil {
		return nil
	}
	// driver.ErrSkip is a control signal for database/sql fallback paths, not a real failure.
	if errors.Is(err, driver.ErrSkip) {
		return err
	}
	if h.OnError == nil {
		return err
	}
	if hookErr := h.OnError(ctx, err, query, args...); hookErr != nil {
		return hookErr
	}
	return err
}

func namedToTnterface(args []driver.NamedValue) []any {
	res := make([]any, 0, len(args))
	for _, arg := range args {
		res = append(res, arg.Value)
	}
	return res
}

func (s *Stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if stmt, ok := s.Stmt.(driver.StmtQueryContext); ok {
		_, err := s.hooks.Before(ctx, s.query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}

		rows, err := stmt.QueryContext(ctx, args)
		if err != nil {
			return rows, callOnError(s.hooks, ctx, err, s.query, namedToTnterface(args)...)
		}
		_, err = s.hooks.After(ctx, s.query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}
		return rows, err
	} else {
		panic("not implement StmtQueryContext")
	}
}

func (s *Stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if stmt, ok := s.Stmt.(driver.StmtExecContext); ok {
		_, err := s.hooks.Before(ctx, s.query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}

		res, err := stmt.ExecContext(ctx, args)
		if err != nil {
			return res, callOnError(s.hooks, ctx, err, s.query, namedToTnterface(args)...)
		}
		_, err = s.hooks.After(ctx, s.query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}
		return res, err
	} else {
		panic("not implement StmtExecContext")
	}
}

func(c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if prepare, ok := c.Conn.(driver.ConnPrepareContext); ok {
		stmt, err := prepare.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		return &Stmt{
			Stmt:  stmt,
			hooks: c.hooks,
			query: query,
		}, nil
	} else {
		panic("not implement ConnPrepareContext")
	}
}


func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if exec, ok := c.Conn.(driver.ExecerContext); ok {
		_, err := c.hooks.Before(ctx, query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}

		res, err := exec.ExecContext(ctx, query, args)
		if err != nil {
			return res, callOnError(c.hooks, ctx, err, query, namedToTnterface(args)...)
		}
		_, err = c.hooks.After(ctx, query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}
		return res, err
	} else {
		panic("not implement ExecerContext")
	}
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if queryer, ok := c.Conn.(driver.QueryerContext); ok {
		_, err := c.hooks.Before(ctx, query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}

		rows, err := queryer.QueryContext(ctx, query, args)
		if err != nil {
			return rows, callOnError(c.hooks, ctx, err, query, namedToTnterface(args)...)
		}
		_, err = c.hooks.After(ctx, query, namedToTnterface(args)...)
		if err != nil {
			return nil, err
		}
		return rows, err
	} else {
		panic("not implement QueryerContext")
	}
}




func (drv *Driver) Open(name string) (driver.Conn, error) {
	conn, err := drv.Driver.Open(name)
	if err != nil {
		return nil, err
	}
	return &Conn{
		Conn:  conn,
		hooks: drv.hooks,
	}, nil
}