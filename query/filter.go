package query

type RelationMapping struct {
	TableName   string // リレーションテーブル名
	JoinKey     string // 結合キー (usersテーブルとリレーションテーブルの結合条件)
	FilterField string // フィルタ対象となるリレーションテーブルのフィールド
}

// 任意のフィールドと値による汎用的なフィルター
type ByFieldFilter struct {
	Field           string
	Value           string
	RelationMapping map[string]RelationMapping // モデルごとのリレーションマッピング
}

// Apply メソッドでリレーションテーブルと通常のフィルターを区別
func (f *ByFieldFilter) Apply(q *Query) {
	// リレーションされたテーブルのフィルタリングかどうかを判別
	if mapping, ok := f.RelationMapping[f.Field]; ok {
		// リレーション先のカラムに基づいたフィルタリング
		q.Filters[mapping.FilterField] = f.Value
	} else {
		// 通常のフィルタリング
		q.Filters[f.Field] = f.Value
	}
}
