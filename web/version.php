<?php
require_once __DIR__ . '/cache.php';
set_cache_headers(__FILE__);

$apiUrl = getenv('API_URL') ?: 'http://localhost:8080/api/calculate';
$apiBase = str_replace('/api/calculate', '', $apiUrl);

$ch = curl_init($apiBase . '/api/version');
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

if ($httpCode !== 200) {
    echo json_encode(['version' => 'n/a']);
    exit;
}

echo $response;
