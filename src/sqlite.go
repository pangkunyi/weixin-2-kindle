package main

import (
	"atom"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const (
	WEIXIN_MP_ACC_DML = `create table weixin_mp_acc(
			id integer not null primary key autoincrement, 
			name text,
			open_id text,
			open_id_ext text)`
	WEIXIN_MP_ARTICLE_DML = `create table weixin_mp_article(
			id integer not null primary key autoincrement, 
			title text,
			identity text,
			url text,
			content text,
			entry_summary text)`
	WEIXIN_MP_ARTICLE_IDX = `create index idx_identity on weixin_mp_article (identity)`
)

func initDb() {
	if _, err := os.Stat(C.WorkDir + "/weixin-2-kindle.db"); os.IsNotExist(err) {
		err := InitDatabase()
		if err != nil {
			panic(err)
		}
	}
}

func InitDatabase() error {
	os.Remove(DB_FILE)
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return err
	}
	defer db.Close()

	sqls := []string{WEIXIN_MP_ACC_DML, WEIXIN_MP_ARTICLE_DML, WEIXIN_MP_ARTICLE_IDX}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveEntry(feedTitle, feedId string, entry *atom.Entry) error {
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`insert into feed(feed_title, feed_id, entry_title, entry_id, entry_summary) values(?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(feedTitle, feedId, entry.Title, entry.Id, entry.Summary)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func GetEntry(feedId, entryId string, entry *atom.Entry) error {
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("select entry_Title, entry_summary from feed where feed_id = ? and entry_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(feedId, entryId)
	entry.Id = entryId
	return row.Scan(&entry.Title, &entry.Summary)
}

func HasEntry(feedId, entryId string) (bool, error) {
	var entry atom.Entry
	err := GetEntry(feedId, entryId, &entry)
	if err == nil { //found entry already in database
		return true, nil
	} else if err != sql.ErrNoRows { //database query error
		return false, err
	}
	return false, nil
}
