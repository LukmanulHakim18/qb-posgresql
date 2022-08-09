# qb-Postgresql

I dedicate this query builder to be used in the `Blue Bird` service.

## Overview

- Featured ORM
- Easy to plug and play query condition
- Easy to do Transactions
- Every feature comes with tests
- Debug Mode

## Instalation

To install qb-postgresql, you need to install Go and set your Go workspace first.

1. The first need Go installed(version 1.14+ is required), then you can use the below Go command to install Gen.

```
go get -u github.com/LukmanulHakim18/qb-posgresql
```

2. Import it in your code:

```
import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/LukmanulHakim18/qb-posgresql"
)
```

## Rulse

Ketentuan dalam penggunaan lib ini adalah sebagai beriut:

- Entity harus memiliki tag db pada setiap field yang di binding ke column table `db:"foo"`
- Entity harus memiliki `Primery Key` dengan format penamaan `db:"id"` dan increment

## Quick start

Seluruh Kasus dibawah ini dilakukan dengan qb, dan jika ingin menggunakan library ini silahkan lihat rules Lib ini

### # Retrieving Single Row of Entity

Selain mengambil semua rekaman yang cocok dengan kueri tertentu, Anda juga dapat mengambil rekaman tunggal menggunakan metode `findFindOne`. methode ini menggambil data tunggal terbaru.

```go

// entity to receive data from db
comment := Comment{}

// make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

// Make query select from table comments and get last data
err := qb.Select("comments").FindOne(&comment)


// Make query select from table comments and get by id
err := qb.Select("comments").Where("id", "=", 82).FindOne(&comment)

// Make query select from table comments and get by email
err := qb.Select("comments").Where("email", "=", "anbukestra@bluebird.com").FindOne(&comment)

```

### # Retrieving Multy Rows of Entity

Mengambil semua rekaman yang cocok dengan kueri tertentu dengan menambahkan function - function filter seperti `AND`, `OR`, dll. metode `find` berfungsi untuk mengexecute query dan melakukan populate data hasil query ke slice of `Entitiy`

```go

// entity to receive datas from db
comments := []Comment{}

// make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

// Make query select from table comments
err := qb.Select("comments").
		Where("email", "LIKE", "%"+"aboy"+"%").
anbu
        // Memberikana filter in
		WhereIn("id", 83, 84, 85).

        // Memberikan filter OR
		OrWhere("id", "<>", 0).

        // Memberikan filter LIMIT
		Limit(20).

        // Memberikan filter OFFSET
		Offset(2).

        // Memberikan filter BETWEEN
		WhereBetween("created_at", time.Now().AddDate(0, -3, 0), time.Now()).

        // Mengurutkan data berdasarkan column dan DESC/ASC
		OrderBy("created_at", "asc").

        // Execute and populate
		Find(&comments)

```

### # Insert Row of Entity

Menyimpan data kedatabase menggunakan function `Insert` ini bertujuan untuk mempermudah insert karena tanpa perlu melakukan query yang masif. `Insert` akan memberikan `Id` terupdate dan `error nil` jika success namun memberikan `id = 0` dan `error` jika gagal.

```go
// Membuat entity yang berisikan value
comment := Comment{
    UUID:      goid.NewV4UUID().String(),
    Email:     "aboy@gmail.com",
    Comment:   "comment barusekai entity",
    CreatedAt: time.Now(),
}

// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

// Inserting data to database
Id, err := qb.Insert("comments", comment)


if err != nil {
    panic(err)
}

// Insert Id from DB to Entity
comment.Id = Id

fmt.Println(comment)
}
```

### # Update Row of Entity

Mengupdate data kedatabase menggunakan function `Update` ini bertujuan untuk mempermudah update karena tanpa perlu melakukan query yang masif. `Update` akan memberikan `error nil` jika success namun memberikan nilai `error` jika gagal.
Sebelum update data sebaiknya melakukan uery select untuk mengambil data terlebih dahulu, agar tidak ada data yang kosong atau berubah kecuali data yang ingin diubah.

```go

// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

comment := Comment{}

err := qb.Select("comments").Where("Id", "=", 86).FindOne(&comment)
if err != nil {
    panic(err)
}

// Melakukan Update data apa saja field yang ingin diubah
comment.Comment = "Updated comment"
comment.Email = "newemail@bluebird.co"

// Update ke database
err = qb.Update("comments", comment)
if err != nil {
    panic(err)
}
```

### # Delete Row of Entity

Untuk mrnghapus Data hanya membutuhkan `Id` kalian bisa membuat object dari entity dan menginputkan hanya Id saja.

```go

// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

// set id pada entity
comment := Comment{Id:86}

// Menghapus data dengan id yang sudah di set
err = qb.Delete("comments", comment)
if err != nil {
    panic(err)
}
```

### # Query RAW

Selain itu semua kami juga membuatkan function query raw untuk melakukan query agar lebih dinamis dan lebih fleksible jik menemukan case case tertentu

```go

// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

// lakukan query apapun disin
query := "SELECT COUNT(id) as total FROM comments"

// Inisialisasi tempat menampung data.
var totalDataComments int

// Eksekusi Query
res, err := qb.Raw(query)
if err != nil {
    panic(err)
}

// Cek Jika data ada dan Parsing ke tempat penampungan data.
if res.Next() {
    res.Scan(&totalDataComments)
}
```

Selain itu juga anda dapat melakukan query RAW menggunakan parameter

```go

// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

// sate date befor 3 mounth
startDate := time.Now().AddDate(0, -3, 0)

// lakukan query apapun disin
query := "SELECT COUNT(id) as total FROM comments WHERE created_at > ?"

// Inisialisasi tempat menampung data.
var totalDataComments int

// Eksekusi Query
res, err := qb.Raw(query, startDate)
if err != nil {
    panic(err)
}

// Cek Jika data ada dan Parsing ke tempat penampungan data.
if res.Next() {
    res.Scan(&totalDataComments)
}
```

## Make Transaction

### # Begin transaction

Untuk memulai transaction kita harus memanggil function `TrxBegin`
setelah itu kita dapat melakuakan query apa saja kedalam transaction ini dan nantinya akan dicommit ataupun di dirollback

```go
// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection())

// Memulai transaksi
qb.TrxBegin()

```

pada tahap ini data akan di simpan di temporary, sampai ada perintah berikutnya atau menunggu connection timeout

### # Commit transaction

Untuk melanjutkan transaction yang telah berjalan lancar atau sudah sesuai kita harus memanggil function `TrxCommit`
untuk menyimpan perubahan apapun perubahan yang telah dilakukan ke db dan menutup transaction

```go
// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection()

// Memulai transaksi
qb.TrxBegin()

// commit transaksi
qb.TrxCommit()

```

pada tahap ini data akan di simpan di DB.

### # Rollback transaction

Untuk membatalkan transaction yang telah berjalan atau merollback seluruh perubahan yang dilakukan didalam transaksi, harus memanggil function `TrxRollback`
untuk membatalkan semua perubahan yang dilakukan dan menutup transaction

```go
// Make instence of interface for use function *note: GetConnection returning *sql.DB
qb := qb.NewQueryBuilder(GetConnection()

// Memulai transaksi
qb.TrxBegin()

// commit transaksi
qb.TrxRollback()


```

## Debug Mode

Untuk mencetak query dan argument pada terminal dapat mengeset enverionment "DEBUG_MODE = true" dapat dilihat dari hasil query berikut

```go
qb := qb.NewQueryBuilder(GetConnection())
comments := []Comment{}
err := qb.Select("comments").
    Where("email", "LIKE", "%"+"aboy"+"%").
    WhereIn("id", 83, 84, 85).
    OrWhere("id", "<>", 0).
    Limit(20).
    Offset(2).
    WhereBetween("created_at", time.Now().AddDate(0, -3, 0), time.Now()).
    OrderBy("created_at", "asc").
    GroupBy("email", "id").
    Find(&comments)

if err != nil {
    panic(err)
}
```

Menghasilkan output sebagai berikut

```
query => SELECT * FROM comments WHERE email LIKE $1 AND id IN ($2, $3, $4) OR id <> $5 AND (created_at BETWEEN ($6) AND ($7)) GROUP BY email, id ORDER BY created_at ASC LIMIT 20 OFFSET 2
param => $1=>%aboy%,
         $2=>83,
         $3=>84,
         $4=>85,
         $5=>0,
         $6=>2022-05-09 17:06:05.7043322 +0700 +07,
         $7=>2022-08-09 17:06:05.7048644 +0700 +07 m=+0.003841101
```
