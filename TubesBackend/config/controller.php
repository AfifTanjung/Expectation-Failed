<?php


// Function PHP


// Fungsi Menampilkan PHP

function select($query)

{

// Memanggil Database

    global $db;

    $result = mysqli_query($db, $query);
    $rows = [];

    while ($row = mysqli_fetch_array($result)) {
        $rows[] = $row;
    }

    return $rows;

}

// Fungsi Menambah Data Sensor

function create_data($post)
{
    global $db;

    $# = $post['#'];

// Tambah Data

    $query = "INSERT INTO sensor VALUES(null, '$#')";
    
    mysqli_query($db, $query);

    return mysqli_affected_rows($db);

}


// Fungsi Mengubah Data


function update_data($post) {

    global $db;

    $# = $post['#'];

// Tambah Data

    $query = "UPDATE sensor SET # = '#' WHERE id_no = '$id_no'";   
    
    mysqli_query($db, $query);

    return mysqli_affected_rows($db);
}




