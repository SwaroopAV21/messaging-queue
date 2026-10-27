package models


type Request struct {
	RequestID	int64		`json:"requestID"`
	Api			string		`json:"requestAPI"`
	Method		string		`json:"requestMethod"`
	ReturnCode  int			`json:"returnCode"`
}


func (req Request) GetRequestBody(requestID int64) (Request, bool)  {
	if req.RequestID == requestID {
		return req, true
	}
	return Request{}, false
}