<?php
function set_cache_headers($file) {
    $mtime = filemtime($file);
    $etag = '"' . $mtime . '"';
    header('Cache-Control: max-age=3600, must-revalidate');
    header('ETag: ' . $etag);
    if (isset($_SERVER['HTTP_IF_NONE_MATCH']) && $_SERVER['HTTP_IF_NONE_MATCH'] == $etag) {
        header('HTTP/1.1 304 Not Modified');
        exit;
    }
}
?>
