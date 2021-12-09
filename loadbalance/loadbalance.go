package loadbalance

import (
	"errors"
	"hash/crc32"
	"math/rand"
)

var counter = make(map[string]int)

func Random(nodes []string) (string, error) {
	count := len(nodes)
	if count > 0 {
		r := generateNumber(0, count-1)
		return nodes[r], nil
	}
	return "", errors.New("404")
}

func RoundRobin(serviceName string, nodes []string) (string, error) {
	count := len(nodes)
	if count > 0 {
		r := counter[serviceName]
		if r >= count {
			r = 0
		}
		u := nodes[r]
		if r == count-1 {
			r = 0
		} else {
			r++
		}
		counter[serviceName] = r
		return u, nil
	}
	return "", errors.New("404")
}

func Hash(key string, nodes []string) (string, error) {
	count := len(nodes)
	if count > 0 {
		r := hash(key) % count
		return nodes[r], nil
	}
	return "", errors.New("404")
}

func ConsistentHash(key string, nodes []string) (string, error) {
	count := len(nodes)
	//make hashring
	if count > 0 {
		hash := newConsistentHash(nodes)
		if server, ok := hash.GetNode(key); ok {
			return server, nil
		}
	}
	return "", errors.New("404")
}

func generateNumber(min, max int) int {
	i := rand.Intn(max-min) + min
	return i
}

func hash(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
