package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
)

func HandleRequests() {
	r := gin.Default()

	//indica ao gin para rendenizar paginas html
	r.LoadHTMLGlob("templates/*")
	//indica ao gin para carregar recursos staticos
	r.Static("/assets", "./assets")

	r.GET("/:nome", controllers.Saudacao)
	r.GET("/alunos", controllers.ExibeTodosAlunos)
	r.POST("/alunos", controllers.CriaNovoAluno)
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)
	r.DELETE("/alunos/:id", controllers.DeletaAluno)
	r.PATCH("/alunos/:id", controllers.EditaAluno)
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)
	r.GET("/index", controllers.ExibePaginaIndex)
	r.NoRoute(controllers.RotaNaoEncontrada)
	r.Run()
}
