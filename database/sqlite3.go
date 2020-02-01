package database

import (
	"database/sql"
	"fmt"

	"github.com/Eilick/mysql-client/common"

	_ "github.com/mattn/go-sqlite3"
)

func GetDb() string {
	return *common.WorkDir + "data.db"
}

func AddArticle(title, content string) (int64, error) {

	tmpDb, err := sql.Open("sqlite3", GetDb())

	if err != nil {
		panic(err)
	}

	stmt, err := tmpDb.Prepare("INSERT INTO markdown(title, content, show_status, create_at, update_at) values(?, ?, ?, ?, ?)")

	if err != nil {
		return 0, err
	}

	nowTime := common.GetNowDateTimeString()
	res, err := stmt.Exec(title, content, 0, nowTime, nowTime)
	stmt.Close()
	tmpDb.Close()
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func UpdateMd(id, title, content string) (bool, error) {

	tmpDb, err := sql.Open("sqlite3", GetDb())

	if err != nil {
		panic(err)
	}

	stmt, err := tmpDb.Prepare("UPDATE markdown SET title=?,content = ?,update_at =? WHERE id= ? ")
	if err != nil {
		return false, err
	}

	nowTime := common.GetNowDateTimeString()
	_, err = stmt.Exec(title, content, nowTime, id)
	stmt.Close()
	tmpDb.Close()
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteMd(id string) (bool, error) {

	tmpDb, err := sql.Open("sqlite3", GetDb())

	if err != nil {
		panic(err)
	}

	stmt, err := tmpDb.Prepare("UPDATE markdown SET show_status=-1,update_at =? WHERE id= ? ")
	if err != nil {
		return false, err
	}

	nowTime := common.GetNowDateTimeString()
	_, err = stmt.Exec(nowTime, id)
	stmt.Close()
	tmpDb.Close()
	if err != nil {
		return false, err
	}

	return true, nil
}

func ArticleList() []map[string]interface{} {
	fmt.Println(GetDb())
	tmpDb, err := sql.Open("sqlite3", GetDb())

	if err != nil {
		panic(err)
	}

	rows, err := tmpDb.Query("SELECT id, title FROM markdown where show_status=0 order by update_at desc")

	if err != nil {
		return []map[string]interface{}{}
	}

	list := []map[string]interface{}{}
	for rows.Next() {
		title := ""
		id := 0
		if err := rows.Scan(&id, &title); err == nil {
			list = append(list, map[string]interface{}{
				"id":    id,
				"title": title,
			})
		}
	}
	rows.Close()
	tmpDb.Close()

	return list
}

func SingleArticle(id string) map[string]interface{} {

	tmpDb, err := sql.Open("sqlite3", GetDb())

	if err != nil {
		panic(err)
	}

	rows, err := tmpDb.Query(fmt.Sprintf("SELECT id,title,content,create_at,update_at FROM markdown where id = %s", id))

	if err != nil {
		return map[string]interface{}{}
	}

	data := map[string]interface{}{}
	for rows.Next() {
		title := ""
		id := 0
		content := ""
		createAt := ""
		updateAt := ""
		if err := rows.Scan(&id, &title, &content, &createAt, &updateAt); err == nil {
			data = map[string]interface{}{
				"id":        id,
				"title":     title,
				"content":   content,
				"update_at": updateAt,
				"create_at": createAt,
			}
			break
		}
	}
	rows.Close()
	tmpDb.Close()

	return data
}
