package cloudinary

import (
	"context"
	"crop_connect/helper"
	"crop_connect/util"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Function interface {
	UploadOneWithFilename(folder string, file *multipart.FileHeader, filename string) (string, error)
	UploadOneWithGeneratedFilename(folder string, file *multipart.FileHeader) (string, error)
	UploadManyWithGeneratedFilename(folder string, files []*multipart.FileHeader) ([]string, error)
	RenameOneByFilename(folder string, oldFilename string, newFilename string) (string, error)
	DeleteOneByFilename(folder string, filename string) error
	DeleteOneByURL(folder string, URL string) error
	DeleteManyByURL(folder string, URLs []string) error
	UpdateArrayImage(folder string, imageURLs []string, updateImage []*helper.UpdateImage) ([]string, error)
}

type Cloudinary struct {
	cloudinary *cloudinary.Cloudinary
}

var (
	folderBase string
)

func Init(folderName string) Function {
	cld, err := cloudinary.NewFromParams(util.GetConfig("CLOUDINARY_CLOUD_NAME"), util.GetConfig("CLOUDINARY_API_KEY"), util.GetConfig("CLOUDINARY_API_SECRET"))
	if err != nil {
		panic(err)
	}

	folderBase = folderName

	cld.Config.URL.Secure = true
	return &Cloudinary{
		cloudinary: cld,
	}
}

func (c *Cloudinary) UploadOneWithFilename(folder string, file *multipart.FileHeader, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pictureBuffer, err := file.Open()
	if err != nil {
		return "", err
	}

	resp, err := c.cloudinary.Upload.Upload(ctx, pictureBuffer, uploader.UploadParams{
		PublicID:       filename,
		UniqueFilename: api.Bool(false),
		Folder:         fmt.Sprintf("%s/%s", folderBase, folder),
		Overwrite:      api.Bool(false),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func (c *Cloudinary) UploadOneWithGeneratedFilename(folder string, file *multipart.FileHeader) (string, error) {
	return c.UploadOneWithFilename(folder, file, util.GenerateUUID())
}

func (c *Cloudinary) UploadManyWithGeneratedFilename(folder string, files []*multipart.FileHeader) ([]string, error) {
	var URLs []string

	for _, file := range files {
		URL, err := c.UploadOneWithFilename(folder, file, util.GenerateUUID())
		if err != nil {
			return nil, err
		}

		URLs = append(URLs, URL)
	}

	return URLs, nil
}

func (c *Cloudinary) RenameOneByFilename(folder string, oldFilename string, newFilename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	resp, err := c.cloudinary.Upload.Rename(ctx, uploader.RenameParams{
		FromPublicID: fmt.Sprintf("%s/%s/%s", folderBase, folder, oldFilename),
		ToPublicID:   fmt.Sprintf("%s/%s/%s", folderBase, folder, newFilename),
	})
	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}

func (c *Cloudinary) DeleteOneByFilename(folder string, filename string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	_, err := c.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fmt.Sprintf("%s/%s/%s", folderBase, folder, filename),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cloudinary) DeleteOneByURL(folder string, URL string) error {
	filename := util.GetFilenameWithoutExtension(URL)
	return c.DeleteOneByFilename(folder, filename)
}

func (c *Cloudinary) DeleteManyByURL(folder string, URLs []string) error {
	for _, URL := range URLs {
		filename := util.GetFilenameWithoutExtension(URL)

		err := c.DeleteOneByFilename(folder, filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cloudinary) UpdateArrayImage(folder string, imageURLs []string, updateImage []*helper.UpdateImage) ([]string, error) {
	for i := 0; i < len(updateImage); i++ {
		if updateImage[i].IsDelete {
			imageURLs[i] = ""
		} else if updateImage[i].IsChange {
			URL, err := c.UploadOneWithFilename(folder, updateImage[i].Image, util.GenerateUUID())
			if err != nil {
				return nil, err
			}

			imageURLs[i] = URL
		}
	}

	return util.RemoveNilStringInArray(imageURLs), nil
}
