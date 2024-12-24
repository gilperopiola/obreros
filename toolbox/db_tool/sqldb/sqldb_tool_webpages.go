package sqldb

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/obreros/core/errs"
	"github.com/gilperopiola/obreros/core/models"
)

var (
	CreateWebpageErr = errs.DBCreatingWebpage
	CountWebpagesErr = errs.DBCountingWebpages
	NoOptionsErr     = errs.DBNoQueryOpts
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - SQL DB Tool: Webpage -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (sdbt *sqlDBTool) InsertWebpage(ctx god.Ctx, url, title, content string) (*models.Webpage, error) {
	webpage := models.Webpage{
		URL:     url,
		Title:   title,
		Content: content,
	}

	var previousVersions int64
	query := sdbt.DB.Model(&models.Webpage{}).WithContext(ctx)
	query = query.Where("url = ?", url)
	if err := query.Count(&previousVersions).Error(); err != nil {
		return nil, &errs.DBErr{err, CountWebpagesErr}
	}

	if previousVersions > 0 {
		webpage.Version = int(previousVersions) // The first version is 0.
	}

	if err := sdbt.DB.WithContext(ctx).Create(&webpage).Error(); err != nil {
		return nil, &errs.DBErr{err, CreateWebpageErr}
	}

	return &webpage, nil
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
