package dot_db

import (
	"database/sql"
	"fmt"
	"dotBankNewsAggregator/dot_temp"
)

type NewItem struct {
	DB *sql.DB
	Note dot_temp.TempNews
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
