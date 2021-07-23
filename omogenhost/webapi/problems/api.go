package problems

import (
	"context"
	"errors"
	"github.com/jsannemo/omogenhost/storage"
	apipb "github.com/jsannemo/omogenhost/webapi/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type problemService struct {
}

func InitProblemService() *problemService {
	return &problemService{}
}

var statementOrder = []string{"sv", "en"}

func (ps *problemService) ViewProblem(ctx context.Context, request *apipb.ViewProblemRequest) (*apipb.ViewProblemResponse, error) {
	shortname := request.ShortName
	problem := storage.Problem{}
	if res := storage.GormDB.Debug().Preload("ProblemStatements").Joins("CurrentVersion").First(&problem, "short_name = ?", shortname); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) || len(problem.ProblemStatements) == 0 {
			return nil, status.Error(codes.NotFound, "No such problem")
		}
		return nil, res.Error
	}
	statements := make(map[string]*storage.ProblemStatement)
	for _, statement := range problem.ProblemStatements {
		statements[statement.Language] = &statement
	}
	var statement *storage.ProblemStatement
	if request.Language != "" {
		if s, found := statements[request.Language]; !found {
			return nil, status.Error(codes.NotFound, "No such problem")
		} else {
			statement = s
		}
	} else {
		for _, lang := range statementOrder {
			if s, found := statements[lang]; found {
				statement = s
				break
			}
		}
	}
	if statement == nil {
		statement = &problem.ProblemStatements[0]
	}
	return &apipb.ViewProblemResponse{
		Statement: &apipb.ProblemStatement{
			Language: statement.Language,
			Title:    statement.Title,
			Html:     statement.Html,
			License:  problem.License,
			Authors:  problem.Author,
		},
		Limits: &apipb.ProblemLimits{
			TimeLimitMs:   problem.CurrentVersion.TimeLimitMs,
			MemoryLimitKb: problem.CurrentVersion.MemoryLimitKb,
		},
	}, nil
}
