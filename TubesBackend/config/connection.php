<?php

$db = mysqli_connect(hostname: 'localhost', username: 'root', password: '', database: 'sensor');

// cek koneksi

if (!$db){
    echo 'Gagal';
}

else {
    echo 'Berhasil';
}



?>