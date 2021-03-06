Some command-line cookbook notes and hints, which might be useful when
trying to diagnose cbcollect-info logs...

Once you're in ungzip'ed/untar'ed your cbcollect_info* directory...

    ln -s $GOPATH/src/github.com/couchbase/cbft/cmd/cbft_cbcollect_info_analyze

To dump the full /api/diag from a cbcollect-info directory...

    cat cbcollect_info*/fts_diag.json
    # OR...
    ./cbft_cbcollect_info_analyze extract /api/diag cbcollect_info*

For example, you can now pipe it to "jq ." (or "python -m json.tool") or other tools...

    cat cbcollect_info*/fts_diag.json | jq .
    # OR...
    ./cbft_cbcollect_info_analyze extract /api/diag cbcollect_info* | jq .

To see the "/api/cfg"...

    cat cbcollect_info*/fts_diag.json* | jq '.["/api/cfg"]'
    # OR...
    ./cbft_cbcollect_info_analyze extract /api/diag cbcollect_info* | jq '.["/api/cfg"]'

To see the "/api/stats"...

    cat cbcollect_info*/fts_diag.json | jq '.["/api/stats"]'
    # OR...
    ./cbft_cbcollect_info_analyze extract /api/diag cbcollect_info* | jq '.["/api/stats"]'

To see the bucketDataSourceStats from cbdatasource across all feeds...

    cat cbcollect_info*/fts_diag.json | jq '.["/api/stats"]?["feeds"]?[]?["bucketDataSourceStats"]?'
    # OR...
    ./cbft_cbcollect_info_analyze extract /api/diag cbcollect_info* | jq '.["/api/stats"]?["feeds"]?[]?["bucketDataSourceStats"]?'

To look for non-zero cbdatasource error counts...

    cat cbcollect_info*/fts_diag.json | jq '.["/api/stats"]?["feeds"]?[]?["bucketDataSourceStats"]?' | grep Err | grep -v " 0,"
    # OR...
    ./cbft_cbcollect_info_analyze extract /api/diag cbcollect_info* | jq '.["/api/stats"]?["feeds"]?[]?["bucketDataSourceStats"]?' | grep Err | grep -v " 0,"
    ...
    "TotWorkerReceiveErr": 1,
    "TotWorkerHandleRecvErr": 1,
    "TotUPRDataChangeStateErr": 1,
    "TotUPRSnapshotStateErr": 1,
    "TotWantCloseRequestedVBucketErr": 1,
    ...

To get the cbft node uuids...

    grep "\-uuid" cbcollect_info_*/*fts.log

To find the node that became the so-called "MCP" during rebalance...

    grep %%% cbc*/*.log | cut -f 1 -d \- | sort | uniq

To see what node was being removed during rebalance...

    grep nodesToRemove cbcollect_info_*/*fts.log

To see what metakv "planPIndexes" updates the cbfts might be making...

    grep metakv */*.log | grep fts | grep PUT | grep planPIndexes

Another, faster way to see the what pindex updates the cbft nodes are making through metakv from Aliaksey...

    grep PUT.*PIndexes */ns_server.http_access_internal.log | sort -k4

To find the couchbase nodes that ns-server might know about...

    grep -h per_node_ets_tables */*.log | cut -f 1 -d , | sort | uniq

To see when rebalance was started by ns-server...

    grep -i ns_rebalancer */*.log | grep -i started
    ...
    cbcollect_info_ns_1@172.23.106.176_20160429-001333/diag.log:2016-04-28T17:07:18.198-07:00, ns_rebalancer:0:info:message(ns_1@172.23.106.139) - Started rebalancing bucket default
    ...

To find the num_bytes_used_ram by FTS...

    mortimint -emitParts=NAME cbcollect_info_* | grep num_bytes_used_ram

To see when rebalance REST API call went to ns-server...

    grep rebalance cbc*/ns_server.http_access.log | grep POST

To get the goroutine dump out of our fts_diag.json in original formatting...

    cat fts_diag.json | jq -r '.["/debug/pprof/goroutine?debug=2"]'

    # The -r is the magic to get it back out of a json string,
    # removing quoting, and evaluate newlines again.

