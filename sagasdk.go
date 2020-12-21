package sagasdk

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	http2 "net/http"

	"github.com/ontio/sagasdk/http"
)

const (
	SUCCESS       = 0
	HASHLENERROR  = 1
	HASHDATAERROR = 2
	SERVERERROR   = 3
)

var CodeMessageMap = map[int32]string{
	SUCCESS:       "success",
	HASHLENERROR:  "hash use sha256. match 32 bytes",
	HASHDATAERROR: "hash error",
	SERVERERROR:   "verify server get error from layer2 server",
}

type VerifyResponse struct {
	Code    int32         `json:"code"`
	Message string        `json:"message"`
	Result  *VerifyResult `json:"result"`
}

/*
// for VerifyResult Code.
const (
	SUCCESS          = 0
	PROCESSING       = 1
	FAILED           = 2
	NORECORDONLAYER2 = 3
)
*/

type VerifyResult struct {
	Code             int32  `json:"code"`
	FailedMsg        string `json:"failedMsg"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Proof            string `json:"proof"`
	Layer2Height     uint32 `json:"layer2Height"`
	CommitHeight     uint32 `json:"commitHeight"`
	WitnessStateRoot string `json:"witnessStateRoot"`
	WitnessContract  string `json:"witnessContract"`
}

type SagaSdk struct {
	Client    *http.Client
	VerifyUrl string
}

func NewSagaSdk(verifyUrl string) *SagaSdk {
	http.NewClient()
	return &SagaSdk{
		Client:    http.NewClient(),
		VerifyUrl: verifyUrl,
	}
}

func (self *SagaSdk) VerifyHash(hash string) (*VerifyResult, error) {
	data, respCode, err := self.Client.GetWithHeader(self.VerifyUrl, nil)
	if err != nil || respCode != http2.StatusOK {
		return nil, fmt.Errorf("%d, %s", respCode, err)
	}

	rsp := &VerifyResponse{}
	err = json.Unmarshal(data, rsp)
	if err != nil {
		return nil, err
	}

	if rsp.Code != SUCCESS {
		return nil, fmt.Errorf("%s", rsp.Message)
	}

	return rsp.Result, nil
}

func (self *SagaSdk) AbstractToHash(message []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(message))
}

func (self *SagaSdk) PostRequest(url string, headers []*http.ApiHeadValues, bodyParam []byte) ([]byte, error) {
	data, code, err := self.Client.PostWithHeader(url, headers, bodyParam)
	if code != http2.StatusOK || err != nil {
		return nil, fmt.Errorf("%d, %s", code, err)
	}

	return data, nil
}
