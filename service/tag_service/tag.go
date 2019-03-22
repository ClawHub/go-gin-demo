package tag_service

import (
	"encoding/json"
	"go-gin-demo/models"
	"go-gin-demo/pkg/export"
	"go-gin-demo/pkg/gredis"
	"go-gin-demo/pkg/logging"
	"go-gin-demo/service/cache_service"
	"go.uber.org/zap"
	"io"
	"strconv"
	"time"
	//excelize 最初的 XML 格式文件的一些结构，是通过 tealeg/xlsx 格式文件结构演化而来的
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
)

//标签
type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

//根据name判断标签是否存在
func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

//根据id判断标签是否存在
func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

//增加标签
func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

//编辑标签
func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}

	return models.EditTag(t.ID, data)
}

//删除标签
func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

//获取标签数量
func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

//分页获取标签列表
func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	//判断是否有缓存
	cache := cache_service.Tag{
		State: t.State,

		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagsKey()
	if gredis.Exists(key).Val() > 0 {
		data, err := gredis.Get(key).Bytes()
		if err != nil {
			logging.HTTPLogger.Info("Get Err", zap.Error(err))
		} else {
			//json解密
			_ = json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	//从mysql中获取
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}

	//入缓存
	gredis.Set(key, tags, 3600)

	return tags, nil
}

//导出标签
func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("标签信息")
	if err != nil {
		return "", err
	}

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, title := range titles {
		cell = row.AddCell()
		cell.Value = title
	}

	for _, v := range tags {
		values := []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.CreatedBy,
			strconv.Itoa(v.CreatedOn),
			v.ModifiedBy,
			strconv.Itoa(v.ModifiedOn),
		}

		row = sheet.AddRow()
		for _, value := range values {
			cell = row.AddCell()
			cell.Value = value
		}
	}

	mow := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + mow + export.EXT

	fullPath := export.GetExcelFullPath() + filename
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}

	return filename, nil
}

//导入标签
func (t *Tag) Import(r io.Reader) error {
	file, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := file.GetRows("标签信息")
	for iRow, row := range rows {
		if iRow > 0 {
			var data []string
			for _, cell := range row {
				data = append(data, cell)
			}

			_ = models.AddTag(data[1], 1, data[2])
		}
	}

	return nil
}

//组装
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}

	return maps
}
