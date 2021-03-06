package article

import (
	"strconv"
	"time"

	"github.com/bxcodec/go-clean-arch-grpc/models"
	"github.com/bxcodec/go-clean-arch-grpc/repository"
	"github.com/bxcodec/go-clean-arch-grpc/usecase"
)

type articleUsecase struct {
	articleRepos repository.ArticleRepository
}

func (a *articleUsecase) Fetch(cursor string, num int64) ([]*models.Article, string, error) {
	if num == 0 {
		num = 10
	}

	listArticle, err := a.articleRepos.Fetch(cursor, num)
	if err != nil {
		return nil, "", err
	}
	nextCursor := ""

	if size := len(listArticle); size == int(num) {
		lastId := listArticle[num-1].ID
		nextCursor = strconv.Itoa(int(lastId))
	}

	return listArticle, nextCursor, nil
}

func (a *articleUsecase) GetByID(id int64) (*models.Article, error) {

	return a.articleRepos.GetByID(id)
}

func (a *articleUsecase) Update(ar *models.Article) (*models.Article, error) {
	_, err := a.articleRepos.GetByID(ar.ID)
	if err != nil {
		return nil, err
	}

	ar.UpdatedAt = time.Now()
	return a.articleRepos.Update(ar)
}

func (a *articleUsecase) GetByTitle(title string) (*models.Article, error) {

	return a.articleRepos.GetByTitle(title)
}

func (a *articleUsecase) Store(m *models.Article) (*models.Article, error) {

	existedArticle, _ := a.GetByTitle(m.Title)
	if existedArticle != nil {
		return nil, models.NewErrorConflict()
	}

	id, err := a.articleRepos.Store(m)
	if err != nil {
		return nil, err
	}

	m.ID = id
	return m, nil
}

func (a *articleUsecase) Delete(id int64) (bool, error) {
	existedArticle, _ := a.GetByID(id)

	if existedArticle == nil {
		return false, models.NewErrorNotFound()
	}

	return a.articleRepos.Delete(id)
}

func NewArticleUsecase(a repository.ArticleRepository) usecase.ArticleUsecase {
	return &articleUsecase{a}
}
