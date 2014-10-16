package buffer

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
			if len(pending) == 0 {
				break
			}
		}

	}

	for _, data := range pending {
		next <- data
	}
}
