package buffer

func Chan(in <-chan []byte, next chan<- []byte) {
	defer close(next)

	var top []byte
	pending := NewUnboundedBuffer(32*1024, 100*1024*1024)

recv:

	for {

		if Empty(pending) {
			data, ok := <-in
			if !ok {
				break
			}

			pending.Write(data)
			top = Top(pending)
		}

		select {
		case data, ok := <-in:
			if !ok {
				break recv
			}
			pending.Write(data)

		case next <- top:
			top = Top(pending)
		}

	}

	for !Empty(pending) || len(top) > 0 {
		next <- top
		top = Top(pending)
	}
}

func Top(buf Buffer) []byte {
	data := make([]byte, 32*1024)
	n, _ := buf.Read(data)
	data = data[:n]
	return data
}

func MemChan(in <-chan []byte, next chan<- []byte) {
	defer close(next)

	pending := make([][]byte, 0)

recv:

	for {

		if len(pending) == 0 {
			data, ok := <-in
			if !ok {
				break
			}

			pending = append(pending, data)
		}

		select {
		case data, ok := <-in:
			if !ok {
				break recv
			}
			pending = append(pending, data)

		case next <- pending[0]:
			pending = pending[1:]
		}

	}

	for _, data := range pending {
		next <- data
	}
}
