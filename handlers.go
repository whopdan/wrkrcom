package wrkrcom

import (
	"log"
	"net/http"
	"time"

	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

// Domain houses data provided to a handler
type Domain struct {
	Mgo       *mgo.Session
	MgoDB     *mgo.Database
	FakClient *faktory.Client
}

// HndlrStatus returns the status of the Faktory client
func (dmn *Domain) HndlrStatus(c *gin.Context) {
	var err error
	var tMap = make(map[string]interface{})
	var rMap = make(map[string]interface{})

	tMap, err = dmn.FakClient.Info()
	if err != nil {
		rMap["msg"] = "error occurred querying the Faktory client"
		rMap["err"] = err.Error()
		c.JSON(http.StatusInternalServerError, rMap)
		return
	}

	rMap["msg"] = "response from the Faktory client"
	rMap["content"] = tMap
}

// DummyFunc executes a dummy worker function
func (dmn *Domain) DummyFunc(ctx worker.Context, args ...interface{}) error {
	log.Println("INFO: executing DummyFunc at time:", time.Now().Format(time.RFC3339))
	return nil
}
