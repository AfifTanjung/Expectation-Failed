<?php
// Pastikan koneksi basis data sudah ada

// Buat query untuk mengambil data dari database
$sql = "SELECT * FROM sensor";
$result = $conn->query($sql);

// Siapkan array untuk menampung data
$data = array();

// Periksa jika ada data yang ditemukan
if ($result->num_rows > 0) {
    // output data setiap baris
    while($row = $result->fetch_assoc()) {
        $data[] = $row;
    }
} else {
    echo "0 hasil";
}

// Mengubah data menjadi JSON dan output
header('Content-Type: application/json');
echo json_encode($data);

// Tutup koneksi database
$conn->close();
?>