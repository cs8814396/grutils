module grutils

replace (
	cloud.google.com/go => github.com/google/go-cloud v0.4.1-0.20181112030950-0b43d4400e53
	go.uber.org/atomic => github.com/uber-go/atomic v1.3.3-0.20181018215023-8dc6146f7569
	go.uber.org/multierr => github.com/uber-go/multierr v1.1.1-0.20180122172545-ddea229ff1df
	go.uber.org/zap => github.com/uber-go/zap v1.9.2-0.20180814183419-67bc79d13d15
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20180904163835-0709b304e793
	golang.org/x/lint => github.com/golang/lint v0.0.0-20181026193005-c67002cb31c3
	golang.org/x/net => github.com/golang/net v0.0.0-20190108225652-1e06a53dbb7e
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20180821212333-d2e6202438be
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190109145017-48ac38b7c8cb
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20181108054448-85acf8d2951c
	golang.org/x/tools => github.com/golang/tools v0.0.0-20180828015842-6cd1fcedba52
	google.golang.org/api => github.com/google/google-api-go-client v0.0.0-20181108001712-cfbc873f6b93
	google.golang.org/appengine => github.com/golang/appengine v1.3.1-0.20181031002003-4a4468ece617
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190108161440-ae2f86662275
	google.golang.org/grpc => github.com/grpc/grpc-go v1.17.0

	honnef.co/go/tools => github.com/dominikh/go-tools v0.0.0-20180920025451-e3ad64cb4ed3
)

require (
	github.com/bradfitz/gomemcache v0.0.0-20190329173943-551aad21a668
	github.com/garyburd/redigo v1.6.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/miekg/dns v1.1.8
	github.com/panjf2000/ants v1.3.0
	github.com/ugorji/go v1.1.4
	golang.org/x/crypto v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/net v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/sys v0.0.0-00010101000000-000000000000 // indirect
	golang.org/x/text v0.0.0-00010101000000-000000000000
)

go 1.13
