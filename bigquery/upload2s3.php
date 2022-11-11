<?php

// logserver conf:
// exec_command = /usr/bin/php /usr/local/logServer/exec/upload2s3.php
// s3_key_config = /usr/local/logServer/exec/s3conf_ffs.cfg
// s3_bucket = ffs-global-log
// s3_object_key = fluentd/%{Y}/%{m}/%{d}/%{H}/%{table}_endTN_global_%{ip}_%{index}_%{ts}.%{file_extension}
// gzip_tmp_dir = /mnt/tmp/gzip_tmp_dir
// upload_file_info = /mnt/logserver/ffs.global/s3.upload.info

set_time_limit(0);
date_default_timezone_set('UTC');

if (count($argv) < 2) {
    echo "upload2s3.php: [Error] parameter num must >= 2\n";
    exit(1);
}

if (is_executable('/usr/bin/s3cmd')) {
    $S3_CMD = '/usr/bin/s3cmd';
} else {
    $S3_CMD = '/usr/local/s3cmd/s3cmd';
}
$GZIP_CMD = '/bin/gzip';
$TMP_DIR = getenv('gzip_tmp_dir');
$S3_KEY_CONFIG = getenv('s3_key_config');
$S3_BUCKET = getenv('s3_bucket');
$S3_OBJECT_KEY = getenv('s3_object_key');
$FILE_NAME = $argv[1];

$Y = getenv('Y');
$m = getenv('m');
$d = getenv('d');
$H = getenv('H');
$M = getenv('M');
$S = getenv('S');
$hour_index = getenv('hour_index');
$ip = getenv('local_ip');
$file_info = getenv('upload_file_info');
$prodPos = strpos($S3_OBJECT_KEY, "_endTN_") + 7;
$prodLength = strpos($S3_OBJECT_KEY, "_%{ip}_") - $prodPos;
$prod = substr($S3_OBJECT_KEY, $prodPos, $prodLength);
$suffix = "{$Y}_{$m}_{$d}_{$H}_{$hour_index}_{$prod}";

register_shutdown_function(function () use ($suffix, $TMP_DIR) {
    exec("cd {$TMP_DIR}");
    exec("find . -name '*".$suffix."*' | xargs rm");
});

$outFiles = splitFile_fluentd($FILE_NAME, $TMP_DIR, $suffix);
exec("cd {$TMP_DIR}");

foreach ($outFiles as $curFileName => $tableName) {
    $search = array('%{Y}', '%{m}', '%{d}', '%{H}', '%{table}', '%{ip}', '%{index}', '%{file_extension}', '%{ts}');
    $replace = array($Y, $m, $d, $H, $tableName, $ip, $hour_index, 'gz', time());
    $name = str_replace($search, $replace, $S3_OBJECT_KEY);
    $uri = "s3://{$S3_BUCKET}/{$name}";

    $tmpName = $TMP_DIR.'/'.$curFileName.'.gz';
    $curFileName = $TMP_DIR.'/'.$curFileName;
    exec("$GZIP_CMD -c $curFileName > $tmpName");
    exec("$S3_CMD --config=$S3_KEY_CONFIG put $tmpName $uri");
    exec("$S3_CMD --config=$S3_KEY_CONFIG ls $uri", $out);
    if (isset($out[0]) && $out[0] != '') {
        // success
    } else {
        echo "upload2s3.php: [Error] upload to S3 fail ($uri)\n";
        exit(1);
    }
}

function splitFile_fluentd($filePath, $output_dir, $suffix)
{
    $outFiles = array();
    if (!file_exists($filePath)) {
        exit(1);
    }
    $handle = fopen($filePath, 'r');
    while (!feof($handle)) {
        $line = fgets($handle);
        if (empty($line)) {
            continue;
        }
        $lineData = json_decode($line, true);
        $fileName = $lineData['$tablename']."_{$suffix}";
        if (!isset($outFiles[$fileName])) {
            $outFiles[$fileName] = $lineData['$tablename'];
        }
        file_put_contents($output_dir.'/'.$fileName, $lineData['msg'], FILE_APPEND);
    }
    return $outFiles;
}
