package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/LQR471814/scavenge/items"
)

type ToMarkdown struct{}

func (c ToMarkdown) HandleItem(ctx context.Context, item items.Item) (items.Item, error) {
	p, ok := items.CastItem[Page](item)
	if !ok {
		return nil, fmt.Errorf("invalid item")
	}

	md, err := htmltomarkdown.ConvertNode(
		p.contentHtml,
		converter.WithDomain(p.Metadata.Url),
	)
	if err != nil {
		return nil, err
	}
	p.Content = string(md)

	return items.Item{p}, nil
}

type ExportVolumes struct {
	ChunkSize uint64
	Dir       string

	written  uint64
	volumeNo int
	file     *os.File
	mutex    sync.Mutex
}

func NewExportVolumes(dir string, chunkSize uint64) *ExportVolumes {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		panic(err)
	}
	return &ExportVolumes{
		Dir:       dir,
		ChunkSize: chunkSize,
	}
}

func (e *ExportVolumes) write(content string) {
	defer e.mutex.Unlock()
	e.mutex.Lock()

	if e.written > e.ChunkSize || e.file == nil {
		if e.file != nil {
			e.file.Close()
		}

		e.volumeNo++
		e.written = 0

		f, err := os.Create(filepath.Join(e.Dir, fmt.Sprintf("volume-%d.md", e.volumeNo)))
		if err != nil {
			panic(err)
		}
		e.file = f
	}

	e.written += uint64(len(content))
	e.file.Write([]byte(content))
}

const header = `<hr>

**BEGIN SOURCE:** `

func (e *ExportVolumes) HandleItem(ctx context.Context, item items.Item) (items.Item, error) {
	page, _ := items.CastItem[Page](item)

	buff := strings.Builder{}
	buff.WriteString(header)
	buff.WriteString(page.Metadata.Url)
	buff.WriteString("\n\n")
	buff.WriteString(page.Content)
	buff.WriteString("\n\n")
	e.write(buff.String())

	return item, nil
}
