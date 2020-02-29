package home

import (
	"math"
	"strconv"
	"strings"

	"github.com/i2eco/ecology/appgo/dao"
	"github.com/i2eco/ecology/appgo/model/constx"
	"github.com/i2eco/ecology/appgo/model/mysql"
	"github.com/i2eco/ecology/appgo/pkg/conf"
	"github.com/i2eco/ecology/appgo/pkg/utils"
	"github.com/i2eco/ecology/appgo/router/core"
)

func index(c *core.Context, bookType string) {
	req := ReqHome{}
	err := c.Bind(&req)
	if err != nil {
		// 首页就用默认值
	}
	cid := req.Cid //分类，如果只是一级分类，则忽略，二级分类，则根据二级分类查找内容
	pageIndex := req.Page
	if pageIndex == 0 {
		pageIndex = 1
	}
	//tab
	var (
		urlPrefix = "/"
		cate      mysql.Category
		tabName   = map[string]string{"recommend": "站长推荐", "latest": "最新发布", "popular": "热门书籍"}
	)
	req.Tab = strings.ToLower(req.Tab)
	switch req.Tab {
	case "recommend", "popular", "latest":
	default:
		req.Tab = "latest"
	}

	cates, _ := dao.Category.GetCates(c.Context, -1, 1)
	pid := cid
	if cid > 0 {
		for _, item := range cates {
			if item.Id == cid {
				if item.Pid > 0 {
					pid = item.Pid
				}
				c.Tpl().Data["Cate"] = item
				cate = item
				break
			}
		}
	}

	c.Tpl().Data["Cates"] = cates
	c.Tpl().Data["Cid"] = cid
	c.Tpl().Data["Pid"] = pid
	c.Tpl().Data["IsHome"] = true

	//每页显示24个，为了兼容Pad、mobile、PC
	pageSize := 24
	books, totalCount, err := dao.Book.HomeData(pageIndex, pageSize, mysql.BookOrder(req.Tab), bookType, cid)
	if err != nil {
		c.Html404()
		return
	}
	if totalCount > 0 {
		urlSuffix := "&tab=" + req.Tab
		if cid > 0 {
			urlSuffix = urlSuffix + "&cid=" + strconv.Itoa(cid)
		}
		urlSuffix = urlSuffix + "&type=" + bookType
		html := utils.NewPaginations(conf.RollPage, totalCount, pageSize, pageIndex, urlPrefix, urlSuffix)
		c.Tpl().Data["PageHtml"] = html
	} else {
		c.Tpl().Data["PageHtml"] = ""
	}

	c.Tpl().Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))

	for _, book := range books {
		book.DealAll()
	}

	c.Tpl().Data["Lists"] = books
	c.Tpl().Data["Tab"] = req.Tab
	c.Tpl().Data["BookType"] = bookType
	title := dao.Global.Get(constx.SITE_NAME)

	if cid > 0 {
		title = "[发现] " + cate.Title + " - " + tabName[req.Tab] + " - " + title
	} else {
		title = "探索，发现新世界，畅想新知识 - " + dao.Global.Get(constx.SITE_NAME)
	}

	c.Html("home/index")
}

func Home(c *core.Context) {
	index(c, "")
}

func Ecology(c *core.Context) {
	index(c, "ecology")
}

func Original(c *core.Context) {
	index(c, "original")
}

func Opensource(c *core.Context) {
	index(c, "opensource")
}
