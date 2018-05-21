package dot_db

import (
	"database/sql"
	"fmt"
	"dotBankNewsAggregator/dot_temp"
	_ "github.com/go-sql-driver/mysql"
)

type NewItem struct {
	DB *sql.DB
	Note dot_temp.TempNews
}

type MysqlTempNews struct {
	ID string
	Arr dot_temp.TempNews
}

func Insert(post NewItem){
	result, err := post.DB.Exec(
		"insert into dot_news.news (site_url, post_link, post_header, post_description) values (?, ?, ?, ?)",
		post.Note.Site,
		post.Note.Link,
		post.Note.Header,
		post.Note.Description)

	if err != nil{
		panic(err)
	}

	//ID добавленного объекта
	fmt.Println(result.LastInsertId())
	//Количество добавленных строк
	fmt.Println(result.RowsAffected())
}

func Select(table string, db sql.DB) dot_temp.TempSetNews  {
	rows, err := db.Query("select * from news")
	if err != nil{
		panic(err)
	}
	defer rows.Close()
	data := dot_temp.TempSetNews{}

	for rows.Next() {
		obj := MysqlTempNews{}
		err := rows.Scan(&obj.ID, &obj.Arr.Site, &obj.Arr.Link, &obj.Arr.Header, &obj.Arr.Description)
		if err != nil{
			fmt.Println(err)
			continue
		}
		data.SetNews = append(data.SetNews, obj.Arr)
	}
	return data
}
