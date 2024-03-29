package controller

import (
	"fmt"
	"gin_study_blog/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	uploadConfig map[string]interface{}
)

func IndexHandle(c *gin.Context) {

	articleRecordList, err := service.GetArticleRecordList(0, 15)
	for _, v := range articleRecordList {
		fmt.Printf("articlerecord:%#v", v)
	}
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	categoryList, err := service.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get category list failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	var data map[string]interface{} = make(map[string]interface{}, 10)
	data["article_list"] = articleRecordList
	data["category_list"] = categoryList

	c.HTML(http.StatusOK, "views/index.html", data)
}

func CategoryList(c *gin.Context) {
	categoryIdStr := c.Query("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		fmt.Printf("invalid parameter, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	articleCategoryList, err := service.GetCategoryArticle(categoryId)
	if err != nil {
		fmt.Printf("get category article list failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	categoryList, err := service.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get category list failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	var data map[string]interface{} = make(map[string]interface{}, 10)
	data["article_list"] = articleCategoryList
	data["category_list"] = categoryList
	c.HTML(http.StatusOK, "views/index.html", data)
}

/*
func CategoryList(c *gin.Context) {

	categoryIdStr := c.Query("category_id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	articleRecordList, err := logic.GetArticleRecordListById(int(categoryId), 0, 15)
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	allCategoryList, err := logic.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get category list failed, err:%v\n", err)
	}

	var data map[string]interface{} = make(map[string]interface{}, 10)
	data["article_list"] = articleRecordList
	data["category_list"] = allCategoryList

	c.HTML(http.StatusOK, "views/index.html", data)
}
*/

func NewArticle(c *gin.Context) {
	categoryList, err := service.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get article failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/post_article.html", categoryList)
}

func ArticleSubmit(c *gin.Context) {
	content := c.PostForm("content")
	author := c.PostForm("author")
	categoryIdStr := c.PostForm("category_id")
	title := c.PostForm("title")

	// 转成10进制，64位的整数
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	err = service.InsertArticle(content, author, title, categoryId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	// 状态码 200 的时候不能重定向，所以用 301
	c.Redirect(http.StatusMovedPermanently, "/")
}

func ArticleDetail(c *gin.Context) {

	// 文章 id 是通过 QueryString 提交过来的
	articleIdStr := c.Query("article_id")
	// 转成整数
	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	// 获取文章信息
	articleDetail, err := service.GetArticleDetail(articleId)
	if err != nil {
		fmt.Printf("get article detail failed,article_id:%d err:%v\n", articleId, err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	fmt.Printf("article detail:%#v\n", articleDetail)

	relativeArticle, err := service.GetRelativeArticleList(articleId)
	if err != nil {
		fmt.Printf("get relative article failed, err:%v\n", err)
	}

	prevArticle, nextArticle, err := service.GetPrevAndNextArticleInfo(articleId)
	if err != nil {
		fmt.Printf("get prev or next article failed, err:%v\n", err)
	}

	allCategoryList, err := service.GetAllCategoryList()
	if err != nil {
		fmt.Printf("get all category failed, err:%v\n", err)
	}

	commentList, err := service.GetCommentList(articleId)
	if err != nil {
		fmt.Printf("get comment list failed, err:%v\n", err)
	}

	fmt.Printf("relative article size:%d article_id:%d\n", len(relativeArticle), articleId)
	var m map[string]interface{} = make(map[string]interface{}, 10)
	m["detail"] = articleDetail
	m["relative_article"] = relativeArticle
	m["prev"] = prevArticle
	m["next"] = nextArticle
	m["category"] = allCategoryList
	m["article_id"] = articleId
	m["comment_list"] = commentList

	c.HTML(http.StatusOK, "views/detail.html", m)
}

func CommentSubmit(c *gin.Context) {
	//
	author := c.PostForm("author")
	comment := c.PostForm("comment")
	articleIdStr := c.PostForm("article_id")

	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}
	err = service.InsertComment(author, comment, articleId)
	if err != nil {
		fmt.Printf("insert comment failed, err:%v\n", err)
	}

	url := fmt.Sprintf("/article/detail/?article_id=%d", articleId)
	c.Redirect(http.StatusMovedPermanently, url)

}

func LeaveNew(c *gin.Context) {
	//name := c.PostForm("author")
	//content := c.PostForm("comment")
	//email := c.PostForm("email")
	//
	//err := service.InsertLeave(name, content, email)
	//if err != nil {
	//	fmt.Printf("insert leave failed, err:%v\n", err)
	//	return
	//}

	leaveList, err := service.GetLeaveList()
	for _, v := range leaveList {
		fmt.Printf("leave:%#v\n", v)
	}
	if err != nil {
		fmt.Printf("get leave list failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/gbook.html", leaveList)

}

func LeaveSubmit(c *gin.Context) {
	name := c.PostForm("author")
	content := c.PostForm("comment")
	email := c.PostForm("email")

	err := service.InsertLeave(name, content, email)
	if err != nil {
		fmt.Printf("insert leave failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/leave/new/")

}

/*
func LeaveNew(c *gin.Context) {
	leaveList, err := logic.GetLeaveList()
	if err != nil {
		fmt.Printf("get leave failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	c.HTML(http.StatusOK, "views/gbook.html", leaveList)
}

func AboutMe(c *gin.Context) {
	c.HTML(http.StatusOK, "views/about.html", gin.H{
		"title": "Posts",
	})
}





func UploadFile(c *gin.Context) {
	// single file
	file, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	log.Println(file.Filename)
	rootPath := util.GetRootDir()
	u2, err := uuid.NewV4()
	if err != nil {
		return
	}

	ext := path.Ext(file.Filename)
	url := fmt.Sprintf("/static/upload/%s%s", u2, ext)
	dst := filepath.Join(rootPath, url)
	// Upload the file to specific dst.
	c.SaveUploadedFile(file, dst)
	c.JSON(http.StatusOK, gin.H{
		"uploaded": true,
		"url":      url,
	})
}

func CommentSubmit(c *gin.Context) {

	comment := c.PostForm("comment")
	author := c.PostForm("author")
	email := c.PostForm("email")
	articleIdStr := c.PostForm("article_id")

	articleId, err := strconv.ParseInt(articleIdStr, 10, 64)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	err = logic.InsertComment(comment, author, email, articleId)
	if err != nil {
		fmt.Printf("insert comment failed, err:%v\n", err)
		c.HTML(http.StatusInternalServerError, "views/500.html", nil)
		return
	}

	url := fmt.Sprintf("/article/detail/?article_id=%d", articleId)
	c.Redirect(http.StatusMovedPermanently, url)
}




*/
