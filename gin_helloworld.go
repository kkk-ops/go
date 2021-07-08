package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11);not null;unique"`
	Password string `gorm:"size:255;not null"`
}

func main() {
	// 初始化DB
	db :=InitDB()
	defer db.Close() //延迟关闭

	// 创建一个默认的路由引擎
	r := gin.Default()

	r.LoadHTMLGlob("../templates/**/*")//html渲染
	//r.LoadHTMLFiles("./simple-doc.html")
	r.Static("/static", "./static")//静态文件处理

	// 当客户端以GET方法请求路径时，会执行后面的匿名函数
	r.GET("/hello", func(c *gin.Context) {
		// c.JSON：返回JSON格式的数据
		c.JSON(http.StatusOK, gin.H{
			"status":  gin.H{
				"code": http.StatusOK,
				"status":      "登录成功",
			},
			"message": "Hello world!",
		})
	})

	r.POST("/api/auth/register", func(c *gin.Context) {
		//获取参数
		name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")


		// 数据验证
		if len(telephone) != 11 {
			c.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"手机号必须为11位"})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"密码不能少于6位"})
			return
		}
         if len(name) == 0 {
         	name = RandomString(10)
         }
         log.Println(name,telephone,password)
         //判断手机号是否存在
        if isTelephoneExist(db,telephone) {
        	c.JSON(http.StatusUnprocessableEntity,gin.H{"code":422,"msg":"用户已存在"})
        	return

		}
		//创建用户
		newUser :=User{
			Name: name,
			Telephone: telephone,
			Password: password,
		}
		db.Create(&newUser)

		//返回结果
		c.JSON(200,gin.H{
			"msg": "注册成功",
		})
	})

	r.GET("users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.html", gin.H{
			"title": "users/index",
		})
	})
	r.GET("/doc", func(c *gin.Context) {
		c.HTML(http.StatusOK, "simple-doc.html", nil)
	})

	// 处理multipart forms提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	r.POST("/uploadone", func(c *gin.Context) {
		// 单个文件
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		log.Println(file.Filename)
		dst := fmt.Sprintf("D:/src/%s", file.Filename)
		// 上传文件到指定的目录
		c.SaveUploadedFile(file, dst)
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
		})
	})

	// 处理multipart forms提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	r.POST("/uploadmul", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["file"]

		for index, file := range files {
			log.Println(file.Filename)
			dst := fmt.Sprintf("D:/src/%s_%d", file.Filename, index)
			// 上传文件到指定的目录
			c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d files uploaded!", len(files)),
		})
	})

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run()
}

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNM")
	result := make([]byte,n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
return string(result)
}

func isTelephoneExist(db *gorm.DB,telephone string) bool {
	var user User
	db.Where("telephone = ?",telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "password"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		username,
		password,
		host,
		port,
		database,
		charset, )

	db,err :=  gorm.Open(driverName,args)
	if err != nil {
		panic("fail to connect database,err:" + err.Error())
	}
	db.AutoMigrate(&User{})
	return db

}