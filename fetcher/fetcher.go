package fetcher

import (
	"sync"
)

type Fetcher interface {
	Get() (string, error)
	List() ([]string, error)
}

var (
	once            sync.Once // используется при создании синглтона Fetcher
	sigletonFetcher *_Fetcher
)

func New() Fetcher {
	// Было бы полезно в конструкторе передать общий контекст (context.Context),
	// для определения завершения программы
	// По событию завершения можно было бы аккуратно завершить воркеры, закрыть
	// каналы т.е. подчистить ресурсы
	// Можно было бы реализовать это через метод Close(), но по ТЗ (как я понял)
	// поддержка одной копии (синглтона) осуществляется изнутри реализации =>
	// либо нужен счетчик ссылок (на вызовы конструтора), который уменьшался бы
	// в методе Close(). В момент, когда был равен нулю - происходила бы
	// очистка ресурсов.
	// Наверное, для этой реализации такого не требуется

	once.Do(func() {
		sigletonFetcher = newFetcher()
	})

	return sigletonFetcher
}

func newFetcher() *_Fetcher {
	f := &_Fetcher{
		// 2 - ограничение на число параллельных запросов
		// берем из ТЗ, могли бы из конфига или по количеству ядер
		sem: make(chan struct{}, 2),
	}

	return f
}

type _Fetcher struct {
	sem chan struct{}

	// можно использовать *http.Client для настройки таймаутов, пула коннектов
	// + для мока в тестах через https://golang.org/pkg/net/http/httptest/#example_Server
	// httpClient *http.Client
}

func (slf *_Fetcher) Get() (string, error) {
	slf.sem <- struct{}{} // ограничение на максимальное число параллельных запросов
	defer func() { <-slf.sem }()

	// тут запрос по http

	return "xxx", nil
}

func (slf *_Fetcher) List() ([]string, error) {
	slf.sem <- struct{}{} // ограничение на максимальное число параллельных запросов
	defer func() { <-slf.sem }()

	// тут запрос по http

	return []string{"xxx", "yyy", "zzz"}, nil
}
