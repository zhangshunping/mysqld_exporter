#export DATA_SOURCE_NAME='exporter:123123@(zsp.vmware.com:3306)/'
./mysql_exporter --collect.info_schema.tables \
 --collect.info_schema.innodb_tablespaces \
 --collect.info_schema.innodb_metrics  \
 --collect.global_status \
 --collect.global_variables \
 --collect.slave_status  \
 --collect.info_schema.processlist \
 --collect.perf_schema.tablelocks \
 --collect.perf_schema.eventsstatements \
 --collect.perf_schema.eventsstatementssum \
 --collect.perf_schema.eventswaits \
 --collect.auto_increment.columns \
 --collect.binlog_size \
 --collect.perf_schema.tableiowaits \
 --collect.perf_schema.indexiowaits \
 --collect.info_schema.userstats \
 --collect.info_schema.clientstats \
 --collect.info_schema.tablestats \
 --collect.info_schema.schemastats \
 --collect.perf_schema.file_events \
 --collect.perf_schema.file_instances \
 --collect.perf_schema.replication_group_member_stats \
 --collect.slave_hosts \
 --collect.info_schema.innodb_cmp \
 --collect.info_schema.innodb_cmpmem \
 --collect.info_schema.query_response_time \
 --collect.engine_innodb_status   \
 --collect.node_exporter_meminfo 

