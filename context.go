package main

import (
	"sync"

	"github.com/RadhiFadlillah/go-sastrawi"
)

type Context struct {
	Filename     string
	Translations []string

	mu     sync.Mutex
	output map[string]string

	tokenizer  sastrawi.Tokenizer
	stemmer    sastrawi.Stemmer
	dictionary sastrawi.Dictionary
}

func NewContext(filename string, translations []string) Context {
	return Context{
		Filename:     filename,
		Translations: translations,
		output:       make(map[string]string),

		tokenizer:  sastrawi.NewTokenizer(),
		stemmer:    sastrawi.NewStemmer(sastrawi.DefaultDictionary),
		dictionary: sastrawi.DefaultDictionary,
	}
}

func (c *Context) Run() map[string]string {
	var wg sync.WaitGroup
	for _, translation := range c.Translations {
		wg.Add(1)
		go c.checkTranslation(translation, &wg)
	}
	wg.Wait()
	return c.output
}

func (c *Context) checkTranslation(translation string, wg *sync.WaitGroup) {
	defer wg.Done()
	token := c.tokenizer.Tokenize(translation)
	for _, t := range token {
		if c.dictionary.Find(t) || len(t) <= 3 {
			continue
		}
		if c.stemmer.Stem(t) == t {
			c.mu.Lock()
			c.output[translation] = t + " bukan merupakan kata baku"
			c.mu.Unlock()
			break
		}
	}
}
