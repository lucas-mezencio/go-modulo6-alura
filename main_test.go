package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/guilhermeonrails/api-go-gin/controllers"
	"github.com/guilhermeonrails/api-go-gin/database"
	"github.com/guilhermeonrails/api-go-gin/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var ID int

func SetupRotasDeTeste() *gin.Engine {
	// simplifica as mensagens de test mudando para o modo de release do gin
	gin.SetMode(gin.ReleaseMode)

	rotas := gin.Default()
	return rotas
}

func CriaAlunoMock() {
	aluno := models.Aluno{
		Nome: "Nome aluno teste",
		CPF:  "12345678910",
		RG:   "123456789",
	}
	database.DB.Create(&aluno)
	ID = int(aluno.ID)
}

func DeletaAlunoMock() {
	var aluno models.Aluno
	database.DB.Delete(&aluno, ID)
}

func TestVerificaStatusCodeSaudacao(t *testing.T) {
	r := SetupRotasDeTeste()
	r.GET("/:nome", controllers.Saudacao)

	req, _ := http.NewRequest("GET", "/lucas", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	mockResposta := `{"API diz":"E ai lucas, tudo beleza?"}`
	resBody, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, mockResposta, string(resBody))
}

func TestListandoTodosAlunosHandler(t *testing.T) {
	database.ConectaComBancoDeDados()

	CriaAlunoMock()
	defer DeletaAlunoMock()

	r := SetupRotasDeTeste()
	r.GET("/alunos", controllers.ExibeTodosAlunos)

	req, _ := http.NewRequest("GET", "/alunos", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestBuscaAlunoPorCPFHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()

	r := SetupRotasDeTeste()
	r.GET("/alunos/cpf/:cpf", controllers.BuscaAlunoPorCPF)

	req, _ := http.NewRequest("GET", "/alunos/cpf/12345678910", nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

// testa comparando com o corpo da requisição
func TestFindAlunoByIdHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()

	r := SetupRotasDeTeste()
	r.GET("/alunos/:id", controllers.BuscaAlunoPorID)

	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("GET", path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	var alunoMock models.Aluno

	//converte todos os bytes para json e armazena no destino
	_ = json.Unmarshal(res.Body.Bytes(), &alunoMock)
	assert.Equal(t, "Nome aluno teste", alunoMock.Nome)
}

func TestDeletaAlunoHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	r := SetupRotasDeTeste()
	r.DELETE("/alunos/:id", controllers.DeletaAluno)

	path := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", path, nil)
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}

func TestEditaAlunoHandler(t *testing.T) {
	database.ConectaComBancoDeDados()
	CriaAlunoMock()
	defer DeletaAlunoMock()

	r := SetupRotasDeTeste()
	r.PATCH("/alunos/:id", controllers.EditaAluno)

	aluno := models.Aluno{
		Nome: "Nome aluno teste",
		CPF:  "02345678910",
		RG:   "123456460",
	}
	alunoJson, _ := json.Marshal(aluno)
	editPath := "/alunos/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("PATCH", editPath, bytes.NewBuffer(alunoJson))
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	var alunoMock models.Aluno
	json.Unmarshal(res.Body.Bytes(), &alunoMock)
	assert.Equal(t, "02345678910", alunoMock.CPF)
	assert.Equal(t, "123456460", alunoMock.RG)
	assert.Equal(t, "Nome aluno teste", alunoMock.Nome)
}
