package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/alex-shch/fetcher-test/fetcher"
)

// вообще, передал бы fetcher параметром, но т.к. необходимо поддержать
// параллельное создание + использование синглтона, конструируем внутри
func executor(ctx context.Context) error {
	f := fetcher.New()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if _, err := f.Get(); err != nil {
			return err
		}

		if _, err := f.List(); err != nil {
			return err
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	// магическая цифра 16 при использовании в одном месте и описанная в
	// комментарии не требует заведения константы
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := executor(ctx); err != nil {
				log.Print(err)
			}
			cancel() // если завершился хотя бы один воркер, завершим и остальные
		}()
	}

	time.Sleep(3 * time.Second) // или, ближе к реальности, завершение через signal.Notify

	cancel()  // событие завершения для воркеров
	wg.Wait() // ожидание, пока все воркеры завершаться
}
