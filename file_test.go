package e365_gateway

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"testing"
	"time"
)

func TestElmaFile(t *testing.T) {

	s := NewStand(testDefaultStandSettings)
	files := NewFileAdapter(s)

	testFileId := "68e8ecab-39e5-4566-ae15-b961a4f2cbee"
	testFileIdNotExisted := "68e8ecab-39e5-4566-ae15-b961a4f2cbef"
	testDirId := "ff715471-f756-4492-bb14-da941c55caf2"
	testDirIdNotExisted := "ff715471-f756-4492-bb14-da941c55caf3"

	ctxBg := context.Background()

	t.Run("get_link", func(t *testing.T) {
		link, err := files.GetDownloadLink(ctxBg, testFileId)
		require.NoError(t, err)
		require.NotEqual(t, "", link)
		require.Contains(t, link, "storage.yandexcloud.net/elma365-production")
	})

	t.Run("download_file", func(t *testing.T) {
		rc, err := files.DownloadFile(ctxBg, testFileId)
		require.NoError(t, err)
		bts, err := io.ReadAll(rc)
		require.NoError(t, err)
		require.Contains(t, string(bts), "text_content")
	})

	t.Run("get_dir_info", func(t *testing.T) {
		dirId := "ff715471-f756-4492-bb14-da941c55caf2"
		di, err := files.NewDirectory(dirId).Info(ctxBg)
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

		fileName := fmt.Sprintf("new_file_%s.txt", time.Now().String())
		file, err := files.NewDirectory(testDirId).Upload(ctxBg, buf, fileName)
		require.NoError(t, err)
		require.Equal(t, fileName, file.Name)
	})

	t.Run("full_success_scenario", func(t *testing.T) {

		fileName := fmt.Sprintf("test_%s.txt", time.Now().String())
		fileContent := []byte("test_content")
		buf := bytes.NewBuffer(fileContent)

		elmaDir := files.NewDirectory(testDirId)
		info, err := elmaDir.Info(ctxBg)
		require.NoError(t, err)
		require.Equal(t, testDirId, info.ID)

		f, err := elmaDir.Upload(ctxBg, buf, fileName)
		require.NoError(t, err)
		require.Equal(t, testDirId, f.Directory)
		require.Len(t, f.ID, uuid4Len)

		rc, err := files.DownloadFile(ctxBg, f.ID)
		require.NoError(t, err)
		defer func() {
			err = rc.Close()
			require.NoError(t, err)
		}()

		content, err := io.ReadAll(rc)
		require.NoError(t, err)
		require.Equal(t, fileContent, content)

	})

	t.Run("table_tests", func(t *testing.T) {

		t.Run("dir_info", func(t *testing.T) {
			testCases := []struct {
				name        string
				isSuccess   bool
				dirId       string
				ctx         context.Context
				expectedErr error
				timeout     time.Duration
			}{
				{name: "invalid_dir_id", dirId: "", ctx: nil, expectedErr: ErrInvalidID},
				{name: "failed_create_request", dirId: testDirId, ctx: nil, expectedErr: ErrCreateRequest},
				{name: "failed_send_request", dirId: testDirId, ctx: ctxBg, expectedErr: ErrSendRequest, timeout: time.Nanosecond},
				{name: "request_!ok", dirId: testDirIdNotExisted, ctx: ctxBg, expectedErr: ErrResponseStatusNotOK},
				{name: "success", isSuccess: true, dirId: testDirId, ctx: ctxBg, expectedErr: nil},
			}

			for _, tc := range testCases {

				t.Run(tc.name, func(t *testing.T) {
					dir := files.NewDirectory(tc.dirId)
					if tc.timeout > 0 {
						dir.SetClientTimeout(tc.timeout)
					}
					di, err := dir.Info(tc.ctx)
					require.ErrorIs(t, err, tc.expectedErr)
					if !tc.isSuccess {
						return
					}
					require.Equal(t, di.ID, tc.dirId)
				})
			}
		})

		t.Run("upload_to_dir", func(t *testing.T) {
			testCases := []struct {
				name        string
				isSuccess   bool
				dirId       string
				ctx         context.Context
				buf         *bytes.Buffer
				fileName    string
				expectedErr error
				timeout     time.Duration
			}{
				{name: "nil_buf", expectedErr: ErrNilItem},
				{name: "empty_buf1", buf: bytes.NewBuffer(nil), expectedErr: ErrEmptyBuffer},
				{name: "empty_buf2", buf: bytes.NewBuffer([]byte("")), expectedErr: ErrEmptyBuffer},
				{name: "invalid_dir_id", buf: bytes.NewBuffer([]byte("1")), expectedErr: ErrInvalidID},
				{name: "failed_create_request", dirId: testDirId, buf: bytes.NewBuffer([]byte("1")), expectedErr: ErrCreateRequest},
				{name: "empty_file_name", dirId: testDirId, ctx: ctxBg, buf: bytes.NewBuffer([]byte("1")), expectedErr: ErrResponseStatusNotOK},
				{name: "dir_not_existed", dirId: testDirIdNotExisted, ctx: ctxBg, buf: bytes.NewBuffer([]byte("1")), expectedErr: ErrResponseStatusNotOK},
				{name: "success", isSuccess: true, fileName: fmt.Sprintf("test_%s", time.Now().Format(time.DateTime)), dirId: testDirId, ctx: ctxBg, buf: bytes.NewBuffer([]byte("1")), expectedErr: nil},
			}

			for _, tc := range testCases {

				t.Run(tc.name, func(t *testing.T) {
					dir := files.NewDirectory(tc.dirId)
					if tc.timeout > 0 {
						dir.SetClientTimeout(tc.timeout)
					}
					f, err := dir.Upload(tc.ctx, tc.buf, tc.fileName)
					require.ErrorIs(t, err, tc.expectedErr)
					if !tc.isSuccess {
						return
					}
					require.Equal(t, f.Directory, tc.dirId)
					require.Contains(t, f.Name, tc.fileName)
					require.Len(t, f.ID, uuid4Len)
				})
			}
		})

		t.Run("get_download_link", func(t *testing.T) {
			testCases := []struct {
				name        string
				isSuccess   bool
				fileId      string
				ctx         context.Context
				expectedErr error
				timeout     time.Duration
			}{

				{name: "invalid_dir_id", expectedErr: ErrInvalidID},
				{name: "failed_create_request", fileId: testFileId, ctx: nil, expectedErr: ErrCreateRequest},
				{name: "file_not_existed", fileId: testFileIdNotExisted, ctx: ctxBg, expectedErr: ErrResponseStatusNotOK},
				{name: "success", isSuccess: true, fileId: testFileId, ctx: ctxBg, expectedErr: nil},
				{name: "failed_send_request", fileId: testFileId, ctx: ctxBg, timeout: time.Nanosecond, expectedErr: ErrSendRequest},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					if tc.timeout > 0 {
						files.SetClientTimeout(tc.timeout)
					}
					link, err := files.GetDownloadLink(tc.ctx, tc.fileId)
					require.ErrorIs(t, err, tc.expectedErr)
					if !tc.isSuccess {
						return
					}
					require.Contains(t, link, "storage.yandexcloud.net/elma365-production")
				})
			}
		})

	})

}
