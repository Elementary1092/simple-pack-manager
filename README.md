<h1>Simple package manager</h1>

# Description
This package manager (pm) can achive files, upload archive to the server using SSH.
Further, it can download archives from the server and extract files from the archive.

Sample package declaration file:

<b>packet.json</b>
```
{
 "name": "packet-1",
 "ver": "1.10",
 "targets": [
  {"path": "./archive_this1/*.txt"},
  {"path": "./archive_this2/*", "exclude": "*.tmp"}
 ],
 packets: {
  {"name": "packet-3", "ver": "<=2.0" }
 }
}
```
Note: files will be collected only from listed directories.
This version does not collect files from subdirectories of a directory.
To include all files in a directory end pattern with "/*".

Sample package description file:

<b>package.json</b>
```
{
 "packages": [
  {"name": "packet-1", "ver": ">=1.10"},
  {"name": "packet-2" },
  {"name": "packet-3", "ver": "<=1.10" }
 ]
}
```
Note: package manager assumes that package description file contains all needed dependencies.
It means that only listed files will be downloaded.

# Usage
pm -create ./packet.json - upload package to the server

pm -update ./packages.json - dowload package from the server

# Build
Prerequisites:
- Go 1.16+
- Create ./internal/adapter/pmssh/ssh-ip.conf file 
    and write there IP-address of the server or use compile time value insertion
- Create ./internal/adapter/pmssh/ssh-port.conf file 
    and write there IP-address of the server or use compile time value insertion
- Create ./internal/adapter/pmssh/ssh-user.conf file 
    and write there IP-address of the server or use compile time value insertion
- Create ./internal/adapter/pmssh/ssh-pass.conf file 
    and write there IP-address of the server or use compile time value insertion.
    If you want to connect to the server using private key, create file ssh-key.pem
    file in ./internal/adapter/pmssh/ directory and paste private key there (or use compile time value insertion).

Note: compile time value insertion is done via -ldflags="-X <module_name>.<variable_name>=Your_value".
