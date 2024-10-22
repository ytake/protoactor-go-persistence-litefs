package persistencelitefs

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/ytake/protoactor-go-persistence-litefs/testdata"
)

func sqlite() (*sql.DB, error) {
	open, err := sql.Open("sqlite3", "./testdata/sqlite/data/data.db")
	return open, err
}

func TestProvider_PersistEvent(t *testing.T) {
	ctx := context.Background()
	conn, err := sqlite()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	t.Cleanup(func() {
		conn.ExecContext(ctx, "DELETE FROM journals")
		conn.Close()
	})
	provider, _ := New(ctx, 3, NewTable(), conn, nil)
	evt := &testdata.UserCreated{
		UserID:   ulid.Make().String(),
		UserName: "test",
		Email:    "",
	}
	provider.PersistEvent("user", 1, evt)
	var evv *testdata.UserCreated
	provider.GetEvents("user", 1, 4, func(e interface{}) {
		ev, ok := e.(*testdata.UserCreated)
		if !ok {
			t.Error("unexpected type")
		}
		evv = ev
	})
	if !reflect.DeepEqual(evt, evv) {
		t.Errorf("unexpected event %v", evv)
	}

	var evv2 *testdata.UserCreated
	provider.GetEvents("user", 1, 0, func(e interface{}) {
		ev, ok := e.(*testdata.UserCreated)
		if !ok {
			t.Error("unexpected type")
		}
		evv2 = ev
	})
	if !reflect.DeepEqual(evt, evv2) {
		t.Errorf("unexpected event %v", evv2)
	}
}

func TestProvider_PersistSnapshot(t *testing.T) {
	ctx := context.Background()
	conn, err := sqlite()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	t.Cleanup(func() {
		conn.ExecContext(ctx, "DELETE FROM snapshots")
		conn.Close()
	})
	provider, _ := New(ctx, 3, NewTable(), conn, nil)
	evt := &testdata.UserCreated{
		UserID:   ulid.Make().String(),
		UserName: "test",
		Email:    "",
	}
	provider.PersistSnapshot("user", 1, evt)
	snapshot, idx, ok := provider.GetSnapshot("user")
	if !ok {
		t.Error("snapshot not found")
	}
	if idx != 1 {
		t.Errorf("unexpected index %d", idx)
	}
	if !reflect.DeepEqual(snapshot, evt) {
		t.Errorf("unexpected snapshot %v", snapshot)
	}
}
