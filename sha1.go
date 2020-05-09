func sha1(message string, i1, i2, i3, i4, i5 int64) (int64, int64, int64, int64, int64) {
	h1 := int64(0x67452301) //1732584193
	h2 := int64(0xEFCDAB89) //4023233417
	h3 := int64(0x98BADCFE) //2562383102
	h4 := int64(0x10325476) //271733878
	h5 := int64(0xC3D2E1F0) //3285377520

	// process message in 512-bit chunks
	leftBound := 0
	rightBound := 512
	var w [80]string

	messageBin := makePadding(message)

	for rightBound <= len(messageBin) {
		chunk := messageBin[leftBound:rightBound]
		leftBound = leftBound + 512
		rightBound = rightBound + 512

		leftChunkBound := 0
		rightChunkBound := 32

		// split chunk into sixteen 32 bit words
		for i := 0; i < 16; i++ {
			w[i] = chunk[leftChunkBound:rightChunkBound]
			leftChunkBound = leftChunkBound + 32
			rightChunkBound = rightChunkBound + 32
		}

		//extend to 80 words
		for i := 16; i < 80; i++ {
			x1, err := strconv.ParseInt(w[i-3], 2, 64)
			if err != nil {
				log.Print(err)
			}
			x2, err := strconv.ParseInt(w[i-8], 2, 64)
			if err != nil {
				log.Print(err)
			}
			x3, err := strconv.ParseInt(w[i-14], 2, 64)
			if err != nil {
				log.Print(err)
			}
			x4, err := strconv.ParseInt(w[i-16], 2, 64)
			if err != nil {
				log.Print(err)
			}

			w[i] = fmt.Sprintf("%064b", (leftrotate(x1^x2^x3^x4, 1)))
		}

		a := h1
		b := h2
		c := h3
		d := h4
		e := h5
		var f int64
		var k int64

		//main loop
		for i := 0; i < 80; i++ {
			if 0 <= i && i <= 19 {
				f = d ^ (b & (c ^ d))
				k = 0x5A827999
			} else if 20 <= i && i <= 39 {
				f = b ^ c ^ d
				k = 0x6ED9EBA1
			} else if 40 <= i && i <= 59 {
				f = (b & c) | (b & d) | (c & d)
				k = 0x8F1BBCDC
			} else if 60 <= i && i <= 79 {
				f = b ^ c ^ d
				k = 0xCA62C1D6
			}

			bitWord, err := strconv.ParseInt(w[i], 2, 64)
			if err != nil {
				log.Fatal(err)
			}

			tmp := (leftrotate(int64(a), 5) + f + e + k + bitWord) & 0xffffffff

			e = d
			d = c
			c = leftrotate(int64(b), 30)
			b = a
			a = tmp

			// log.Printf("i=%d: %d %d %d %d %d", i, a, b, c, d, e)
		}

		// add chunk to hash result
		h1 = (h1 + a) & 0xffffffff
		h2 = (h2 + b) & 0xffffffff
		h3 = (h3 + c) & 0xffffffff
		h4 = (h4 + d) & 0xffffffff
		h5 = (h5 + e) & 0xffffffff

	}

	return h1, h2, h3, h4, h5
}

func makePadding(message string) string {
	// message to binary
	messageBin := stringToBin(message)
	messageBin = strings.Replace(messageBin, " ", "", -1)

	// we consider the binary message length
	messageLength := len(messageBin)
	log.Printf("Bitlength of message: %d", messageLength)

	// append bit 1
	messageBin = messageBin + "1"

	for len(messageBin)%512 != 448 {
		messageBin = messageBin + "0"
	}

	// append length of message in 64 bit
	originalLength := fmt.Sprintf("%064b", messageLength)
	messageBin = messageBin + originalLength

	return messageBin
}

func stringToBin(s string) (binString string) {
	for _, c := range s {
		binString = fmt.Sprintf("%s%.8b", binString, c)
	}
	return
}

func leftrotate(n, b int64) int64 {
	return ((n << b) | (n >> (32 - b))) & 0xffffffff

}
