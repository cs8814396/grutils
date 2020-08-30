package service

import (
	"fmt"
	"github.com/gdgrc/grutils/grapps/config/log"
	"github.com/gdgrc/grutils/grserver/data_fetcher/data_fetcherconf"
	model "github.com/gdgrc/grutils/grserver/data_fetcher/model"
	"strings"
	//"data_fetcher/pb/data_fetcher"
)

type Statement struct {
	Params            []interface{}
	Request           *model.FetchDataReq
	PreparedStatement string
}

func (this *Statement) GetParams() []interface{} {
	return this.Params
}
func (this *Statement) GetRecordPreparedStatement() string {
	return fmt.Sprintf("%s limit %d,%d", this.PreparedStatement, (this.Request.PageNo-1)*this.Request.PageSize, this.Request.PageSize)
}
func (this *Statement) GetCountPreparedStatement() string {
	return fmt.Sprintf("select count(*) as count  from ( %s ) as a", this.PreparedStatement)
}
func ConstructMainStatment(req *model.FetchDataReq, dataConf *data_fetcherconf.Query) (s *Statement, err error) {
	s = &Statement{}
	s.Request = req

	//params = make([]interface{}, 0, len(dataConf.Conditions))
	s.PreparedStatement = dataConf.Statement
	placeHolderMap := make(map[string]*model.PlaceHolder)
	maxPlaceHolderLength := 0
	minPlaceHolderLength := int(^uint(0) >> 1)

	//---------construct placeholder map-----------
	for inputName, condition := range dataConf.Conditions {
		//-----------find out if condition is input
		rule, inputOk := req.Condition[inputName]

		if !inputOk {
			// condition is not input

			if condition.Default == "" {
				err = fmt.Errorf("inputName: %s should not be empty", inputName)
				return
			}

			placeHolder := fmt.Sprintf("$%s", inputName)
			placeHolderMap[placeHolder] = &model.PlaceHolder{ReplacedStatement: condition.Default, Params: []interface{}{}}

		} else {
			//condition is  input
			for ruleName, ruleValueList := range rule {
				innerStatement, ruleOk := model.RuleTable[ruleName]
				if ruleOk {
					partStatement := ""
					partStatement += fmt.Sprintf("%s %s ", condition.ColumnName, innerStatement)
					valueListLength := len(ruleValueList)
					tmpParams := make([]interface{}, 0, valueListLength)
					if valueListLength == 0 {
						// empty value list
						err = fmt.Errorf("ruleName: %s has empty value list", ruleName)
						return
					} else {
						// found rule, and the rule is legal

						if innerStatement == "in" {

							partStatement += "("

							for _, ruleValue := range ruleValueList {

								partStatement += "?,"
								tmpParams = append(tmpParams, ruleValue)

							}

							partStatement = partStatement[:len(partStatement)-1]
							partStatement += ") "

						} else {
							// if not in statement, just fetch the first value

							partStatement += "? "
							tmpParams = append(tmpParams, ruleValueList[0])
						}

					}
					placeHolder := fmt.Sprintf("$%s", inputName)
					placeHolderLength := len(placeHolder)
					if placeHolderLength > maxPlaceHolderLength {
						maxPlaceHolderLength = placeHolderLength
					}
					if placeHolderLength < minPlaceHolderLength {
						minPlaceHolderLength = placeHolderLength
					}
					placeHolderMap[placeHolder] = &model.PlaceHolder{ReplacedStatement: partStatement, Params: tmpParams}
					/*
						placeHolder := fmt.Sprintf("$%s", inputName)
						if !strings.Contains(sm, placeHolder) {
							err = fmt.Errorf("placeHolder: %s does not exist in statement?!", placeHolder)
							return
						}
						sm = strings.Replace(sm, placeHolder, partStatement, -1)
						params = append(params, tmpParamsList...)
						// only fetch the first legal rule*/
					break
				}

			}

		}
	}
	//--------------use placeholder map to replace placehold in the statement in order---------------

	// bruce force --- wil be optimized later
	sIndexI := 0
	sIndexJ := sIndexI + minPlaceHolderLength
	smLength := len(s.PreparedStatement)
	tmpSm := s.PreparedStatement
	log.Debug("smLength: ", smLength)

	for sIndexI < smLength && sIndexJ <= smLength {
		suspectString := s.PreparedStatement[sIndexI:sIndexJ]
		//log.Debug("suspectStringaaaaaaaaaaaaa:", suspectString)
		ph, ok := placeHolderMap[suspectString]
		if ok {
			//match and pass.  min match
			sIndexI = sIndexJ
			sIndexJ = sIndexI + minPlaceHolderLength
			tmpSm = strings.Replace(tmpSm, suspectString, ph.ReplacedStatement, 1)
			s.Params = append(s.Params, ph.Params...)
		} else {
			//match fail and pass.
			sIndexI = sIndexI + 1
			sIndexJ = sIndexI + minPlaceHolderLength
		}

	}
	s.PreparedStatement = tmpSm

	return
}
