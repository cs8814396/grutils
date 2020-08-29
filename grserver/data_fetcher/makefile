





all: compile

compile:
	echo "========================";
	cd data_fetcher && go build . && mv data_fetcher ../bin/;
	echo "========================";
	
	

clean:
	echo "gaaaood"

package: compile
	rm -rf data_fetcher_pkg;
	mkdir -p data_fetcher_pkg/etc/;
	mkdir -p data_fetcher_pkg/bin/;
	mv data_fetcher/data_fetcher data_fetcher_pkg/bin/;
	cp etc/config.debug.xml data_fetcher_pkg/etc/;
	cp etc/data_fetcher_conf.toml data_fetcher_pkg/etc/;
	cp README.md data_fetcher_pkg/;
	rm -f data_fetcher_pkg.tar.gz;
	tar -zcvf data_fetcher_pkg.tar.gz data_fetcher_pkg/;



test: test_prepare test_data_fetcherconf test_data_handler
	
	

test_data_fetcherconf:
	cd data_fetcherconf;go test -v -count=1 .;
	

test_data_handler:
	cd data_fetcher;go test -v -count=1 .;



test_db_create_sql :="create database if not exists test_database charset=utf8mb4"

test_tb_create_sql :="create table if not exists test_table ( aa varchar(128) not null default '', bb int not null,primary key(aa)) ENGINE=InnoDB  DEFAULT CHARSET=utf8"


test_tb_insert_sql :="insert ignore into test_table values ('aaa',111);"
test_tb_insert_sql2 :="insert ignore into test_table values ('bb',11);"
test_tb_insert_sql3 :="insert ignore into test_table values ('cc',201);"
test_tb_insert_sql4 :="insert ignore into test_table values ('ccc',1201);"
test_prepare:

	mysql -udev -puroad -e $(test_db_create_sql);
	mysql -udev -puroad test_database -e $(test_tb_create_sql);
	mysql -udev -puroad test_database -e $(test_tb_insert_sql);
	mysql -udev -puroad test_database -e $(test_tb_insert_sql2);
	mysql -udev -puroad test_database -e $(test_tb_insert_sql3);
	mysql -udev -puroad test_database -e $(test_tb_insert_sql4);
	

	# ./data_fetcher -c ../etc/config.debug.xml -sc ../data_fetcherconf/data_fetcher_conf.toml 

	# curl -d'{"data_name":"test_data","page_size":50,"condition":{"test_condition1":{"in":["aa","bb","ccc"]},"test_condition2":{"gte":["110"]}}}' "http://127.0.0.1:1112/data_fetcher/FetchData"
