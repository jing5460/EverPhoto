package model

import "fmt"

type (
	CommonResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	AuthResponse struct {
		CommonResponse
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}

	SpaceReward int64

	CheckinResult struct {
		Success        bool        `json:"checkin_result"`
		Continuity     int         `json:"continuity"`
		TotalReward    SpaceReward `json:"total_reward"`
		TomorrowReward SpaceReward `json:"tomorrow_reward"`
	}

	CheckinResponse struct {
		CommonResponse
		Data *CheckinResult `json:"data"`
	}
)

const MB = 1 << 20

func (sr SpaceReward) String() string {
	return fmt.Sprintf("%dMB", sr/MB)
}
