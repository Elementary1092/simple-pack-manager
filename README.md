<h1>Simple package manager</h1>

# Description
This package manager (pm) can achieve files, upload archieve to the server using SSH.
Further, it can download archieves from the server and extract files from the archieve.

Sample package declaration file:

<b>packet.json</b>
```
{
 "name": "packet-1",
 "ver": "1.10",
 "targets": [
  "./archive_this1/*.txt",
  {"path", "./archive_this2/*", "exclude": "*.tmp"},
 ]
 packets: {
  {"name": "packet-3", "ver": "<="2.0" },
 }
}
```

Sample package description file:

<b>package.json</b>
```
{
 "packages": [
  {"name": "packet-1", "ver": ">=1.10"},
  {"name": "packet-2" },
  {"name": "packet-3", "ver": "<="1.10" },
 ]
}
```

# Usage
pm create ./packet.json - upload package to the server

pm update ./packages.json - dowload package from the server
