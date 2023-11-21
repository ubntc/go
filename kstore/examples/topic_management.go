package examples

import (
	"context"
	"errors"
	"log"

	"github.com/ubntc/go/kstore/kstore"
	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/provider/kafkago"
)

func RunTopicManagement(ctx context.Context, tm *kstore.SchemaManager, cfg *config.KeyFile) (result error) {
	log.Println("start demo: RunTopicManagement")

	tbl, err := kstore.NewTableSchema("table1", kstore.Field{Name: "col1", Type: kstore.FieldTypeString})
	if err != nil {
		return err
	}

	group := "demo1"
	c := kafkago.NewClient(cfg, config.DefaultProperties(), config.Group{ID: group})
	db := kstore.NewDatabase(tm, c)

	errch, err := db.StartTableReader(ctx, tbl)
	if err != nil {
		return err
	}
	defer func() {
		result = errors.Join(
			result,  // join the returned error
			<-errch, // with the async error
		)
	}()

	if err := db.CreateOrUpdateTable(ctx, tbl); err != nil {
		return err
	}

	addRows := func() (err error) {
		if err := db.WriteRows(ctx, tbl, kstore.GenerateRows(tbl, 10)...); err != nil {
			return err
		}
		log.Println("wrote 10 messages to table:", tbl.Name)
		return nil
	}

	var l Steps
	l.Add("create", func() error { return db.CreateOrUpdateTable(ctx, tbl) })
	l.Add("write rows 1", addRows)
	l.Add("update 1", func() error {
		tbl.Schema = append(tbl.Schema, kstore.Field{Name: "col2", Type: kstore.FieldTypeString})
		return db.CreateOrUpdateTable(ctx, tbl)
	})
	l.Add("write rows 1", addRows)
	l.Add("update 2", func() error {
		tbl.Schema = append(tbl.Schema, kstore.Field{Name: "col3", Type: kstore.FieldTypeString})
		return db.CreateOrUpdateTable(ctx, tbl)
	})
	l.Add("write rows 2", addRows)
	l.Add("delete all", func() error {
		tbl := *tbl
		tbl.Schema = nil
		return errors.Join(
			db.CreateOrUpdateTable(ctx, &tbl),
			tm.DeleteTopic(ctx, tbl.GetTopic()),
			tm.DeleteTopic(ctx, config.DefaultSchemasTopic),
		)
	})

	return l.Run()
}
