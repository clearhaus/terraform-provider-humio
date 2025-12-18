resource "humio_parser" "emojis_loglevels" {
  repository    = "sandbox"
  name          = "emojis"
  parser_script = <<PARSERSCRIPT
case {
  🔥 | loglevel:="error";
  😕 | loglevel:="warning";
  🙂 | loglevel:="info";

//Everything else
  * | loglevel:="unknown";
}
PARSERSCRIPT
}

resource "humio_parser" "filebeat" {
  repository    = "sandbox"
  name          = "filebeat-log"
  parser_script = <<PARSERSCRIPT
regex("(?<@timestamp>\\S+)\\s+(?<loglevel>\\S+)\\s+\\[?(?<source>[^\\ \\]]+)") | parseTimestamp(field=@timestamp, format="yyyy-MM-dd'T'HH:mm:ss[.SSS]XXX")
PARSERSCRIPT
}

resource "humio_parser" "gc" {
  repository    = "sandbox"
  name          = "gc"
  test_data     = ["[2018-12-05T08:21:06.750+0000][gc             ] GC(5016) Pause Young (Allocation Failure) 2435M->1587M(2553M) 5.766ms"]
  parser_script = <<PARSERSCRIPT
/\[(?<@timestamp>[^\]]+)\]/ |
parseTimestamp(format="yyyy-MM-dd'T'HH:mm:ss.SSSX", field=@timestamp) |
kvParse() |
regex("(?<@gctime>[0-9.]+)ms", strict=false)
PARSERSCRIPT
}

resource "humio_parser" "gc_kafka" {
  repository = "sandbox"
  name       = "gc-kafka"
  test_data = [
    "[2018-12-05T07:37:48.685+0000][gc            ] GC(111957) Pause Young (G1 Evacuation Pause) 711M->559M(796M) 9.860ms",
    "[2018-12-05T07:38:20.790+0000][gc,start      ] GC(133000) Pause Young (G1 Evacuation Pause)",
  ]
  parser_script = <<PARSERSCRIPT
/\[(?<@timestamp>[^\]]+)\]/ |
parseTimestamp(format="yyyy-MM-dd'T'HH:mm:ss.SSSX", field=@timestamp) |
kvParse() |
regex("(?<@gctime>[0-9.]+)ms", strict=false)
PARSERSCRIPT
}

resource "humio_parser" "gc_zk" {
  repository    = "sandbox"
  name          = "gc-zk"
  test_data     = ["[2018-12-05T08:11:25.249+0000][gc,cpu       ] GC(6) User=0.03s Sys=0.00s Real=0.00s"]
  parser_script = <<PARSERSCRIPT
/\[(?<@timestamp>[^\]]+)\]/ |
parseTimestamp(format="yyyy-MM-dd'T'HH:mm:ss.SSSX", field=@timestamp) |
kvParse() |
regex("(?<@gctime>[0-9.]+)ms", strict=false)
PARSERSCRIPT
}

resource "humio_parser" "http_error" {
  repository    = "sandbox"
  name          = "http-error"
  parser_script = <<PARSERSCRIPT
regex("(?<@timestamp>[\\d\\/\\:]+\\s+[\\d\\/\\:]+)\\s+(\\[(?<severity>[^\\]]+)\\])?") | parseTimestamp(field=@timestamp, format="yyyy/MM/dd' 'HH:mm:ss", timezone="UTC")
PARSERSCRIPT
}

resource "humio_parser" "humio_stdout" {
  repository    = "sandbox"
  name          = "humio-stdout"
  parser_script = <<PARSERSCRIPT
regex("") | parseTimestamp(field=@timestamp, format="yyyy-MM-dd'T'HH:mm:ss.SSSZ")
PARSERSCRIPT
}

resource "humio_parser" "jenkins_build_log" {
  repository = "sandbox"
  name       = "jenkins-build-log"
  tag_fields = ["build", "task"]
  test_data = [
    "2018-10-15T12:51:40+00:00 [INFO] This is an example log entry. id=123 fruit=banana",
    "2018-10-15T12:52:42+01:30 [ERROR] Here is an error log entry. class=c.o.StringUtil fruit=pineapple",
    "2018-10-15T12:53:12+01:00 [INFO] User logged in. user_id=1831923 protocol=http",
  ]
  parser_script = <<PARSERSCRIPT
@source=/jenkins-data\/jobs\/(?<task>[^\/]+)\/builds\/(?<build>\d+)/
PARSERSCRIPT
}

resource "humio_parser" "kafka" {
  repository    = "sandbox"
  name          = "kafka"
  parser_script = <<PARSERSCRIPT
regex("\\[(?<@timestamp>[^\\]]+)\\]\\s+(?<loglevel>\\S+)\\s+(\\[(?<thread>[^]]+)\\]\\:)?") | parseTimestamp(field=@timestamp, format="yyyy-MM-dd' 'HH:mm:ss,SSS", timezone="UTC") | kvParse()
PARSERSCRIPT
}
resource "humio_parser" "kafka_gc" {
  repository = "sandbox"
  name       = "kafka-gc"
  test_data = [
    "[2018-12-05T07:37:48.685+0000][gc            ] GC(111957) Pause Young (G1 Evacuation Pause) 711M->559M(796M) 9.860ms",
    "[2018-12-05T07:38:20.790+0000][gc,start      ] GC(133000) Pause Young (G1 Evacuation Pause)",
  ]
  parser_script = <<PARSERSCRIPT
/\[(?<@timestamp>[^\]]+)\]/ |
parseTimestamp(format="yyyy-MM-dd'T'HH:mm:ss.SSSX", field=@timestamp) |
kvParse() |
regex("(?<@gctime>[0-9.]+)ms", strict=false)
PARSERSCRIPT
}
resource "humio_parser" "metricbeat" {
  repository    = "sandbox"
  name          = "metricbeat-log"
  parser_script = <<PARSERSCRIPT
regex("(?<@timestamp>\\S+)\\s+(?<loglevel>\\S+)\\s+\\[?(?<source>[^\\ \\]]+)") | parseTimestamp(field=@timestamp, format="yyyy-MM-dd'T'HH:mm:ss[.SSS]XXX")
PARSERSCRIPT
}

resource "humio_parser" "unattended_upgrades" {
  repository    = "sandbox"
  name          = "unattended-upgrades"
  parser_script = <<PARSERSCRIPT
/^(?<ts>\S+?\s\S+?)\s(?<loglevel>\w+?)/

| @timestamp := parseTimestamp("yyyy-MM-dd HH:mm:ss[,SSS]", field=ts, timezone=UTC)
| drop([ts])
| kvParse()
PARSERSCRIPT
}

resource "humio_parser" "webfront" {
  repository = "sandbox"
  name       = "webfront"
  tag_fields = ["@host", "type"]
  test_data = [
    "85.2.10.134 - - [20/Feb/2019:09:03:55 +0100] \"GET /?page=1&search=&tab=dataspaces HTTP/1.1\" 200 389 \"-\" \"'\"><svg/onload=(new(Image)).src='//0yzm97bocjcdyi9k3v7jdromadgc4a6yyomj99xy\\56burpcollaborator.net'>\" 3 - \"-\" \"Caddy,Humio-1.4.2--build-4731--sha-4f2a2ce2d\" \"-\"",
    "95.142.1.129 - - [23/Jan/2019:12:44:04 +0100] \"POST /api/v1/dataspaces/flaf-dev/ingest HTTP/1.1\" 200 2 \"-\" \"HttpWebRequestSender\" 34 2963 \"http://10.0.2.5:8080\" \"Caddy,Humio-1.2.9--build-4799--sha-73345945c\"",
    "13.80.5.2 - - [23/Jan/2019:12:40:30 +0100] \"POST /api/v1/ingest/elastic-bulk/_bulk HTTP/1.1\" 200 27 \"-\" \"-\" 59 19657 \"http://10.0.1.4:8080\" \"Caddy,Humio-1.2.10--build-4840--sha-3d524e860\" \"foo\"",
    "Dec 3 21:35:06 webfront02 kernel: [ 0.126594] e820: reserve RAM buffer [mem 0xc3f69000-0xc3ffffff]",
  ]
  parser_script = <<PARSERSCRIPT
case {
//System logs
  @source="/var/log/syslog" or @source="/var/log/auth.log" | /^(?<ts>\w+\s\s?\d+\s+\S+?) (?<host>\S+?) (?<app>.+?):/
    | parseTimestamp("MMM [ ]d HH:mm:ss", field=ts, timezone="Europe/Berlin");
  @source="/var/log/caddy/*.requests.log"
    | type:="requests"
    | kvparse()
    | /^(?<ts>\S+\s\S+)\]/
    | @timestamp := parseTimestamp("dd/MMM/yyyy:HH:mm:ss Z", field=ts, timezone="Europe/Berlin");
//Caddy access logs
  @source="/var/log/caddy/*.log"
    | /^(?<client>\S+) - (?<userid>\S+?) \[(?<ts>\S+\s\S+)\] "(?<method>\S+?) (?<url>\S+) (?<httpversion>\S+)" (?<statuscode>\d+) (?<body_bytes_sent_to_client>\S+) "(?<referrer>.*?)" "(?<useragent>.+?)" (?<responsetime_ms>\d+?) (?<request_length>.+?) "(?<upstream>\S+)" "(?<upstream_server>\S+?)" "(?<humio_query_session>.+?)" (?<tlsversion>\S+?) (?<tlscipher>\S+)/
    | type:=accesslog
    | responsetime:=responsetime_ms/1000
    | parseUrl(upstream)
    | regex("^/var/log/caddy/(?<host>.+)\.(?<protocol>.+?)\.log", field=@source)
    | @timestamp := parseTimestamp("dd/MMM/yyyy:HH:mm:ss Z", field=ts, timezone="Europe/Berlin")
    | humiotime:="-" | machinetype:="unknown" /* backwards compatibility*/;
//Nginx access logs
  @source="/var/log/nginx/*.log"
    | /^(?<client>\S+) - (?<userid>\S+?) \[(?<ts>\S+\s\S+)\] "(?<method>\S+?) (?<url>\S+) (?<httpversion>\S+)" (?<statuscode>\d+) (?<body_bytes_sent_to_client>\S+) "(?<referrer>.*?)" "(?<useragent>.+?)" (?<responsetime>\d+?(\.\d+)?) (?<request_length>\d+?) "(?<upstream>\S+?(, \S+)?)" "(?<upstream_server>\S+?(, \S+)?)" "(?<humio_query_session>.+?)" (?<tlsversion>\S+?) (?<tlscipher>\S+) "(?<upstream_connect_time>\S+?(, \S+)?)" "(?<upstream_header_time>\S+?(, \S+)?)" "(?<upstream_response_time>\S+?(, \S+)?)" "(?<upstream_status>\S+?(, \S+)?)"/
    | type:=accesslog
    | responsetime_ms:=responsetime*1000
    | regex("^/var/log/nginx/(?<host>.+)\.(?<protocol>.+?)\.log", field=@source)
    | @timestamp := parseTimestamp("dd/MMM/yyyy:HH:mm:ss Z", field=ts, timezone="Europe/Berlin")
    | humiotime:="-" | machinetype:="unknown" /* backwards compatibility*/;
//Everything else
  * | parseerror:="Fallback to kvparse" | kvparse();
}

| drop([ts])
PARSERSCRIPT
}

resource "humio_parser" "zookeeper" {
  repository    = "sandbox"
  name          = "zookeeper"
  parser_script = <<PARSERSCRIPT
regex("\\[(?<@timestamp>[^\\]]+)\\]\\s+(?<loglevel>\\S+)\\s+(\\[(?<thread>[^]]+)\\]\\:)?") | parseTimestamp(field=@timestamp, format="yyyy-MM-dd' 'HH:mm:ss,SSS", timezone="UTC") | kvParse()
PARSERSCRIPT
}
