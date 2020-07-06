package wrkrcom

import (
	"errors"
	"log"
	"net/http"

	faktory "github.com/contribsys/faktory/client"
	"github.com/danoand/utils"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

var (
	// CntStopAllWorkers directive to halt worker jobs from continuing
	CntStopAllWorkers = "STOP_WORKERS"
	// CntWorkerNameSendSMSMessages sms msg job name
	CntWorkerNameSendSMSMessages = "bv-job-worker-smsmsgs"
	// CntWorkerNameQRCodeGen qrcode generation job name
	CntWorkerNameQRCodeGen = "bv-job-worker-qrcode-gen"
	// CntWorkerNameSMSSnaps sms snapshot job name
	CntWorkerNameSMSSnaps = "bv-job-worker-smssnaps"
	// ErrClientNotSet indicates that a Faktory client has not been set
	ErrClientNotSet = errors.New("Faktory client not set")

	cntWorkerSMSMsgURL   = "https://bv-job-worker-dev-smsmsgs.herokuapp.com"
	cntWorkerQRCodeURL   = "https://bv-job-worker-dev-qrcode-gen.herokuapp.com"
	cntWorkerSMSSnapsURL = "https://bv-job-worker-dev-smssnaps.herokuapp.com"
)

// Domain houses data provided to a handler
type Domain struct {
	Mgo          *mgo.Session
	MgoDB        *mgo.Database
	FakClientOK  bool
	FakClient    *faktory.Client
	FakClientErr error
}

// NewDomain constructs a domain object
func NewDomain() (Domain, error) {
	var rtn Domain

	rtn.FakClient, rtn.FakClientErr = faktory.Open()
	if rtn.FakClientErr == nil {
		// successfully created a Faktory client
		rtn.FakClientOK = true
	}
	return rtn, rtn.FakClientErr
}

// HndlrStatus returns the status of the Faktory client
func (dmn Domain) HndlrStatus(c *gin.Context) {
	var err error
	var tMap = make(map[string]interface{})
	var rMap = make(map[string]interface{})

	// Has a Faktory client object been generated and
	if !dmn.FakClientOK {
		// Faktory client error has occurred
		rMap["msg"] = "Faktory client error has occurred"
		rMap["err"] = dmn.FakClientErr
		c.JSON(http.StatusInternalServerError, rMap)
		return
	}

	tMap, err = dmn.FakClient.Info()
	if err != nil {
		rMap["msg"] = "error occurred querying the Faktory client"
		rMap["err"] = err.Error()
		c.JSON(http.StatusInternalServerError, rMap)
		return
	}

	rMap["msg"] = "response from the Faktory client"
	rMap["content"] = tMap
	c.JSON(http.StatusOK, rMap)
}

// GenerateHandlrStatus return a gin handler
func (dmn Domain) GenerateHandlrStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var tMap = make(map[string]interface{})
		var rMap = make(map[string]interface{})

		// Has a Faktory client object been generated and
		if !dmn.FakClientOK {
			// Faktory client error has occurred
			rMap["msg"] = "Faktory client error has occurred"
			rMap["err"] = dmn.FakClientErr
			c.JSON(http.StatusInternalServerError, rMap)
			return
		}

		tMap, err = dmn.FakClient.Info()
		if err != nil {
			rMap["msg"] = "error occurred querying the Faktory client"
			rMap["err"] = err.Error()
			c.JSON(http.StatusInternalServerError, rMap)
			return
		}

		rMap["msg"] = "response from the Faktory client"
		rMap["content"] = tMap
		c.JSON(http.StatusOK, rMap)
	}
}

// TODO: commented the lines below due to compilation errors... do I need this function?
// DummyFunc executes a dummy worker function
// func (dmn Domain) DummyFunc(ctx worker.Context, args ...interface{}) error {
//	log.Println("INFO: executing DummyFunc at time:", time.Now().Format(time.RFC3339))
//	return nil
//}

// WakeWorkerApps sends a 'GET' to non-production worker instances to "wake" them up (if sleeping)
func WakeWorkerApps() {
	var (
		err   error
		links = []string{
			cntWorkerQRCodeURL,
			cntWorkerSMSMsgURL,
			cntWorkerSMSSnapsURL,
		}
	)

	// Iterate through the pertinent URLs
	for _, ul := range links {
		// Execute a 'GET'
		_, err = http.Get(ul)
		if err != nil {
			// error fetching a URL (ping a Heroku instance)
			log.Printf("ERROR: %v - error fetching a URL (ping a Heroku instance): %v. See: %v\n",
				utils.FileLine(),
				ul,
				err)
		}
	}
	return
}
