package lib

import (
	"math/rand"
	"regexp"
	"time"
)

type Dest struct {
	URL   string
	Words []string
}

type Counter struct {
	Dest  *Dest
	Count int
}

func (c *Counter) GetWord() string {
	l := len(c.Dest.Words)
	if c.Count > l {
		return ""
	}
	if c.Count < 0 {
		return ""
	}
	return c.Dest.Words[l-c.Count-1]
}

func (c *Counter) GetURL() string {
	return c.Dest.URL
}

func (c *Counter) GetCount() int {
	return c.Count
}

var links map[string]*Dest = make(map[string]*Dest)
var m map[string]*Counter = make(map[string]*Counter)

// var m map[string][]string = make(map[string][]string)

func GetDest(slug string) *Dest {
	dest := links[slug]
	return dest
}

func AddDest(link string, words []string) string {
	slug := GenSlug()
	links[slug] = &Dest{link, words}
	return slug
}

func NewCounter(dest *Dest) (*Counter, string) {
	l := len(dest.Words) - 1
	c := &Counter{dest, l}
	slug := GenSlug()
	m[slug] = c
	return c, slug
}

func PopNext(slug string) (*Counter, string) {
	c := m[slug]
	if c == nil {
		return nil, ""
	}
	defer delete(m, slug)

	c.Count = c.Count - 1
	if c.Count <= 0 {
		return c, ""
	}

	nextSlug := GenSlug()
	m[nextSlug] = c
	return c, nextSlug
}

func Get(slug string) (*Counter, string) {
	c, s := PopNext(slug)
	if c != nil {
		return c, s
	}

	dest := GetDest(slug)
	if dest == nil {
		return nil, ""
	}

	c, s = NewCounter(dest)
	return c, s
}

// const alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func GenSlug() string {
	n := 6

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

///////////////////////////////////////////////////////////

func ParseSlugFromUrl(url string) (string, bool) {
	re := regexp.MustCompile("/([a-zA-Z0-9-_]+)")
	matches := re.FindStringSubmatch(url)

	if !(len(matches) == 0) {
		return "", false
	}

	return matches[1], true
}
