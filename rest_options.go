//  Copyright (c) 2017 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package cbft

import (
	"net/http"

	"github.com/couchbase/cbgt"
	"github.com/couchbase/cbgt/rest"
	log "github.com/couchbase/clog"
)

// List of log levels that maps strings to integers.
// (Works with clog - https://github.com/couchbase/clog)
var logLevels map[string]uint32

func init() {
	logLevels = make(map[string]uint32)

	logLevels["INFO"] = 0
	logLevels["CRIT"] = 3
	logLevels["ERRO"] = 2
	logLevels["FATA"] = 3
	logLevels["WARN"] = 1
}

// ManagerOptionsExt is a REST handler that serves as a wrapper for
// ManagerOptions - where it sets the manager options, and updates
// the logLevel upon request.
type ManagerOptionsExt struct {
	mgr        *cbgt.Manager
	mgrOptions *rest.ManagerOptions
}

func NewManagerOptionsExt(mgr *cbgt.Manager) *ManagerOptionsExt {
	return &ManagerOptionsExt{
		mgr:        mgr,
		mgrOptions: rest.NewManagerOptions(mgr),
	}
}

func (h *ManagerOptionsExt) ServeHTTP(
	w http.ResponseWriter, req *http.Request) {
	h.mgrOptions.ServeHTTP(w, req)

	// Update log level if requested
	logLevel := h.mgr.Options()["logLevel"]
	if len(logLevel) > 0 {
		logLevelInt, exists := logLevels[logLevel]
		if exists {
			log.SetLevel(log.LogLevel(logLevelInt))
		} else {
			log.Warnf("Unrecognized log level setting: %v", logLevel)
		}
	}
}
