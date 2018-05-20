package main

import (
	"net/http"
	"io/ioutil"
	"regexp"
	"fmt"
	"encoding/xml"
	"html/template"
	"dotBankNewsAggregator/goquery"
	"dotBankNewsAggregator/dot_temp"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"dotBankNewsAggregator/dot_db"
)

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type msg string

func (m msg) ServeHTTP(resp http.ResponseWriter, req *http.Request){
	fmt.Fprint(resp, m)
}

/*
 * Получить код веб-страницы
 */
func getHtmlCode(page string) (string) {

	url := page

	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return BytesToString(html)
}

/*
 * Возвращает строковое представление байтового массива
 */
func BytesToString(data []byte) string {
	return string(data[:])
}

/*
 * Получает структуру, в которой массив url-ов
 * по конкретному xml.
 */
func getArrayUrls(childXml string) dot_temp.TempSitemapUrls  {
	var urlFromXml dot_temp.TempSitemapUrls
	resp, _ := http.Get(childXml)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &urlFromXml)
	return urlFromXml
}

/*
 * Получает url карты сайта из файла robots.txt
 */
func getSitemap(url string) (sitemap string){
	robots, err := http.Get(url + "/robots.txt")

	//Ошибка в синтаксисе url..
	if err != nil{
		panic(err)
	}

	defer robots.Body.Close()
	html, err := ioutil.ReadAll(robots.Body)

	//Не удалось прочитать robots.txt..
	if err != nil{
		panic(err)
	}

	str := BytesToString(html)
	re := regexp.MustCompile(`(https?://.+/.\S+)`)
	res := re.FindString(str)

	return res
}

/*
 * Возвращает конечную структуру новости (RamblerNews)
 */
func getContentSiteNews(page string, title string, description string) dot_temp.TempNews {
	res, err := http.Get(page)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		panic(res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	rn := dot_temp.TempNews{
		Link:page,
		Header:doc.Find(title).Text(),
		Description:doc.Find(description).Text(),
	}

	return rn
}

/*
 * Получает новости с сайта
 */
func getSetNews(anySite dot_temp.TempSitemap, infoSite map[string]string) dot_temp.TempSetNews  {
	var setNews dot_temp.TempSetNews

	infoHeader := infoSite["header"]
	infoDescription := infoSite["description"]

	//Подключение к базе MySQL
	con := "root@/dot_news?"
	con += "&charset=utf8"
	con += "&interpolateParams=true"

	db, err := sql.Open("mysql", con)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	for index, value := range anySite.Sitemap {
		fmt.Println(index, value)
		for index, value := range getArrayUrls(value).SitemapUrls {
			post := getContentSiteNews(value, infoHeader, infoDescription)
			fmt.Println(index, value)
			fmt.Println("Header: " + post.Header)
			fmt.Println("Description: " + post.Description)
			one_news := dot_temp.TempNews{
				Site: infoSite["site"],
				Link: value,
				Header: post.Header,
				Description: post.Description,
			}
			dot_db.Insert(dot_db.NewItem{ db, one_news })
			setNews.SetNews = append(setNews.SetNews, one_news)
		}
		break
	}

	defer db.Close()
	return setNews
}

/*
 * Главный обработчик карты сайта
 */
func HandleSitemap(w http.ResponseWriter, r *http.Request){
	site := r.FormValue("site1")
	sitemap := getSitemap(site)

	var smp dot_temp.TempSitemap
	resp, _ := http.Get(sitemap)
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &smp)

	infoSite := make(map[string]string)

	infoSite["site"] = site

	if site == "https://news.rambler.ru"{
		infoSite["header"] = ".big-title__title"
		infoSite["description"] = ".gallery__annotation"
	}else if site == "https://rb.ru"{
		infoSite["header"] = ".article__title"
		infoSite["description"] = ".post-announce p"
	}

	setNews := getSetNews(smp, infoSite)

	tmpl, _:= template.ParseFiles("views/result.html")
	tmpl.Execute(w, setNews)

	resp.Body.Close()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "views/index.html")
	})

	http.HandleFunc("/result", HandleSitemap)

	fmt.Println("Сервер запущен..")
	http.ListenAndServe(":8585", nil)
}
