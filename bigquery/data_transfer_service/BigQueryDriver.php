<?php
namespace app\driver;
use Google\Cloud\BigQuery\BigQueryClient;
class BigQueryDriver
{
    protected $_connection = null;

    public function __construct(array $aParams) {
        $this->_connection = new BigQueryClient([
            'projectId' => $aParams['project_id'],
        ]);
    }

    public function queryResultAsArray($sql) {
        $queryJobConfig = $this->_connection->query($sql);
        try {
            $queryResults = $this->_connection->runQuery($queryJobConfig);
        }catch (\Exception $e) {
            var_dump($e->getMessage());
        }
        if ($queryResults->isComplete()) {
            $rows = $queryResults->rows();
            $result = [];
            foreach ($rows as $row) {
                $result[] = $row;
            }
            return $result;
        } else {
            die(pg_last_error($this->_connection));
        }
    }
}
