# Cond Queue(管程队列)

Go's Channel is faster than cond, however, channel isn't resizeable.

Like the following case,
```
ch := make(chan int, 128)
for i := 0; i < 256; i++ {
   // when i is above the 128,
   // the channel got blocked.
   ch <- i
}
```

In some cases, we don't want the producer got blocked.

So this is the usage of cond queue, it's resizeable and will never block the producer.

# Example

```
// producer
queue := NewCondQueue[int]()
for i := 0; i < 256; i++ {
     queue.Push(i)
}
...
// consumer
for {
    // Pop will block until producer pushing value into the queue
    value, ok := queue.Pop()
    if !ok {
       // queue is closed
    }
}

// consumer non-block
// in some cases, we need to do fast check of queue.
// if passing true to Pop(), it will not block anymore.
// if there's no element in the queue, ok is false.
value, ok := queue.Pop(true)

// Close
// Like buffered channel, producer cannot push after Close()
// But consumers can pop until no element.
queue.Close()

// Foreach
// return true to stop while false to continue
queue.ForEach(func (T) bool {
   
})
```

