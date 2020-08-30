#data_fetcher


启动方式:

./data_fetcher -c ../etc/config.debug.xml -sc ../etc/data_fetcher_conf.toml  指定两个配置文件. -sc 的是数据库配置文件toml格式

curl -d'{"data_name":"test_data","page_size":50,"condition":{"test_condition1":{"in":["aa","bb","ccc"]},"test_condition2":{"gte":["110"]}}}' "http://127.0.0.1:9096/fetch_data"

{"errcode":0,"errmsg":"","data":{"total_page_num":1,"record_num":1,"total_record_num":1,"page_size":50,"page_no":1,"record_list":[{"aa":"ccc","bb":"1201"}]}}



curl -d'{"data_name":"query_web_device_env","page_size":50,"condition":{"device_fp1":{"eq":["5f8a50afa1ca41e98017ca4ddbc848cf2eca8426ce82b981cfd033d89fc520d4"]} }}' "http://127.0.0.1:9096/fetch_data"

配置文件编写:

现在暂时 config.debug.xml 没用．但是还是要指定. 关键定义好: data_fetcher_conf.toml,具体参考包里面的样例即可