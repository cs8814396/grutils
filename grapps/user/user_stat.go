package user

//drsql "database/sql"

//"grutils/grdatabase"
//"reflect"

var LastUid int64

/*
func FrontUserStat(appid string) (err error) {
	sql := "SELECT `uid`,`front_uid`,`extra_info`  FROM `user` where `appid`= ? AND `uid`>0"

	paramSlice := make([]interface{}, 0)
	paramSlice = append(paramSlice, appid)

	adminConn, err := config.DataAdminGet(false)
	if err != nil {
		return
	}

	rows, err := adminConn.Query(sql, paramSlice...)
	if err != nil {

		return
	}
	defer rows.Close()

	err, conn := grcache.RedisConnOPGet(config.GRAPPS_REDIS, nil)

	if err != nil {

		return

	}
	defer conn.Close()

	for rows.Next() {
		var newItem User
		err = rows.Scan(&newItem.Uid, &newItem.FrontUid, &newItem.ExtraInfoString)
		if err != nil {
			return
		}
		if err = json.Unmarshal([]byte(newItem.ExtraInfoString), &newItem.ExtraInfo); err != nil {
			return
		}
		var chnId int64
		chnId = 0
		if newItem.FrontUid == 0 {
			chnId = newItem.ExtraInfo.ChnId
		}
		redisSetName := fmt.Sprintf("FrontUserStat?appid=%s&uid=%d&chn_id=%d", appid, newItem.FrontUid, chnId)

		_, _, err = conn.Do("SADD", redisSetName, newItem.Uid)
		if err != nil {
			return
		}

		//fmt.Println(result)

		LastUid = newItem.Uid

	}

	return
}
func GetNextUser(appid string, uid int64, chnId int64, iType int, layer int) (nextUserMap map[string][]User, nextUserCountMap map[string]int, err error) {
	err, conn := grcache.RedisConnOPGet(config.GRAPPS_REDIS, nil)

	if err != nil {

		return

	}
	defer conn.Close()

	//redisSetName := fmt.Sprintf("FrontUserStat?appid=%s&uid=%d", appid, uid)

	nextUserMap = make(map[string][]User)
	nextUserMap["l0"] = []User{}
	nextUserMap["l0"] = append(nextUserMap["l0"], User{Appid: appid, Uid: uid, ExtraInfo: UserExtraInfo{ChnId: chnId}})

	nextUserCountMap = make(map[string]int)
	nextUserCountMap["l0"] = 1

	for i := 0; i < layer; i++ {
		var result interface{}
		redisScanCursor := "0"

		lastMapKey := fmt.Sprintf("l%d", i)
		mapKey := fmt.Sprintf("l%d", i+1)

		nextUserMap[mapKey] = make([]User, 0)
		nextUserCountMap[mapKey] = 0

		//fmt.Printf("%+v\n", nextUserMap)

		for _, frontUser := range nextUserMap[lastMapKey] {
			var chnId int64
			chnId = 0
			if frontUser.Uid == 0 { // notice that it should be uid not frontuid
				chnId = frontUser.ExtraInfo.ChnId
			}
			redisSetName := fmt.Sprintf("FrontUserStat?appid=%s&uid=%d&chn_id=%d", appid, frontUser.Uid, chnId)

			for {

				result, _, err = conn.Do("SSCAN", redisSetName, redisScanCursor)
				if err != nil || result == nil {
					err = fmt.Errorf("sscan result nil or err: %s", err)
					return
				}
				if result != nil {
					tmpResult := result.([]interface{})

					cursorBytes, ok := tmpResult[0].([]byte)
					if !ok {
						err = fmt.Errorf("cursor transform err,data: %+v", tmpResult[0])
						return
					}

					dataSlice, ok := tmpResult[1].([]interface{})
					if !ok {
						err = fmt.Errorf("cursor transform err,data: %+v", tmpResult[0])
						return
					}

					for _, data := range dataSlice {

						dataBytes, ok := data.([]byte)
						if !ok {
							err = fmt.Errorf("data to databytes err,data: %+v", data)
							return
						}

						uidString := string(dataBytes)
						var uid int64
						uid, err = strconv.ParseInt(uidString, 10, 64)
						if err != nil {
							err = fmt.Errorf("ParseInt err,data: %s", uidString)
							return
						}
						var newItem User
						var exist bool
						newItem, exist, err = GetUserFromDBorCache(appid, uid, true)
						if err != nil {
							return
						}

						if exist {

							var baseInfoItem User
							baseInfoItem.Appid = newItem.Appid
							baseInfoItem.Uid = newItem.Uid
							baseInfoItem.FrontUid = newItem.FrontUid
							baseInfoItem.ExtraInfo.Nickname = newItem.ExtraInfo.Nickname
							baseInfoItem.ExtraInfo.HeadImgUrl = newItem.ExtraInfo.HeadImgUrl

							nextUserMap[mapKey] = append(nextUserMap[mapKey], baseInfoItem)

							nextUserCountMap[mapKey] += 1
						} else {
							config.DefaultLogger.Error("GetNextUser not found redisSetName: %s, appid: %s,uid: %d", redisSetName, appid, uid)
						}

					}

					redisScanCursor = string(cursorBytes)

					if redisScanCursor == "0" {
						break
					}

				}
			}
		}

	}
	delete(nextUserMap, "l0")

	delete(nextUserCountMap, "l0")
	return

}
*/
