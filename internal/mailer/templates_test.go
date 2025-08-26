package mailer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHtmlTemplateWatcher(t *testing.T) {
	dir, err := os.MkdirTemp("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	require.NoError(t, os.Mkdir(filepath.Join(dir, "templates"), 0o700))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "templates", "foo.html"), []byte(`{{ define "foo" }}foo{{ end }}`), 0o600))

	preDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(preDir)
	os.Chdir(dir)

	watcher, errCh, err := NewTemplateContainerWatcher("templates")
	require.NoError(t, err)
	require.NotNil(t, watcher)
	select {
	case msg := <-errCh:
		err = msg
	default:
		err = nil
	}

	require.NoError(t, err)
	defer watcher.Close()

	text, err := watcher.RenderToString("foo", &TemplateData{})
	require.NoError(t, err)
	assert.Equal(t, "foo", text)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "templates", "foo.html"), []byte(`{{ define "foo" }}bar{{ end }}`), 0o600))

	require.Eventually(t, func() bool {
		text, err := watcher.RenderToString("foo", &TemplateData{})
		return text == "bar" && err == nil
	}, time.Millisecond*100, time.Millisecond*50)
}

func TestNewTemplateContainerWatcher_BadDirectory(t *testing.T) {
	watcher, errChan, err := NewTemplateContainerWatcher("not_exists")
	require.Error(t, err)
	assert.Nil(t, watcher)
	assert.Nil(t, errChan)
}

func TestRender(t *testing.T) {
	tpl := template.New("test")
	_, err := tpl.Parse(`{{ define "foo" }}foo{{ .Props.Bar }}{{ end }}`)
	require.NoError(t, err)

	mt := NewTemplateContainerFromTemplate(tpl)
	data := &TemplateData{
		Props: map[string]any{"Bar": "bar"},
	}

	text, err := mt.RenderToString("foo", data)
	require.NoError(t, err)
	assert.Equal(t, "foobar", text)

	buf := &bytes.Buffer{}
	require.NoError(t, mt.Render(buf, "foo", data))
	assert.Equal(t, "foobar", buf.String())
}

func TestRenderError(t *testing.T) {
	tpl := template.New("test")
	_, err := tpl.Parse(`{{ define "foo" }}foo{{ .Foo.Bar }}bar{{ end }}`)
	require.NoError(t, err)

	mt := NewTemplateContainerFromTemplate(tpl)
	text, err := mt.RenderToString("foo", &TemplateData{})
	require.Error(t, err)
	assert.Equal(t, "", text)

	buf := bytes.Buffer{}
	assert.Error(t, mt.Render(&buf, "foo", &TemplateData{}))
	assert.Equal(t, "foo", buf.String())
}

func TestRenderUnknownTemplate(t *testing.T) {
	tpl := template.New("")
	mt := NewTemplateContainerFromTemplate(tpl)

	text, err := mt.RenderToString("foo", &TemplateData{})
	require.Error(t, err)
	assert.Equal(t, "", text)

	buf := &bytes.Buffer{}
	assert.Error(t, mt.Render(buf, "foo", &TemplateData{}))
	assert.Equal(t, "", buf.String())
}
