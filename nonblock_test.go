package condqueue

import (
	"sync"
	"testing"
	"time"
)

func TestQueue(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(10)
	q := NewCondQueue[int]()
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			if i&1 == 0 {
				q.Push(i)
			} else {
				n, _ := q.Pop()
				if n&1 != 0 {
					t.Errorf("misbehave: %d", n)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	wg.Add(2)
	go func() {
		time.AfterFunc(time.Second, func() {
			q.Push(1)
		})
		wg.Done()
	}()

	go func() {
		t.Log("wait")
		t1 := time.Now()
		q.Pop()
		t.Log("woken: ", time.Since(t1))
		wg.Done()
	}()
	wg.Wait()

	for i := 0; i < 10; i++ {
		q.Push(i)
	}
	q.Close()
	for i := 0; i < 10; i++ {
		n, _ := q.Pop()
		if n != i {
			t.Errorf("misbehave: %d", n)
		}
	}
	wg.Wait()

	n, ok := q.Pop()
	if ok {
		t.Errorf("misbehave: %d", n)
	}
	for i := 0; i < 10; i++ {
		q.Push(i)
	}
	q.ForEach(func(i int) bool {
		t.Logf("%d", i)
		return false
	})
}

func BenchmarkCond(b *testing.B) {
	cond := NewCondQueue[int]()

	var wg sync.WaitGroup
	wg.Add(2 * b.N)
	for i := 0; i < b.N; i++ {
		go func(id int) {
			cond.Push(id)
			wg.Done()
		}(i)
		go func() {
			cond.Pop()
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	ch := make(chan int)
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(2 * b.N)
	for i := 0; i < b.N; i++ {
		go func(id int) {
			ch <- id
			wg.Done()
		}(i)
		go func() {
			<-ch
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCondSelective(b *testing.B) {
	cond := NewCondQueue[int]()

	var wg sync.WaitGroup
	wg.Add(2 * b.N)
	for i := 0; i < b.N; i++ {
		go func(id int) {
			cond.Push(id)
			wg.Done()
		}(i)
		go func() {
			cond.Pop(true)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkChannelSelective(b *testing.B) {
	ch := make(chan int, b.N)
	b.ResetTimer()
	var wg sync.WaitGroup
	wg.Add(2 * b.N)
	for i := 0; i < b.N; i++ {
		go func(id int) {
			select {
			case ch <- id:
			default:
			}
			wg.Done()
		}(i)
		go func() {
			select {
			case <-ch:
			default:
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
