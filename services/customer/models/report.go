package models

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
	// "log"
)

// RecentMonthSignupReport comment
type RecentMonthSignupReport struct {
	Date  null.Time `gorm:"column:date;" sql:"type:datetime" json:"date"`
	Value null.Int  `gorm:"column:value;default:0;" sql:"type:int(11);" json:"value"`
}

// CustomerDetailStatsReport comment
type CustomerDetailStatsReport struct {
	EarnedAchievementID            null.String `gorm:"column:earnedAchievementId;" sql:"type:varchar" json:"earnedAchievementId"`
	EarnedAchievementFulfillmentID null.String `gorm:"column:earnedAchievementFulfillmentId;" sql:"type:varchar" json:"earnedAchievementFulfillmentId"`
}

// ReportRecentMonthSignup ...
func ReportRecentMonthSignup(db *gorm.DB) ([]*map[string]interface{}, error) {
	var (
		err     error
		list    []*RecentMonthSignupReport
		results []*map[string]interface{}
	)

	list = make([]*RecentMonthSignupReport, 0)
	results = make([]*map[string]interface{}, 0)

	db.
		// Debug().
		Table(
			"(?) AS `dates`",
			gorm.Expr(`
SELECT `+"`"+`a`+"`"+`.`+"`"+`date`+"`"+`
FROM (
	SELECT LAST_DAY(CURRENT_DATE) - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS `+"`"+`date`+"`"+`
	FROM (SELECT 0 AS `+"`"+`a`+"`"+` UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS `+"`"+`a`+"`"+`
	CROSS JOIN (SELECT 0 AS `+"`"+`a`+"`"+` UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS `+"`"+`b`+"`"+`
	CROSS JOIN (SELECT 0 AS `+"`"+`a`+"`"+` UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS `+"`"+`c`+"`"+`
) AS `+"`"+`a`+"`"+`
WHERE (
	`+"`"+`a`+"`"+`.`+"`"+`date`+"`"+` BETWEEN DATE_SUB(CURRENT_DATE,INTERVAL 7 DAY) AND CURRENT_DATE
)
ORDER BY `+"`"+`a`+"`"+`.`+"`"+`date`+"`"+`
			`),
		).
		Select("`dates`.`date` AS `date`, COUNT(`customers`.`id`) AS `value`").
		Joins("LEFT JOIN `customers` ON `dates`.`date` = DATE(`customers`.`created_at`)").
		Group("`dates`.`date`").
		Order("`dates`.`date` ASC").
		Find(&list)

	for _, item := range list {
		info, err := FetchDetailInfo(db, item)
		if err != nil {
			break
		}
		results = append(results, &info)
	}

	return results, err
}

// SELECT
//     `dates`.`date` AS `date`,
//     COUNT(customers.id) AS `count`
// FROM
//     (
//         SELECT `a`.`date`
//         FROM (
//             SELECT LAST_DAY(CURRENT_DATE) - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY AS `date`
//             FROM (SELECT 0 AS `a` UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS `a`
//             CROSS JOIN (SELECT 0 AS `a` UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS `b`
//             CROSS JOIN (SELECT 0 AS `a` UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) AS `c`
//         ) AS `a`
//         WHERE (
//             `a`.`date` BETWEEN DATE_ADD(CURRENT_DATE,INTERVAL - DAY(CURRENT_DATE)+1 DAY) AND LAST_DAY(CURRENT_DATE)
//         )
//         ORDER BY `a`.`date`
//     ) AS `dates`
// LEFT JOIN `customers` ON `dates`.`date` = DATE(`customers`.`created_at`)
// GROUP BY `dates`.`date`
// ORDER BY `dates`.`date` ASC

// ReportCustomerDetailStatistics ...
func ReportCustomerDetailStatistics(db *gorm.DB, customerID string) (map[string]interface{}, error) {
	var (
		err     error
		results map[string]interface{} = make(map[string]interface{})
		// wg      sync.WaitGroup
	)

	earnedAchievementIDs := make(chan []string, 0)
	go func(resp chan []string) {
		var (
			earnedAchievementIDs []string = make([]string, 0)
			list []*CustomerDetailStatsReport = make([]*CustomerDetailStatsReport, 0)
		)

		db.
			// Debug().
			Table("`customers` AS `customers`").
			Select("`customers`.`id` AS `customerId`, `earned_achievements`.`achievement_id` AS `earnedAchievementId`").
			Joins("LEFT JOIN `earned_achievements` ON `earned_achievements`.`customer_id` = `customers`.`id`").
			Where("`customers`.`id` = ?", customerID).
			Find(&list)

		for _, item := range list {
			earnedAchievementID := item.EarnedAchievementID.ValueOrZero()
			earnedAchievementIDs = append(earnedAchievementIDs, earnedAchievementID)
		}
		resp <- earnedAchievementIDs
	}(earnedAchievementIDs)

	results["earnedAchievements"] = <-earnedAchievementIDs
	close(earnedAchievementIDs)

	earnedAchievementFulfillmentIDs := make(chan []string, 0)
	go func(resp chan []string) {
		var (
			earnedAchievementFulfillmentIDs []string = make([]string, 0)
			list []*CustomerDetailStatsReport = make([]*CustomerDetailStatsReport, 0)
		)

		db.
			// Debug().
			Table("`customers` AS `customers`").
			Select("`customers`.`id` AS `customerId`, `earned_achievement_fulfillments`.`fulfillment_id` AS `earnedAchievementFulfillmentId`").
			Joins("LEFT JOIN `earned_achievement_fulfillments` ON `earned_achievement_fulfillments`.`customer_id` = `customers`.`id`").
			Where("`customers`.`id` = ?", customerID).
			Find(&list)

		for _, item := range list {
			earnedAchievementFulfillmentID := item.EarnedAchievementFulfillmentID.ValueOrZero()
			earnedAchievementFulfillmentIDs = append(earnedAchievementFulfillmentIDs, earnedAchievementFulfillmentID)
		}
		resp <- earnedAchievementFulfillmentIDs
	}(earnedAchievementFulfillmentIDs)

	results["earnedAchievementFulfillments"] = <-earnedAchievementFulfillmentIDs
	close(earnedAchievementFulfillmentIDs)

	return results, err
}
