package e365_gateway

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
	"time"
)

func TestElmaFile(t *testing.T) {

	s := NewStand("https://q3bamvpkvrulg.elma365.ru", "", "33ef3e66-c1cd-4d99-9a77-ddc4af2893cf")
	files := NewFileAdapter(s)

	testFileId := "68e8ecab-39e5-4566-ae15-b961a4f2cbee"

	t.Run("get_link", func(t *testing.T) {
		link, err := files.GetDownloadLink(testFileId)
		require.NoError(t, err)
		require.NotEqual(t, "", link)
		require.Contains(t, link, "storage.yandexcloud.net/elma365-production")
	})

	t.Run("download_file", func(t *testing.T) {
		rc, err := files.DownloadFile(testFileId)
		require.NoError(t, err)
		bts, err := io.ReadAll(rc)
		require.NoError(t, err)
		require.Contains(t, string(bts), "text_content")
	})

	t.Run("get_dir_info", func(t *testing.T) {
		dirId := "ff715471-f756-4492-bb14-da941c55caf2"
		di, err := files.NewDirectory(dirId).Info()
		require.NoError(t, err)
		require.Equal(t, "test", di.Name)
		require.Equal(t, dirId, di.ID)
	})

	t.Run("upload_file_2_dir", func(t *testing.T) {

		f, err := os.Open("test/file.txt")
		require.NoError(t, err)
		bts, err := io.ReadAll(f)
		require.NoError(t, err)
		buf := bytes.NewBuffer(bts)

		dirId := "ff715471-f756-4492-bb14-da941c55caf2"
		fileName := fmt.Sprintf("new_file_%s.txt", time.Now().String())
		file, err := files.NewDirectory(dirId).Upload(buf, fileName)
		require.NoError(t, err)
		require.Equal(t, fileName, file.Name)
	})

}
