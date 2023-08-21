package service

import (
	"context"
	"log"
	"shellrean.id/belajar-auth/domain"
	"shellrean.id/belajar-auth/dto"
	"shellrean.id/belajar-auth/internal/util"
	"time"
)

type fdsService struct {
	ipCheckerService   domain.IpCheckerService
	loginLogRepository domain.LoginLogRepository
}

func NewFds(ipCheckerService domain.IpCheckerService,
	loginLogRepository domain.LoginLogRepository) domain.FdsService {
	return &fdsService{
		ipCheckerService:   ipCheckerService,
		loginLogRepository: loginLogRepository,
	}
}

func (f fdsService) IsAuthorized(ctx context.Context, ip string, userId int64) bool {
	locationCheck, err := f.ipCheckerService.Query(ctx, ip)
	if err != nil || locationCheck == (dto.IpChecker{}) {
		return false
	}

	newAccess := domain.LoginLog{
		UserID:       userId,
		IsAuthorized: false,
		IpAddress:    ip,
		Timezone:     locationCheck.Timezone,
		Lon:          locationCheck.Lon,
		Lat:          locationCheck.Lat,
		AccessTime:   time.Now(),
	}

	lastLogin, err := f.loginLogRepository.FindLastAuthorized(ctx, userId)
	if err != nil {
		_ = f.loginLogRepository.Save(ctx, &newAccess)
		return false
	}
	if lastLogin == (domain.LoginLog{}) {
		newAccess.IsAuthorized = true
		_ = f.loginLogRepository.Save(ctx, &newAccess)
		return true
	}

	distanceHour := newAccess.AccessTime.Sub(lastLogin.AccessTime)
	distanceChange := util.GetDistance(lastLogin.Lat, lastLogin.Lon, newAccess.Lat, newAccess.Lon)

	log.Printf("hour: %f, distance: %f\n", distanceHour.Hours(), distanceChange)

	if (distanceChange / distanceHour.Hours()) > 400 {
		_ = f.loginLogRepository.Save(ctx, &newAccess)
		return false
	}

	newAccess.IsAuthorized = true
	_ = f.loginLogRepository.Save(ctx, &newAccess)
	return true
}
