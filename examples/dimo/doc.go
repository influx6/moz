// Package dimo contains a template within it's comments to generate code based off.
// 
//
/* All types of services and adapters are generated from the template below.
 @templater(id => Vars, gen => Partial.Go, {
		var (
			// defaultSendWithBeforeAbort defines the time to await receiving a message else aborting.
			defaultSendWithBeforeAbort = 3 * time.Second
		)
 })
 @templater(id => Service, gen => Partial.Go, {

 	import (
		"context"
		"time"
		"sync/atomic"
	)

 	//go:generate moz generate-file -fromFile ./{{sel "filename"}} -toDir ./impl/{{sel "Name" | lower}}

	{{ $noadapter := sel "NoAdapter" }}
	{{ if eq $noadapter "" }}
	// {{ sel "Name" }}FromByteAdapter defines a function that that will take a channel of bytes and return a channel of {{ sel "Type"}}.
 	type {{ sel "Name" }}FromByteAdapterWithContext func(context.Context, chan []byte) chan {{sel "Type"}}

	// {{ sel "Name" }}ToByteAdapter defines a function that that will take a channel of bytes and return a channel of {{ sel "Type"}}.
 	type {{ sel "Name" }}ToByteAdapter func(context.Context, chan {{sel "Type"}}) chan []byte

	{{ end}}

	// {{sel "Name"}}PartialCollect defines a function which returns a channel where the items of the incoming channel
	// are buffered until the channel is closed or the context expires returning whatever was collected, and closing the returning channel.
  // This function does not guarantee complete data.
	func {{sel "Name"}}PartialCollect(ctx context.Context, waitTime time.Duration, in chan {{ sel "Type" }}) chan []{{sel "Type"}} {
  	res := make(chan []{{sel "Type"}}, 0)

		go func() {
			var buffer []{{sel "Type"}}

			t := time.NewTimer(waitTime)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
					res <- buffer
          close(res)
					return
				case data, ok := <-in:
					if !ok {
						res <- buffer
            close(res)
						return
					}

					buffer = append(buffer, data)
					continue
				case <-t.C:
					t.Reset(waitTime)
					continue
				}
			}
		}()

		return res
	}

	// {{sel "Name"}}Collect defines a function which returns a channel where the items of the incoming channel
	// are buffered until the channel is closed, nothing will be returned if the channel given is not closed  or the context expires.
  // Once done, returning channel is closed.
  // This function guarantees complete data.
	func {{sel "Name"}}Collect(ctx context.Context, waitTime time.Duration, in chan {{ sel "Type" }}) chan []{{sel "Type"}} {
  	res := make(chan []{{sel "Type"}}, 0)

		go func() {
			var buffer []{{sel "Type"}}

			t := time.NewTimer(waitTime)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
          close(res)
					return
				case data, ok := <-in:
					if !ok {
						res <- buffer
            close(res)
						return
					}

					buffer = append(buffer, data)
					continue
				case <-t.C:
					t.Reset(waitTime)
					continue
				}
			}
		}()

		return res
	}

	// {{sel "Name"}}Mutate defines a function which returns a channel where the items of the incoming channel
	// are mutated based on a function, till the provided channel is closed.
  // If the given channel is closed or if the context expires, the returning channel is closed as well.
  // This function guarantees complete data.
	func {{sel "Name"}}Mutate(ctx context.Context, waitTime time.Duration, mutateFn func({{ sel "Type" }}) {{sel "Type"}}, in chan {{ sel "Type" }}) chan {{sel "Type"}} {
  	res := make(chan {{sel "Type"}}, 0)

		go func() {
			t := time.NewTimer(waitTime)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
          close(res)
					return
				case data, ok := <-in:
					if !ok {
            close(res)
						return
					}

					res <- mutateFn(data)
				case <-t.C:
					t.Reset(waitTime)
					continue
				}
			}
		}()

		return res
	}

	// {{sel "Name"}}Filter defines a function which returns a channel where the items of the incoming channel
	// are filtered based on a function, till the provided channel is closed.
  // If the given channel is closed or if the context expires, the returning channel is closed as well.
  // This function guarantees complete data.
	func {{sel "Name"}}Filter(ctx context.Context, waitTime time.Duration, filterFn func({{ sel "Type" }}) bool, in chan {{ sel "Type" }}) chan {{sel "Type"}} {
  	res := make(chan {{sel "Type"}}, 0)

		go func() {
			t := time.NewTimer(waitTime)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
          close(res)
					return
				case data, ok := <-in:
					if !ok {
            close(res)
						return
					}

          if !filterFn(data){
            continue
          }

					res <- data
				case <-t.C:
					t.Reset(waitTime)
					continue
				}
			}
		}()

		return res
	}

	// {{sel "Name"}}CollectUntil defines a function which returns a channel where the items of the incoming channel
	// are buffered until the data matches a given requirement provided through a function. If the function returns true
  // then currently buffered data is returned and a new buffer is created. This is useful for batch collection based on
  // specific criteria. If the channel is closed before the criteria is met, what data is left is sent down the returned channel,
  // closing that channel. If the context expires then data gathered is returned and returning channel is closed.
  // This function guarantees some data to be delivered.
	func {{sel "Name"}}CollectUntil(ctx context.Context, waitTime time.Duration, condition func([]{{sel "Type"}}) bool, in chan {{ sel "Type" }}) chan []{{sel "Type"}} {
  	res := make(chan []{{sel "Type"}}, 0)

		go func() {
			var buffer []{{sel "Type"}}

			t := time.NewTimer(waitTime)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
          close(res)
					return
				case data, ok := <-in:
					if !ok {
						res <- buffer
            close(res)
						return
					}

					buffer = append(buffer, data)

          // If we do not match the given criteria, then continue buffering.
          if condition(buffer) {
            continue
          }

          // We do match criteria, send buffered data, and reset buffer
					res <- buffer
          buffer = nil
				case <-t.C:
					t.Reset(waitTime)
					continue
				}
			}
		}()

		return res
	}

  // {{sel "Name"}}MergeWithoutOrder merges the incoming data from the {{sel "Type"}} into a single stream of {{sel "Type"}},
  // merge collects data from all provided channels in turns, each giving a specified time to deliver, else
  // skipped until another turn. Once all data is collected from each sender, then the data set is merged into
  // a single slice and delivered to the returned channel.
  // MergeWithoutOrder makes the following guarantees:
  // 1. Items will be received in order and return in order of channels provided.
  // 2. Returned channel will never return incomplete data where a channels return value is missing from its index.
  // 3. If any channel is closed then all operation is stopped and the returned channel is closed.
  // 4. If the context expires before items are complete then returned channel is closed.
  // 5. It will continously collect order data for all channels until any of the above conditions are broken.
  // 6. All channel data are collected once for the receiving scope, i.e a channels data will not be received twice into the return slice
  //    but all channels will have a single data slot for a partial data collection session.
  // 7. Will continue to gather data from provided channels until all are closed or the context has expired.
  // 8. If any of the senders is nil then the returned channel will be closed, has this leaves things in an unstable state.
  func {{sel "Name"}}MergeWithoutOrder(ctx context.Context, maxWaitTime time.Duration, senders ...chan {{ sel "Type" }}) chan []{{ if hasPrefix (sel "Type") "[]" }} {{ (trimPrefix (sel "Type") "[]") }} {{ else }} {{ sel "Type" }} {{end}} {
  	res := make(chan []{{ if hasPrefix (sel "Type") "[]" }} {{ (trimPrefix (sel "Type") "[]") }} {{ else }} {{ sel "Type" }} {{end}}, 0)

  	for _, elem := range senders {
  		if elem == nil {
  			close(res)
  			return res
  		}
  	}

  	go func() {
  		var index int

  		total := len(senders)
  		filled := make(map[int]{{ sel "Type" }}, 0)

  		for {
  			// if the current index has being filled, shift forward and reattempt loop.
  			if _, ok := filled[index]; ok {
  				index++
  				continue
  			}

  			if len(filled) == total {
  				var content []{{ if hasPrefix (sel "Type") "[]" }} {{ (trimPrefix (sel "Type") "[]") }} {{ else }} {{ sel "Type" }} {{end}}

  				for _, item := range filled {
  					content = append(content, {{ if hasPrefix (sel "Type") "[]" }} item... {{ else }} item {{end}})
  				}

  				res <- content

  				index = 0
  				filled = make(map[int]{{ sel "Type" }}, 0)
  			}

  			timer := time.NewTimer(maxWaitTime)

  			select {
  			case <-ctx.Done():
  				close(res)
  				timer.Stop()
  				return
  			case <-timer.C:
  				switch index >= total {
  				case true:
  					index = 0
  				case false:
  					index++
  				}
  			case data, ok := <-senders[index]:
  				if !ok {
  					close(res)
  					timer.Stop()
  					return
  				}

  				filled[index] = data
  				index++
  			}

  			timer.Stop()
  		}
  	}()

  	return res
  }

  // {{sel "Name"}}MergeInOrder merges the incoming data from the {{sel "Type"}} into a single stream of {{sel "Type"}},
  // merge collects data from all provided channels in turns, each giving a specified time to deliver, else
  // skipped until another turn. Once all data is collected from each sender, then the data set is merged into
  // a single slice and delivered to the returned channel.
  // MergeInOrder makes the following guarantees:
  // 1. Items will be received in order and return in order of channels provided.
  // 2. Returned channel will never return incomplete data where a channels return value is missing from its index.
  // 3. If any channel is closed then all operation is stopped and the returned channel is closed.
  // 4. If the context expires before items are complete then returned channel is closed.
  // 5. It will continously collect order data for all channels until any of the above conditions are broken.
  // 6. All channel data are collected once for the receiving scope, i.e a channels data will not be received twice into the return slice
  //    but all channels will have a single data slot for a partial data collection session.
  // 7. Will continue to gather data from provided channels until all are closed or the context has expired.
  // 8. If any of the senders is nil then the returned channel will be closed, has this leaves things in an unstable state.
  func {{sel "Name"}}MergeInOrder(ctx context.Context, maxWaitTime time.Duration, senders ...chan {{ sel "Type" }}) chan []{{ if hasPrefix (sel "Type") "[]" }} {{ (trimPrefix (sel "Type") "[]") }} {{ else }} {{ sel "Type" }} {{end}} {
  	res := make(chan []{{ if hasPrefix (sel "Type") "[]" }} {{ (trimPrefix (sel "Type") "[]") }} {{ else }} {{ sel "Type" }} {{end}}, 0)

  	for _, elem := range senders {
  		if elem == nil {
  			close(res)
  			return res
  		}
  	}

  	go func() {
  		var index int

  		total := len(senders)
  		filled := make(map[int]{{ sel "Type" }}, 0)

  		for {
  			// if the current index has being filled, shift forward and reattempt loop.
  			if _, ok := filled[index]; ok {
  				index++
  				continue
  			}

  			if len(filled) == total {
  				var content []{{ if hasPrefix (sel "Type") "[]" }} {{ (trimPrefix (sel "Type") "[]") }} {{ else }} {{ sel "Type" }} {{end}}

  				for index := range senders {
            item := filled[index]
  					content = append(content, {{ if hasPrefix (sel "Type") "[]" }} item... {{ else }} item {{end}})
  				}

  				res <- content

  				index = 0
  				filled = make(map[int]{{ sel "Type" }}, 0)
  			}

  			timer := time.NewTimer(maxWaitTime)

  			select {
  			case <-ctx.Done():
  				close(res)
  				timer.Stop()
  				return
  			case <-timer.C:
  				switch index >= total {
  				case true:
  					index = 0
  				case false:
  					index++
  				}
  			case data, ok := <-senders[index]:
  				if !ok {
  					close(res)
  					timer.Stop()
  					return
  				}

  				filled[index] = data
  				index++
  			}

  			timer.Stop()
  		}
  	}()

  	return res
  }

  // {{sel "Name"}}CombineParitallyWithoutOrder receives a giving stream of content from multiple channels, returning a single channel of a
  // 2d slice, it sequentially tries to recieve data from each sender within a provided time duration, else skipping until it's next turn.
  // This ensures every sender has adquate time to receive and reducing long blocked waits for a specific sender, more so, this reduces the
  // overhead of managing multiple go-routined receving channels which are prone to goroutine ophaning or memory leaks.
  // {{sel "Name"}}CombineParitallyWithoutOrder makes the following guarantees:
  // 1. Items will be received in any order received from the channels provided.
  // 2. Returned channel will never return incomplete data where a channels return value is missing from its index.
  // 3. If any channel is closed then all operation is stopped and the returned channel is closed.
  // 4. If the context expires before items are complete then returned channel is closed.
  // 5. It will continously collect data in any order for all channels until any of the above conditions are broken.
  // 6. All channel data are collected once for the receiving scope, i.e a channels data will not be received twice into the return slice
  //    but all channels will have a single data slot for a partial data collection session.
  // 7. Will continue to gather data from provided channels until all are closed or the context has expired.
  // 8. If any of the senders is nil then the returned channel will be closed, has this leaves things in an unstable state.
  func {{sel "Name"}}CombinePartiallyWithoutOrder(ctx context.Context, maxItemWait time.Duration, senders ...chan {{sel "Type"}}) chan []{{sel "Type"}} {
  	res := make(chan []{{sel "Type"}}, 0)

  	for _, elem := range senders {
  		if elem == nil {
  			close(res)
  			return res
  		}
  	}

  	go func() {
  		content := make([]{{sel "Type"}}, 0)

  		var index int

  		total := len(senders)
  		filled := make(map[int]bool, 0)
  		closed := make(map[int]bool, 0)

      var sendersClosed int

  		for {
  			// if the current index has being filled, shift forward and re-attempt loop.
  			if filled[index] || closed[index]{
  				index++
  				continue
  			}

  			if len(content) == total {
  				res <- content

  				index = 0
  				filled = make(map[int]bool, 0)
  				content = make([]{{sel "Type"}}, len(senders))
  			}

  			timer := time.NewTimer(maxItemWait)

  			select {
  			case <-ctx.Done():
  				close(res)
  				timer.Stop()
  				return
  			case <-timer.C:
  				switch index >= total {
  				case true:
  					index = 0
  				case false:
  					index++
  				}
  			case data, ok := <-senders[index]:
  				if !ok {
  					timer.Stop()

            sendersClosed++
            closed[index] = true
  					continue
  				}

  				content = append(content, data)
  				filled[index] = true
  				index++
  			}

  			timer.Stop()
  		}

  	}()

  	return res
  }

  // {{sel "Name"}}CombineWithoutOrder receives a giving stream of content from multiple channels, returning a single channel of a
  // 2d slice, it sequentially tries to recieve data from each sender within a provided time duration, else skipping until it's next turn.
  // This ensures every sender has adquate time to receive and reducing long blocked waits for a specific sender, more so, this reduces the
  // overhead of managing multiple go-routined receving channels which are prone to goroutine ophaning or memory leaks.
  // {{sel "Name"}}CombineWithoutOrder makes the following guarantees:
  // 1. Items will be received in any order received from the channels provided.
  // 2. Returned channel will never return incomplete data where a channels return value is missing from its index.
  // 3. If any channel is closed then all operation is stopped and the returned channel is closed.
  // 4. If the context expires before items are complete then returned channel is closed.
  // 5. It will continously collect data in any order for all channels until any of the above conditions are broken.
  // 6. All channel data are collected once for the receiving scope, i.e a channels data will not be received twice into the return slice
  //    but all channels will have a single data slot for a partial data collection session.
  // 7. Will continue to gather data from provided channels until all are closed or the context has expired.
  // 8. If any of the senders is nil then the returned channel will be closed, has this leaves things in an unstable state.
  func {{sel "Name"}}CombineWithoutOrder(ctx context.Context, maxItemWait time.Duration, senders ...chan {{sel "Type"}}) chan []{{sel "Type"}} {
  	res := make(chan []{{sel "Type"}}, 0)

  	for _, elem := range senders {
  		if elem == nil {
  			close(res)
  			return res
  		}
  	}

  	go func() {
  		content := make([]{{sel "Type"}}, 0)

  		var index int

  		total := len(senders)
  		filled := make(map[int]bool, 0)

  		for {
  			// if the current index has being filled, shift forward and reattempt loop.
  			if filled[index] {
  				index++
  				continue
  			}

  			if len(content) == total {
  				res <- content
  				index = 0
  				filled = make(map[int]bool, 0)
  				content = make([]{{sel "Type"}}, len(senders))
  			}

  			timer := time.NewTimer(maxItemWait)

  			select {
  			case <-ctx.Done():
  				close(res)
  				timer.Stop()
  				return
  			case <-timer.C:
  				switch index >= total {
  				case true:
  					index = 0
  				case false:
  					index++
  				}
  			case data, ok := <-senders[index]:
  				if !ok {
  					close(res)
  					timer.Stop()
  					return
  				}

  				content = append(content, data)
  				filled[index] = true
  				index++
  			}

  			timer.Stop()
  		}

  	}()

  	return res
  }

  // {{sel "Name"}}CombineInPartialOrder receives a giving stream of content from multiple channels, returning a single channel of a
  // 2d slice, it sequentially tries to recieve data from each sender within a provided time duration, else skipping until it's next turn.
  // This ensures every sender has adquate time to receive and reducing long blocked waits for a specific sender, more so, this reduces the
  // overhead of managing multiple go-routined receving channels which are prone to goroutine ophaning or memory leaks.
  // {{sel "Name"}}CombineInPartialOrder makes the following guarantees:
  // 1. Items will be received in order and return in order of channels provided.
  // 2. Returned channel will never return incomplete data where a channels return value is missing from its index.
  // 3. If any channel is closed then all operation is not stopped and will continue but there will be an empty slot in the slice returned.
  // 4. If the context expires before items are complete then returned channel is closed.
  // 5. It will continously collect order data for all channels until any of the above conditions are broken.
  // 6. All channel data are collected once for the receiving scope, i.e a channels data will not be received twice into the return slice
  //    but all channels will have a single data slot for a partial data collection session.
  // 7. Will continue to gather data from provided channels until all are closed or the context has expired.
  // 8. If any of the senders is nil then the returned channel will be closed, has this leaves things in an unstable state.
  func {{sel "Name"}}CombineInPartialOrder(ctx context.Context, maxItemWait time.Duration, senders ...chan {{sel "Type"}}) chan []{{sel "Type"}} {
  	res := make(chan []{{sel "Type"}}, 0)

  	for _, elem := range senders {
  		if elem == nil {
  			close(res)
  			return res
  		}
  	}

  	go func() {
  		content := make([]{{sel "Type"}}, len(senders))

  		var index int

  		total := len(senders)
  		filled := make(map[int]bool, 0)
  		closed := make(map[int]bool, 0)

      var sendersClosed int

  		for {
  			// if the current index has being filled, shift forward and reattempt loop.
  			if filled[index] || closed[index] {
  				index++
  				continue
  			}

        if sendersClosed >= total {
  				res <- content
          return
        }

  			if len(filled) == total {
  				res <- content
  				index = 0
  				filled = make(map[int]bool, 0)
  				content = make([]{{sel "Type"}}, len(senders))
  			}

  			timer := time.NewTimer(maxItemWait)

  			select {
  			case <-ctx.Done():
  				close(res)
  				timer.Stop()
  				return
  			case <-timer.C:
  				switch index >= total {
  				case true:
  					index = 0
  				case false:
  					index++
  				}
  			case data, ok := <-senders[index]:
  				if !ok {
  					timer.Stop()

            sendersClosed++
            closed[index] = true
  					continue
  				}

  				content[index] = data
  				filled[index] = true
  				index++
  			}

  			timer.Stop()
  		}

  	}()

  	return res
  }

  // {{sel "Name"}}CombineInOrder receives a giving stream of content from multiple channels, returning a single channel of a
  // two slice type, more so, it will use the maxItemWait duration to give every channel a opportunity for delivery
  // data. If the time is passed it will cycle to the next item, till the complete set is retrieved from, unless
  // the context timout expires and causes a complete stop of operation. CombineInOrder guarantees that the data
  // retrieved will be stored in order of passed in channels. If any of the channels is nil, then the returned
  // channel will be closed.
  // CombineInOrder makes the following guarantees:
  // 1. Items will be received in order and return in order of channels provided.
  // 2. Returned channel will never return incomplete data where a channels return value is missing from its index.
  // 3. If any channel is closed then all operation is stopped and the returned channel is closed.
  // 4. If the context expires before items are complete then returned channel is closed.
  // 5. It will continously collect order data for all channels until any of the above conditions are broken.
  // 6. All channel data are collected once for the receiving scope, i.e a channels data will not be received twice into the return slice
  //    but all channels will have a single data slot for a partial data collection session.
  // 7. Will continue to gather data from provided channels until all are closed or the context has expired.
  // 8. If any of the senders is nil then the returned channel will be closed, has this leaves things in an unstable state.
  func {{sel "Name"}}CombineInOrder(ctx context.Context, maxItemWait time.Duration, senders ...chan {{sel "Type"}}) chan []{{sel "Type"}} {
  	res := make(chan []{{sel "Type"}}, 0)

  	for _, elem := range senders {
  		if elem == nil {
  			close(res)
  			return res
  		}
  	}

  	go func() {
  		content := make([]{{sel "Type"}}, len(senders))

  		var index int

  		total := len(senders)
  		filled := make(map[int]bool, 0)

  		for {
  			// if the current index has being filled, shift forward and reattempt loop.
  			if filled[index] {
  				index++
  				continue
  			}

  			if len(filled) == total {
  				res <- content
  				index = 0
  				filled = make(map[int]bool, 0)
  				content = make([]{{sel "Type"}}, len(senders))
  			}

  			timer := time.NewTimer(maxItemWait)

  			select {
  			case <-ctx.Done():
  				close(res)
  				timer.Stop()
  				return
  			case <-timer.C:
  				switch index >= total {
  				case true:
  					index = 0
  				case false:
  					index++
  				}
  			case data, ok := <-senders[index]:
  				if !ok {
  					close(res)
  					timer.Stop()
  					return
  				}

  				content[index] = data
  				filled[index] = true
  				index++
  			}

  			timer.Stop()
  		}

  	}()

  	return res
  }

	// {{ sel "Name" }}Distributor delivers messages to subscription channels which it manages internal,
	// ensuring every subscriber gets the delivered message, it guarantees that every subscribe will get the chance to receive a
	// messsage unless it takes more than a giving duration of time, and if passing that duration that the operation
	// to deliver will be cancelled, this ensures that we do not leak goroutines nor
	// have eternal channel blocks.
	type {{ sel "Name" }}Distributor struct {
		running             int64
		messages            chan {{ sel "Type" }}
		close               chan struct{}
		clear               chan struct{}
		subscribers         []chan {{ sel "Type" }}
		newSub              chan chan {{ sel "Type" }}
		sendWaitBeforeAbort time.Duration
	}

	// New{{sel "Name"}}Disributor returns a new instance of a {{ sel "Name" }}Distributor.
	func New{{sel "Name"}}Disributor(buffer int, sendWaitBeforeAbort time.Duration) *{{ sel "Name" }}Distributor {
		if sendWaitBeforeAbort <= 0 {
			sendWaitBeforeAbort = defaultSendWithBeforeAbort
		}

		return &{{ sel "Name" }}Distributor{
			clear:               make(chan struct{}, 0),
			close:               make(chan struct{}, 0),
			subscribers:         make([]chan {{ sel "Type" }}, 0),
			newSub:              make(chan chan {{ sel "Type" }}, 0),
			messages:            make(chan {{ sel "Type" }}, buffer),
			sendWaitBeforeAbort: sendWaitBeforeAbort,
		}
	}

	// PublishDeadline sends the message into the distributor to be delivered to all subscribers if it has not
	// passed the provided deadline.
	func (d *{{ sel "Name" }}Distributor) PublishDeadline(message {{ sel "Type" }}, dur time.Duration) {
		if atomic.LoadInt64(&d.running) == 0 {
			return
		}

		timer := time.NewTimer(dur)
		defer timer.Stop()

		select {
		case <-timer.C:
			return
		case d.messages <- message:
			return
		}
	}

	// Publish sends the message into the distributor to be delivered to all subscribers.
	func (d *{{ sel "Name" }}Distributor) Publish(message {{ sel "Type" }}) {
		if atomic.LoadInt64(&d.running) == 0 {
			return
		}

		d.messages <- message
	}

	// Subscribe adds the channel into the distributor subscription lists.
	func (d *{{ sel "Name" }}Distributor) Subscribe(sub chan {{ sel "Type" }}) {
		if atomic.LoadInt64(&d.running) == 0 {
			return
		}

		d.newSub <- sub
	}

	// Clear removes all subscribers from the distributor.
	func (d *{{ sel "Name" }}Distributor) Clear() {
		if atomic.LoadInt64(&d.running) == 0 {
			return
		}

		d.clear <- struct{}{}
	}

	// Stop halts internal delivery behaviour of the distributor.
	func (d *{{ sel "Name" }}Distributor) Stop() {
		if atomic.LoadInt64(&d.running) == 0 {
			return
		}

		d.close <- struct{}{}
	}

	// Start initializes the distributor to deliver messages to subscribers.
	func (d *{{ sel "Name" }}Distributor) Start() {
		if atomic.LoadInt64(&d.running) != 0 {
			return
		}

		atomic.AddInt64(&d.running, 1)
		go d.manage()
	}

	// manage implements necessary logic to manage message delivery and
	// subscriber adding
	func (d *{{ sel "Name" }}Distributor) manage() {
		defer atomic.AddInt64(&d.running, -1)

		for {
			select {
			case <-d.clear:
				d.subscribers = nil

			case newSub, ok := <-d.newSub:
				if !ok {
					return
				}

				d.subscribers = append(d.subscribers, newSub)
			case message, ok := <-d.messages:
				if !ok {
					return
				}

				for _, sub := range d.subscribers {
					go func(c chan {{ sel "Type" }}) {
						tick := time.NewTimer(d.sendWaitBeforeAbort)
						defer tick.Stop()

						select {
						case sub <- message:
							return
						case <-tick.C:
							return
						}
					}(sub)
				}
			case <-d.close:
				return
			}
		}
	}

	// Mono{{ sel "Name" }}Service defines a interface for underline systems which want to communicate like
	// in a stream using channels to send/recieve only {{sel "Type"}} values from a single source.
	// It allows different services to create adapters to
	// transform data coming in and out from the Service.
	// Auto-Generated using the moz code-generator https://github.com/influx6/moz.
	// @iface
 	type Mono{{ sel "Name" }}Service interface {
		// Send will return a channel which will allow reading from the Service it till it is closed.
		Send() (<-chan {{ sel "Type"}}, error)

		// Receive will take the channel, which will be written into the Service for it's internal processing
		// and the Service will continue to read form the channel till the channel is closed.
		// Useful for collating/collecting services.
		Receive(<-chan {{sel "Type"}}) error

		// Done defines a signal to other pending services to know whether the Service is still servicing
		// request.
		Done() chan struct{}

		// Errors returns a channel which signals services to know whether the Service is still servicing
		// request.
		Errors() (chan error)

		// Service defines a function to be called to stop the Service internal operation and to close
		// all read/write operations.
		Stop() error
	}

	// {{ sel "Name" }}Service defines a interface for underline systems which want to communicate like
	// in a stream using channels to send/recieve {{sel "Type"}} values. It allows different services to create adapters to
	// transform data coming in and out from the Service.
	// Auto-Generated using the moz code-generator https://github.com/influx6/moz.
	// @iface
 	type {{ sel "Name" }}Service interface {
		// Send will return a channel which will allow reading from the Service it till it is closed.
		Send(string) (<-chan {{ sel "Type"}}, error)

		// Receive will take the channel, which will be written into the Service for it's internal processing
		// and the Service will continue to read form the channel till the channel is closed.
		// Useful for collating/collecting services.
		Receive(string, <-chan {{sel "Type"}}) error

		// Done defines a signal to other pending services to know whether the Service is still servicing
		// request.
		Done() chan struct{}

		// Errors returns a channel which signals services to know whether the Service is still servicing
		// request.
		Errors() (chan error)

		// Service defines a function to be called to stop the Service internal operation and to close
		// all read/write operations.
		Stop() error
	}

})
*/
//
// @templaterTypesFor(id => Vars, filename => vars.go)
//
package dimo


//go:generate moz generate --toDir=./