package main

import (
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
	WEIXIN_MP_ACC_SAVE_SQL     = `insert into weixin_mp_acc(name,open_id,open_id_ext) values(?,?,?)`
	WEIXIN_MP_ACC_LIST_ALL_SQL = `select id, name, open_id, open_id_ext from weixin_mp_acc`
	WEIXIN_MP_ARTICLE_DML      = `create table weixin_mp_article(
			id integer not null primary key autoincrement, 
			title text,
            acc_id integer,
			identity text,
			url text,
			content text)`
	WEIXIN_MP_ARTICLE_IDX            = `create unique index udx_identity on weixin_mp_article (acc_id,identity)`
	WEIXIN_MP_ARTICLE_SAVE_SQL       = `insert into weixin_mp_article(title, acc_id, identity, url, content) values(?,?,?,?,?)`
	WEIXIN_MP_ARTICLE_GET_BY_UDX_SQL = `select id, title, acc_id, identity, url, content from weixin_mp_article where acc_id=? and identity=?`
)

var (
	DB_FILE string
	db      *sql.DB
)

func initDb() {
	DB_FILE = C.WorkDir + "/weixin-2-kindle.db"
	needInitDb := false
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		needInitDb = true
	}
	var err error
	if db, err = sql.Open("sqlite3", DB_FILE); err != nil {
		panic(err)
	}

	if needInitDb {
		sqls := []string{WEIXIN_MP_ACC_DML, WEIXIN_MP_ARTICLE_DML, WEIXIN_MP_ARTICLE_IDX}
		for _, sql := range sqls {
			if _, err = db.Exec(sql); err != nil {
				panic(err)
			}
		}
		for _, acc := range C.WeixinMpAccs {
			SaveWeixinMpAcc(acc.Name, acc.OpenId, acc.OpenIdExt)
		}
	}
}

func CloseDb() {
	db.Close()
}

func QueryDb(fn func(*sql.Rows) error, query string, args ...interface{}) (err error) {
	var rows *sql.Rows
	if rows, err = db.Query(query, args...); err != nil {
		return
	}
	defer rows.Close()
	if err = fn(rows); err != nil {
		return
	}
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func SaveWeixinMpAcc(name, openId, openIdExt string) error {
	_, err := db.Exec(WEIXIN_MP_ACC_SAVE_SQL, name, openId, openIdExt)
	return err
}

func ListAllWeixinMpAcc() (accs []*WeixinMpAcc, err error) {
	err = QueryDb(func(rows *sql.Rows) error {
		for rows.Next() {
			var id int
			var name, openId, openIdExt sql.NullString
			if err := rows.Scan(&id, &name, &openId, &openIdExt); err != nil {
				return err
			}
			accs = append(accs, &WeixinMpAcc{Id: id, Name: name.String, OpenId: openId.String, OpenIdExt: openIdExt.String})
		}
		return nil
	}, WEIXIN_MP_ACC_LIST_ALL_SQL)
	return
}

func SaveWeixinMpArticle(title string, accId int, identity, url, content string) error {
	_, err := db.Exec(WEIXIN_MP_ARTICLE_SAVE_SQL, title, accId, identity, url, content)
	return err
}

func GetOneWeixinMpArticle(accId int, identity string) (article *WeixinMpArticle, err error) {
	err = QueryDb(func(rows *sql.Rows) error {
		for rows.Next() {
			var id, accId int
			var title, identity, url, content sql.NullString
			if err := rows.Scan(&id, &title, &accId, &identity, &url, &content); err != nil {
				return err
			}
			article = &WeixinMpArticle{Id: id, AccId: accId, Title: title.String, Identity: identity.String, Url: url.String, Content: content.String}
			break
		}
		return nil
	}, WEIXIN_MP_ARTICLE_GET_BY_UDX_SQL, accId, identity)
	return
}
