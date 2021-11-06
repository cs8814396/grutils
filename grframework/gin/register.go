package gin

import (
	"encoding/json"
	"fmt"
	"http"

	"github.com/gdgrc/grutils/grapps/config"
	"github.com/gin-gonic/gin"
)

func ResponseMap(c *gin.Context, resultMap map[string]interface{}, isBeauty bool) {
	var data []byte
	var err error
	if isBeauty {
		data, err = json.MarshalIndent(resultMap, "", "      ")
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			c.String(http.StatusOK, msg)
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	} else {
		data, err = json.Marshal(resultMap)
		if err != nil {
			msg := fmt.Sprintf(`{"result":0, "msg":"last decode err"}`)
			c.String(http.StatusOK, msg)
			config.DefaultLogger.Error(msg + err.Error())
			return
		}

	}

	c.String(http.StatusOK, string(data))

}
