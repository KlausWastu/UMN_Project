Service ini digunakan untuk menunjang aplikasi yang bertujuan melakukan pembersihan data mentah dari excel dan menyimpan data yang sudah bersih kedalam aplikasi.

# Petunjuk Penggunaan : 
1. Mengambil data Test-Backend.xlsx
2. Setelah melakukan pengambilan data, jalankan server
3. User diharapkan mendaftarkan diri pada API register
4. Setelah mendaftarkan diri, user diharapkan untuk login terlebih dahulu menggunakan email yang didaftarkan melalui API login
   (User diwajibkan login karena terdapat middleware yang membatasi akses kedalam service ini)
5. Data Test-Backend.xlsx yang diambil waktu pertama kali diimport kedalam field body dengan type file dan value berisi excel di API import-
6. Setelah itu terdapat pengecekan duplikasi pada file Test-Backend dan pengecekan ke database, jika data sudah clear maka terdapat response json dari cleanData
7. User dapat memasukan data yang sudah bersih itu kedalam database dengan cara melakukan input select single/multiple menggunakan no(number) pada kolom excel yang diimport, user dapat melakukan insert ini menggunakan API insert-

Berikut cara penggunaan service cleaning raw data, link dibawah ini merupakan contoh implementasi menggunakan postman (local)
https://documenter.getpostman.com/view/25757044/2sAXxV4pXv
