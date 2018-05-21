Агрегатор новостей на Golang

Для ввода на главной странице:

1. https://news.rambler.ru
header: .big-title__title
description: .gallery__annotation

2. https://rb.ru
header: .article__title
description: .post-announce p

Чтобы избежать дублирования данных в базе, можно в insert добавить ignore и расставить необходимые индексы в самой таблице.
