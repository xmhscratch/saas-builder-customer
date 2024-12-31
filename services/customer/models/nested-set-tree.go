package models

import (
	null "gopkg.in/guregu/null.v4"
)

// NestedSetTreeModel ...
type NestedSetTreeModel struct {
	NodeLeft  null.Int `gorm:"column:node_left" sql:"type:smallint(6)" json:"nodeLeft"`
	NodeRight null.Int `gorm:"column:node_right" sql:"type:smallint(6)" json:"nodeRight"`
	NodeLevel null.Int `gorm:"column:node_level" sql:"type:smallint(6)" json:"nodeLevel"`
	NodeDepth null.Int `gorm:"column:node_depth" sql:"type:smallint(6)" json:"nodeDepth"`
}
