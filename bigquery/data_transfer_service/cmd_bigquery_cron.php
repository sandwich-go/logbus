<?php
// import
require __DIR__ . '/vendor/autoload.php';
require __DIR__ . '/BigQueryDriver.php';
use Google\Cloud\BigQuery\DataTransfer\V1\DataTransferServiceClient;
use Google\Cloud\BigQuery\DataTransfer\V1\TransferConfig;
use Google\Cloud\BigQuery\DataTransfer\V1\EmailPreferences;
use Google\Cloud\BigQuery\DataTransfer\V1\ScheduleOptions;
use Google\Protobuf\Struct;
use Google\Protobuf\Value;
use Google\Protobuf\Timestamp;

// config
define("BUCKET_NAME", "ffs-global-log");
define("AWS_ACCESS_KEY_ID", "xxxxx");
define("AWS_SECRET_ACCESS_KEY", "yyyyyy");
define("BIGQUERY_TRANSFER_PROJECT_ID", 363858487727);
$aTable = [];

// create table if not exist
$oCon = new app\driver\BigQueryDriver(['projectId' => 'ffseaside']);
$aData = $oCon->queryResultAsArray("SELECT * FROM public.INFORMATION_SCHEMA.TABLES;");
$aExistTable = [];
foreach ($aData as $aItem) {
    $aExistTable[$aItem['table_name']] = 1;
}
$sql = <<<EOT
create table public.%s (
    user_id STRING,
    optime TIMESTAMP,
    data STRING,
    app INT64
) PARTITION BY DATE(optime) CLUSTER BY user_id
EOT;
foreach ($aTable as $k => $v) {
    if (!$v) {
        continue;
    }
    if (!empty($aExistTable[$k])) {
        continue;
    }
    $oCon->queryResultAsArray(sprintf($sql, $k));
    executeLog("[Table] create table $sql");
}

// get transfer config
$client = new DataTransferServiceClient();
$parent = sprintf('projects/%s/locations/us', BIGQUERY_TRANSFER_PROJECT_ID);
$aExistTransfer = [];
try {
    $pagedResponse = $client->listTransferConfigs($parent);
    foreach ($pagedResponse->iterateAllElements() as $dataSource) {
        $aExistTransfer[$dataSource->getDisplayName()] = 1;
    }
} catch (Exception $e) {
    executeLog($e->getTraceAsString(), "[ERROR]");
    exit(1);
}

// create transfer config if not exist
foreach ($aTable as $k => $v) {
    if (!empty($aExistTransfer[$k])) {
        continue;
    }
    if (empty($aExistTable[$k])) {
        continue;
    }
    $transferConfig = new TransferConfig();
    $transferConfig->setDisplayName($k);
    $transferConfig->setDataSourceId("amazon_s3");
    $transferConfig->setDestinationDatasetId("public");
    $oReferences = new EmailPreferences();
    $oReferences->setEnableFailureEmail(true);
    $transferConfig->setEmailPreferences($oReferences);
    $transferConfig->setNotificationPubsubTopic("projects/ffseaside/topics/ffs_oplog_transfer");
    $oOptions = new ScheduleOptions();
    $oOptions->setDisableAutoScheduling(true);
    $transferConfig->setScheduleOptions($oOptions);
    $aParam = new Struct();
    $aParam->setFields(
        [
            "data_path" => StrValue('s3://'.BUCKET_NAME.'/fluentd/{run_time-1h|"%Y"}/{run_time-1h|"%m"}/{run_time-1h|"%d"}/{run_time-1h|"%H"}/' . $k . '_end*'),
            "destination_table_name_template" => StrValue($k),
            "access_key_id" => StrValue(AWS_ACCESS_KEY_ID),
            "secret_access_key" => StrValue(AWS_SECRET_ACCESS_KEY),
            "file_format" => StrValue("JSON"),
            "max_bad_records" => StrValue("200"),
            "ignore_unknown_values" => (new Value)->setBoolValue(true),
            "skip_leading_rows" => StrValue("0"),
            "allow_quoted_newlines" => StrValue("false"),
        ]
    );
    $transferConfig->setParams($aParam);
    $aRet = $client->createTransferConfig($parent, $transferConfig);
    executeLog("[Transfer] create transfer " . $k);
}

// run transfer job
try {
    $pagedResponse = $client->listTransferConfigs($parent);
    foreach ($pagedResponse->iterateAllElements() as $dataSource) {
        $client->startManualTransferRuns(['parent' => $dataSource->getName(), 'requestedRunTime' => (new Timestamp())->setSeconds(time())]);
        executeLog("[Schedule] run transfer " . $dataSource->getDisplayName());
        sleep(10);
    }
} catch (Exception $e) {
    executeLog($e->getTraceAsString(), "[ERROR]");
    exit(1);
} finally {
    $client->close();
}




function StrValue($str) {
    return (new Value())->setStringValue($str);
}

function executeLog($str, $prefix='[INFO]') {
    echo "[" . date('Y-m-d H:i:s') . "] " . $prefix . ' ' . $str . "\n";
}