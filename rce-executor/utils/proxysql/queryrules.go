package proxysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

type (
	//QueryRules 查询规则
	QueryRules struct {
		RuleID               uint64 `db:"rule_id" json:"rule_id"`
		Active               uint64 `db:"active" json:"active"`
		Username             string `db:"username" json:"username"`
		Schemaname           string `db:"schemaname" json:"schemaname"`
		FlagIN               uint64 `db:"flagIN" json:"flagIN"`
		ClientAddr           string `db:"client_addr" json:"client_addr"`
		ProxyAddr            string `db:"proxy_addr" json:"proxy_addr"`
		ProxyPort            string `db:"proxy_port" json:"proxy_port"`
		Digest               string `db:"digest" json:"digest"`
		MatchDigest          string `db:"match_digest" json:"match_digest"`
		MatchPattern         string `db:"match_pattern" json:"match_pattern"`
		NegateMatchPattern   uint64 `db:"negate_match_pattern" json:"negate_match_pattern"`
		FlagOUT              string `db:"flagOUT" json:"flagOUT"`
		ReplacePattern       string `db:"replace_pattern" json:"replace_pattern"`
		DestinationHostgroup string `db:"destination_hostgroup" json:"destination_hostgroup"`
		CacheTTL             string `db:"cache_ttl" json:"cache_ttl"`
		Reconnect            string `db:"reconnect" json:"reconnect"`
		Timeout              string `db:"timeout" json:"timeout"`
		Retries              string `db:"retries" json:"retries"`
		Delay                string `db:"delay" json:"delay"`
		MirrorFlagOUT        string `db:"mirror_flagOUT" json:"mirror_flagOUT"`
		MirrorHostgroup      string `db:"mirror_hostgroup" json:"mirror_hostgroup"`
		ErrorMsg             string `db:"error_msg" json:"error_msg"`
		Log                  string `db:"log" json:"log"`
		Apply                uint64 `db:"apply" json:"apply"`
		Comment              string `db:"comment" json:"comment"`
	}
)

//define query rules
const (
	StmtAddOneQr = `
	INSERT 
	INTO 
		mysql_query_rules(rule_id,username) 
	VALUES(%d,%s)`

	/*delete a query rules*/
	StmtDeleteOneQr = `
	DELETE 
	FROM 
		mysql_query_rules 
	WHERE rule_id = %d`

	/*query all query rules.*/
	StmtFindAllQr = `
	SELECT 
		ifnull(rule_id,0) as rule_id,
		ifnull(active,0) as active,
		ifnull(username,"") as username,
		ifnull(schemaname,"") as schemaname,
		ifnull(flagIN,0) as flagIN,
		ifnull(client_addr,"") as client_addr,
		ifnull(proxy_addr,"") as proxy_addr,
		ifnull(proxy_port,0) as proxy_port,
		ifnull(digest,"") as digest,
		ifnull(match_digest,"") as match_digest,
		ifnull(match_pattern,"") as match_pattern,
		ifnull(negate_match_pattern,0) as negate_match_pattern,
		ifnull(flagOUT,0) as flagOUT,
		ifnull(replace_pattern,"") as replace_pattern,
		ifnull(destination_hostgroup,0) as destination_hostgroup,
		ifnull(cache_ttl,0) as cache_ttl,
		ifnull(reconnect,0) as reconnect,
		ifnull(timeout,0) as timeout,
		ifnull(retries,0) as retries,
		ifnull(delay,0) as delay,
		ifnull(mirror_flagOUT,0) as mirror_flagOUT,
		ifnull(mirror_hostgroup,0) as mirror_hostgroup,
		ifnull(error_msg,"") as error_msg,
		ifnull(log,0) as log,
		ifnull(apply,0) as apply,
		ifnull(comment,"") as comment 
	FROM mysql_query_rules 
	LIMIT %d 
	OFFSET %d`

	/*find last insert rule_id*/
	StmtFindLastRuleID = `
	SELECT 
		max(rule_id)
	FROM mysql_query_rules
	WHERE 
		username = %s`

	/*update a query rules.*/
	StmtUpdateOneQr = `
	UPDATE 
		mysql_query_rules 
	SET 
		active=%d,
		username=%s,
		schemaname=%s,
		flagIN=%d,
		client_addr=%s,
		proxy_addr=%s,
		proxy_port=%s,
		digest=%s,
		match_digest=%s,
		match_pattern=%s,
		negate_match_pattern=%d,
		flagOUT=%s,
		replace_pattern=%s,
		destination_hostgroup=%s,
		cache_ttl=%s,
		reconnect=%s,
		timeout=%s,
		retries=%s,
		delay=%s,
		mirror_flagOUT=%s,
		mirror_hostgroup=%s,
		error_msg=%s,
		log=%s,
		apply=%d,
		comment=%s
	WHERE 
		rule_id=%d`
)

//FindAllQr select * from mysql_query_rules limit n offset n
func FindAllQr(db *sql.DB, limit uint64, skip uint64) ([]QueryRules, error) {

	var AllQr []QueryRules
	Query := fmt.Sprintf(StmtFindAllQr, limit, skip)

	// exec query statement
	rows, err := db.Query(Query)
	if err != nil {
		return []QueryRules{}, errors.Trace(err)
	}
	defer rows.Close()

	// scan results.
	for rows.Next() {

		var tmpqr QueryRules

		err = rows.Scan(
			&tmpqr.RuleID,
			&tmpqr.Active,
			&tmpqr.Username,
			&tmpqr.Schemaname,
			&tmpqr.FlagIN,
			&tmpqr.ClientAddr,
			&tmpqr.ProxyAddr,
			&tmpqr.ProxyPort,
			&tmpqr.Digest,
			&tmpqr.MatchDigest,
			&tmpqr.MatchPattern,
			&tmpqr.NegateMatchPattern,
			&tmpqr.FlagOUT,
			&tmpqr.ReplacePattern,
			&tmpqr.DestinationHostgroup,
			&tmpqr.CacheTTL,
			&tmpqr.Reconnect,
			&tmpqr.Timeout,
			&tmpqr.Retries,
			&tmpqr.Delay,
			&tmpqr.MirrorFlagOUT,
			&tmpqr.MirrorHostgroup,
			&tmpqr.ErrorMsg,
			&tmpqr.Log,
			&tmpqr.Apply,
			&tmpqr.Comment,
		)

		if err != nil {
			continue
		}

		AllQr = append(AllQr, tmpqr)
	}
	return AllQr, nil
}

//NewQr new mysql query rules
func NewQr(username string) (*QueryRules, error) {
	newqr := new(QueryRules)

	if username == "" {
		return nil, errors.BadRequestf(username)
	}
	if strings.Index(username, "\"") == -1 {
		newqr.Username = fmt.Sprintf("\"%s\"", username)
	} else {
		newqr.Username = username
	}

	newqr.DestinationHostgroup = "NULL"
	newqr.Schemaname = "NULL"
	newqr.FlagIN = 0
	newqr.ClientAddr = "NULL"
	newqr.ProxyAddr = "NULL"
	newqr.ProxyPort = "NULL"
	newqr.Digest = "NULL"
	newqr.MatchDigest = "NULL"
	newqr.MatchPattern = "NULL"
	newqr.NegateMatchPattern = 0
	newqr.FlagOUT = "NULL"
	newqr.ReplacePattern = "NULL"
	newqr.CacheTTL = "NULL"
	newqr.Reconnect = "NULL"
	newqr.Timeout = "NULL"
	newqr.Retries = "NULL"
	newqr.Delay = "NULL"
	newqr.MirrorFlagOUT = "NULL"
	newqr.MirrorHostgroup = "NULL"
	newqr.ErrorMsg = "NULL"
	newqr.Log = "NULL"
	newqr.Apply = 0
	newqr.Active = 0
	newqr.Comment = "NULL"

	return newqr, nil
}

//SetQrRuleid set qr rule_id
func (qr *QueryRules) SetQrRuleid(ruleID uint64) {
	qr.RuleID = ruleID
}

//SetQrActive set qr active
func (qr *QueryRules) SetQrActive(active uint64) {
	switch active {
	case 0:
		qr.Active = 0
	case 1:
		qr.Active = 1
	default:
		qr.Active = 1
	}
}

//SetQrApply set qr apply
func (qr *QueryRules) SetQrApply(apply uint64) {
	switch apply {
	case 0:
		qr.Apply = 0
	case 1:
		qr.Apply = 1
	default:
		qr.Apply = 1
	}
}

//SetQrSchemaName set qr schemaname
func (qr *QueryRules) SetQrSchemaName(schemaName string) {
	if schemaName == "" || len(schemaName) == 0 {
		qr.Schemaname = "NULL"
	} else {
		if strings.Index(schemaName, "\"") == -1 {
			qr.Schemaname = fmt.Sprintf("\"%s\"", schemaName)
		} else {
			qr.Schemaname = schemaName
		}
	}
}

//SetQrFlagIN set qr flagIN
func (qr *QueryRules) SetQrFlagIN(flagIN uint64) {
	qr.FlagIN = flagIN
}

//SetQrClientAddr set qr client_addr
func (qr *QueryRules) SetQrClientAddr(clientAddr string) {
	if clientAddr == "" || len(clientAddr) == 0 {
		qr.ClientAddr = "NULL"
	} else {
		if strings.Index(clientAddr, "\"") == -1 {
			qr.ClientAddr = fmt.Sprintf("\"%s\"", clientAddr)
		} else {
			qr.ClientAddr = clientAddr
		}
	}
}

//SetQrProxyAddr set qr proxy_addr
func (qr *QueryRules) SetQrProxyAddr(proxyAddr string) {
	if proxyAddr == "" || len(proxyAddr) == 0 {
		qr.ProxyAddr = "NULL"
	} else {
		if strings.Index(proxyAddr, "\"") == -1 {
			qr.ProxyAddr = fmt.Sprintf("\"%s\"", proxyAddr)
		} else {
			qr.ProxyAddr = proxyAddr
		}
	}
}

//SetProxyPort set qr proxy_port
func (qr *QueryRules) SetProxyPort(proxyPort string) {
	if proxyPort == "" || len(proxyPort) == 0 {
		qr.ProxyPort = "NULL"
	} else {
		if strings.Index(proxyPort, "\"") == -1 {
			qr.ProxyPort = fmt.Sprintf("\"%s\"", proxyPort)
		} else {
			qr.ProxyPort = proxyPort
		}
	}
}

//SetQrDigest set qr digest
func (qr *QueryRules) SetQrDigest(digest string) {
	if digest == "" || len(digest) == 0 {
		qr.Digest = "NULL"
	} else {
		if strings.Index(digest, "\"") == -1 {
			qr.Digest = fmt.Sprintf("\"%s\"", digest)
		} else {
			qr.Digest = digest
		}
	}
}

//SetQrMatchDigest set qr match_digest
func (qr *QueryRules) SetQrMatchDigest(matchDigest string) {
	if matchDigest == "" || len(matchDigest) == 0 {
		qr.MatchDigest = "NULL"
	} else {
		if strings.Index(matchDigest, "\"") == -1 {
			qr.MatchDigest = fmt.Sprintf("\"%s\"", matchDigest)
		} else {
			qr.MatchDigest = matchDigest
		}
	}
}

//SetQrMatchPattern set qr match_pattern
func (qr *QueryRules) SetQrMatchPattern(matchPattern string) {
	if matchPattern == "" || len(matchPattern) == 0 {
		qr.MatchPattern = "NULL"
	} else {
		if strings.Index(matchPattern, "\"") == -1 {
			qr.MatchPattern = fmt.Sprintf("\"%s\"", matchPattern)
		} else {
			qr.MatchPattern = matchPattern
		}
	}
}

//SetQrNegateMatchPattern set qr mnegate_match_pattern
func (qr *QueryRules) SetQrNegateMatchPattern(negateMatchPattern uint64) {
	switch negateMatchPattern {
	case 0:
		qr.NegateMatchPattern = 0
	case 1:
		qr.NegateMatchPattern = 1
	default:
		qr.NegateMatchPattern = 0
	}
}

//SetQrFlagOut set qr flagout
func (qr *QueryRules) SetQrFlagOut(flagOut string) {
	if flagOut == "" || len(flagOut) == 0 {
		qr.FlagOUT = "NULL"
	} else {
		if strings.Index(flagOut, "\"") == -1 {
			qr.FlagOUT = fmt.Sprintf("\"%s\"", flagOut)
		} else {
			qr.FlagOUT = flagOut
		}
	}
}

//SetQrReplacePattern  set qr replace_pattern
func (qr *QueryRules) SetQrReplacePattern(replacePattern string) {
	if replacePattern == "" || len(replacePattern) == 0 {
		qr.ReplacePattern = "NULL"
	} else {
		if strings.Index(replacePattern, "\"") == -1 {
			qr.ReplacePattern = fmt.Sprintf("\"%s\"", replacePattern)
		} else {
			qr.ReplacePattern = replacePattern
		}
	}
}

//SetQrDestHostGroup set qr destination_hostgroup
func (qr *QueryRules) SetQrDestHostGroup(destinationHostgroup string) {
	if destinationHostgroup == "" || len(destinationHostgroup) == 0 {
		qr.DestinationHostgroup = "NULL"
	} else {
		if strings.Index(destinationHostgroup, "\"") == -1 {
			qr.DestinationHostgroup = fmt.Sprintf("\"%s\"", destinationHostgroup)
		} else {
			qr.DestinationHostgroup = destinationHostgroup
		}
	}
}

//SetQrCacheTTL set qr cache_ttl
func (qr *QueryRules) SetQrCacheTTL(cacheTTL string) {
	if cacheTTL == "" || len(cacheTTL) == 0 {
		qr.CacheTTL = "NULL"
	} else {
		if strings.Index(cacheTTL, "\"") == -1 {
			qr.CacheTTL = fmt.Sprintf("\"%s\"", cacheTTL)
		} else {
			qr.CacheTTL = cacheTTL
		}
	}
}

//SetQrReconnect set qr reconnect
func (qr *QueryRules) SetQrReconnect(reconnect string) {
	if reconnect == "" || len(reconnect) == 0 {
		qr.Reconnect = "NULL"
	} else {
		if strings.Index(reconnect, "\"") == -1 {
			qr.Reconnect = fmt.Sprintf("\"%s\"", reconnect)
		} else {
			qr.Reconnect = reconnect
		}
	}
}

//SetQrTimeOut  set qr timeout
func (qr *QueryRules) SetQrTimeOut(timeout string) {
	if timeout == "" || len(timeout) == 0 {
		qr.Timeout = "NULL"
	} else {
		if strings.Index(timeout, "\"") == -1 {
			qr.Timeout = fmt.Sprintf("\"%s\"", timeout)
		} else {
			qr.Timeout = timeout
		}
	}
}

//SetQrRetries set qr retries
func (qr *QueryRules) SetQrRetries(retries string) {
	if retries == "" || len(retries) == 0 {
		qr.Retries = "NULL"
	} else {
		if strings.Index(retries, "\"") == -1 {
			qr.Retries = fmt.Sprintf("\"%s\"", retries)
		} else {
			qr.Retries = retries
		}
	}
}

//SetQrDelay set qr delay
func (qr *QueryRules) SetQrDelay(delay string) {
	if delay == "" || len(delay) == 0 {
		qr.Delay = "NULL"
	} else {
		if strings.Index(delay, "\"") == -1 {
			qr.Delay = fmt.Sprintf("\"%s\"", delay)
		} else {
			qr.Delay = delay
		}
	}
}

//SetQrMirrorFlagOUT  set qr mirror_flagout
func (qr *QueryRules) SetQrMirrorFlagOUT(mirrorFlagout string) {
	if mirrorFlagout == "" || len(mirrorFlagout) == 0 {
		qr.MirrorFlagOUT = "NULL"
	} else {
		if strings.Index(mirrorFlagout, "\"") == -1 {
			qr.MirrorFlagOUT = fmt.Sprintf("\"%s\"", mirrorFlagout)
		} else {
			qr.MirrorFlagOUT = mirrorFlagout
		}
	}
}

//SetQrMirrorHostgroup set qr mirror_hostgroup
func (qr *QueryRules) SetQrMirrorHostgroup(mirrorHostgroup string) {
	if mirrorHostgroup == "" || len(mirrorHostgroup) == 0 {
		qr.MirrorHostgroup = "NULL"
	} else {
		if strings.Index(mirrorHostgroup, "\"") == -1 {
			qr.MirrorHostgroup = fmt.Sprintf("\"%s\"", mirrorHostgroup)
		} else {
			qr.MirrorHostgroup = mirrorHostgroup
		}
	}
}

//SetQrErrorMsg set qr error_msg
func (qr *QueryRules) SetQrErrorMsg(errorMsg string) {
	if errorMsg == "" || len(errorMsg) == 0 {
		qr.ErrorMsg = "NULL"
	} else {
		if strings.Index(errorMsg, "\"") == -1 {
			qr.ErrorMsg = fmt.Sprintf("\"%s\"", errorMsg)
		} else {
			qr.ErrorMsg = errorMsg
		}
	}
}

//SetQrLog set qr log
func (qr *QueryRules) SetQrLog(log string) {
	if log == "" || len(log) == 0 {
		qr.Log = "NULL"
	} else {
		if strings.Index(log, "\"") == -1 {
			qr.Log = fmt.Sprintf("\"%s\"", log)
		} else {
			qr.Log = log
		}
	}
}

//AddOneQr add a new query rules.
func (qr *QueryRules) AddOneQr(db *sql.DB) error {

	Query := fmt.Sprintf(StmtAddOneQr, qr.RuleID, qr.Username)

	_, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err) //add user failed
	}

	Query = fmt.Sprintf(StmtFindLastRuleID, qr.Username)
	rows := db.QueryRow(Query)

	/*
		FIX:
		It will always return 0 when you use sql.Result.LastInsertId() function to get last inserted row id.
		the proxysql not support transaction.
		So,I Query a max(id) after insert a row.
	*/
	err = rows.Scan(&qr.RuleID)
	if err != nil {
		return errors.Trace(err)
	}

	LoadQueryRulesToRuntime(db)
	SaveQueryRulesToDisk(db)

	return nil
}

//DeleteOneQr delete a query rules.
func (qr *QueryRules) DeleteOneQr(db *sql.DB) error {

	Query := fmt.Sprintf(StmtDeleteOneQr, qr.RuleID)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(qr.RuleID)))
	}

	LoadQueryRulesToRuntime(db)
	SaveQueryRulesToDisk(db)

	return nil
}

//UpdateOneQrInfo update a query rules.
func (qr *QueryRules) UpdateOneQrInfo(db *sql.DB) error {

	var Query string

	Query = fmt.Sprintf(StmtUpdateOneQr,
		qr.Active,
		qr.Username,
		qr.Schemaname,
		qr.FlagIN,
		qr.ClientAddr,
		qr.ProxyAddr,
		qr.ProxyPort,
		qr.Digest,
		qr.MatchDigest,
		qr.MatchPattern,
		qr.NegateMatchPattern,
		qr.FlagOUT,
		qr.ReplacePattern,
		qr.DestinationHostgroup,
		qr.CacheTTL,
		qr.Reconnect,
		qr.Timeout,
		qr.Retries,
		qr.Delay,
		qr.MirrorFlagOUT,
		qr.MirrorHostgroup,
		qr.ErrorMsg,
		qr.Log,
		qr.Apply,
		qr.Comment,
		qr.RuleID)

	result, err := db.Exec(Query)
	if err != nil {
		return errors.Trace(err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.NotFoundf(strconv.Itoa(int(qr.RuleID)))
	}

	LoadQueryRulesToRuntime(db)
	SaveQueryRulesToDisk(db)

	return nil
}
