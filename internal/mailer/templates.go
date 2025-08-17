package mailer

import (
	"bytes"
	"html/template"
	"io"
	"path/filepath"
	"sync"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/fsnotify/fsnotify"
)

// TemplateContainer represents a set of templates that can be render
type TemplateContainer struct {
	templates *template.Template
	mutex     sync.RWMutex
	stop      chan struct{}
	stopped   chan struct{}
	watch     bool
}

// TemplateData contains the data used to populate the template variables, it has Props
// that can be of any type and HTML that only can be `template.HTML` types.
type TemplateData struct {
	Props map[string]any
	HTML  map[string]template.HTML
}

func (m *Mailer) NewTemplateData(lang string) (*TemplateData, error) {
	footer, err := models.Tr(lang, "templates.footer.part1", map[string]any{
		"SupportEmail": m.config().GetSupport().GetSupportEmail(),
		"SiteName":     m.config().Main.GetSiteName(),
	})
	if err != nil {
		return nil, err
	}

	return &TemplateData{Props: map[string]any{"Footer": footer}}, nil
}

func NewTemplateContainerFromTemplate(t *template.Template) *TemplateContainer {
	return &TemplateContainer{templates: t}
}

// NewTemplateContainer creates a new templates container
func NewTemplateContainer(dir string) (*TemplateContainer, error) {
	c := &TemplateContainer{}

	htmlTemplates, err := template.ParseGlob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}

	c.templates = htmlTemplates
	return c, nil
}

// NewTemplateContainerWatcher creates a new templates container scanning a directory and
// watch the directory filesystem changes to apply them to the loaded
// templates. This function returns the container and an errors channel to pass
// all errors that can happen during the watch process, or an regular error if
// we fail to create the templates or the watcher. The caller must consume the
// returned errors channel to ensure not blocking the watch process.
func NewTemplateContainerWatcher(dir string) (*TemplateContainer, <-chan error, error) {
	htmlTemplates, err := template.ParseGlob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, nil, err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, nil, err
	}

	err = watcher.Add(dir)
	if err != nil {
		watcher.Close()
		return nil, nil, err
	}

	c := &TemplateContainer{
		templates: htmlTemplates,
		watch:     true,
		stop:      make(chan struct{}),
		stopped:   make(chan struct{}),
	}
	errors := make(chan error)

	go func() {
		defer close(errors)
		defer close(c.stopped)
		defer watcher.Close()

		for {
			select {
			case <-c.stop:
				return
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if htmlTemplates, err := template.ParseGlob(filepath.Join(dir, "*.html")); err != nil {
						errors <- err
					} else {
						c.mutex.Lock()
						c.templates = htmlTemplates
						c.mutex.Unlock()
					}
				}
			case err := <-watcher.Errors:
				errors <- err
			}
		}
	}()

	return c, errors, nil
}

// Close stops the templates watcher of the container in case you have created
// it with watch parameter set to true
func (c *TemplateContainer) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.watch {
		close(c.stop)
		<-c.stopped
	}
}

// RenderToString renders the template referenced with the template name using
// the data provided and return a string with the result
func (c *TemplateContainer) RenderToString(templateName string, data TemplateData) (string, error) {
	var text bytes.Buffer
	if err := c.Render(&text, templateName, data); err != nil {
		return "", err
	}

	return text.String(), nil
}

// Render renders the template referenced with the template name using
// the data provided and write it to the writer provided
func (c *TemplateContainer) Render(w io.Writer, templateName string, data TemplateData) error {
	c.mutex.Lock()
	ht := c.templates
	c.mutex.Unlock()

	if err := ht.ExecuteTemplate(w, templateName, data); err != nil {
		return err
	}

	return nil
}
