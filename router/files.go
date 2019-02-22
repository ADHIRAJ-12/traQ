package router

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/traPtitech/traQ/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/traPtitech/traQ/model"
)

// PostFile POST /files
func (h *Handlers) PostFile(c echo.Context) error {
	userID := getRequestUserID(c)

	uploadedFile, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// アクセスコントロールリスト作成
	aclRead := repository.ACL{}
	if s := c.FormValue("acl_readable"); len(s) != 0 && s != "all" {
		for _, v := range strings.Split(s, ",") {
			uid, err := uuid.FromString(v)
			if err != nil || uid == uuid.Nil {
				return echo.NewHTTPError(http.StatusBadRequest, err)
			}

			if ok, err := h.Repo.UserExists(uid); err != nil {
				c.Logger().Error(err)
				return c.NoContent(http.StatusInternalServerError)
			} else if !ok {
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("unknown acl user id: %s", uid))
			}
			aclRead[uid] = true
		}
	} else {
		aclRead[uuid.Nil] = true
	}

	src, err := uploadedFile.Open()
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer src.Close()

	file, err := h.Repo.SaveFileWithACL(uploadedFile.Filename, src, uploadedFile.Size, uploadedFile.Header.Get(echo.HeaderContentType), model.FileTypeUserFile, userID, aclRead)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, file)
}

// GetFileByID GET /files/:fileID
func (h *Handlers) GetFileByID(c echo.Context) error {
	fileID := getRequestParamAsUUID(c, paramFileID)
	dl := c.QueryParam("dl")

	meta, err := h.Repo.GetFileMeta(fileID)
	if err != nil {
		c.Logger().Error()
		return c.NoContent(http.StatusInternalServerError)
	}

	// 直接アクセスURLが発行できる場合は、そっちにリダイレクト
	url, _ := h.Repo.GetFS().GenerateAccessURL(meta.GetKey())
	if len(url) > 0 {
		return c.Redirect(http.StatusFound, url)
	}

	file, err := h.Repo.GetFS().OpenFileByKey(meta.GetKey())
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer file.Close()

	c.Response().Header().Set(echo.HeaderContentLength, strconv.FormatInt(meta.Size, 10))
	c.Response().Header().Set(headerCacheControl, "private, max-age=31536000") //1年間キャッシュ
	c.Response().Header().Set(headerFileMetaType, meta.Type)

	switch meta.Type {
	case model.FileTypeStamp, model.FileTypeIcon:
		c.Response().Header().Set(headerCacheFile, "true")
	}

	if dl == "1" {
		c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", meta.Name))
	}

	return c.Stream(http.StatusOK, meta.Mime, file)
}

// DeleteFileByID DELETE /files/:fileID
func (h *Handlers) DeleteFileByID(c echo.Context) error {
	fileID := getRequestParamAsUUID(c, paramFileID)

	if err := h.Repo.DeleteFile(fileID); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMetaDataByFileID GET /files/:fileID/meta
func (h *Handlers) GetMetaDataByFileID(c echo.Context) error {
	fileID := getRequestParamAsUUID(c, paramFileID)

	meta, err := h.Repo.GetFileMeta(fileID)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, meta)
}

// GetThumbnailByID GET /files/:fileID/thumbnail
func (h *Handlers) GetThumbnailByID(c echo.Context) error {
	fileID := getRequestParamAsUUID(c, paramFileID)

	_, file, err := h.Repo.OpenThumbnailFile(fileID)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer file.Close()

	c.Response().Header().Set(headerCacheControl, "private, max-age=31536000") //1年間キャッシュ
	return c.Stream(http.StatusOK, mimeImagePNG, file)
}
