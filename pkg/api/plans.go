package api

import (
	"fmt"
	"net/http"
)

type Plan struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	CPU          float64 `json:"cpu"`
	Memory       float64 `json:"memory"`
	Storage      float64 `json:"storage"`
	Available    bool    `json:"available"`
	Sort         int     `json:"sort"`
	HourlyPrice  int     `json:"hourly_price"`
	MonthlyPrice int     `json:"monthly_price"`
}

func (c *Client) PlansList() ([]*Plan, error) {
	url := fmt.Sprintf("plans")

	req, err := c.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	v := make([]*Plan, 0)
	_, err = c.Do(req, &v)
	if err != nil {
		return nil, err
	}

	return v, err
}
