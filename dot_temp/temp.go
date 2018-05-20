package dot_temp

/*
 * Шаблон новости
 */
type TempNews struct {
	Site string
	Link string
	Header string
	Description string
}

/*
 * Бокс для новостей
 */

type TempSetNews struct {
	SetNews []TempNews
}

/*
 * Шаблон карты сайта
 */
type TempSitemap struct {
	Sitemap []string `xml:"sitemap>loc"`
}

/*
 * Дочерние элементы в sitemap
 */

type TempSitemapUrls struct {
	SitemapUrls []string `xml:"url>loc"`
}