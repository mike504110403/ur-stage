package transection

type (
	TransectionSetReq struct {
		MemberId           int     `json:"member_id"`
		Amount             float64 `json:"amount"`
		TransectionSrcType int     `json:"transection_src_type"`
		TransectionRelate  int     `json:"transection_relate"`
	}
)
