package services

import (
	"thera-api/models"
	"thera-api/repositories"
	"time"
)

type SchedulesService struct {
	SchedulesRepo *repositories.SchedulesRepository
}

func (s *SchedulesService) GetAllSchedules(tenantId string) ([]models.Schedules, error) {
	return s.SchedulesRepo.FindAll(tenantId)
}

func (s *SchedulesService) GetScheduleByID(id string, tenantId string) (*models.Schedules, error) {
	return s.SchedulesRepo.FindByID(id, tenantId)
}

func (s *SchedulesService) CreateSchedule(dateTime time.Time, categoryId, status string, tenantId string) (*models.Schedules, error) {
	schedule := &models.Schedules{
		DateTime:   dateTime,
		CategoryId: categoryId,
		Status:     status,
		TenantId:   &tenantId,
	}

	if err := s.SchedulesRepo.Create(schedule); err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *SchedulesService) UpdateSchedule(
	id string,
	dateTime *time.Time,
	categoryId *string,
	status *string,
	tenantId string,
) (*models.Schedules, error) {
	schedule, err := s.SchedulesRepo.FindByID(id, tenantId)
	if err != nil {
		return nil, err
	}

	if dateTime != nil {
		schedule.DateTime = *dateTime
	}
	if categoryId != nil {
		schedule.CategoryId = *categoryId
	}
	if status != nil {
		schedule.Status = *status
	}

	if err := s.SchedulesRepo.Update(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *SchedulesService) DeleteSchedule(id string, tenantId string) error {
	return s.SchedulesRepo.Delete(id, tenantId)
}
