# Mirror
## Installation:
copy mirror executable file and servers.yml to the directory you wan't your backup to be in
## Config:

edit servers.yml file and add your servers like this:

```yaml
serverName:
  host: 127.0.0.1
  port: 3306
  user: thabet
  pass: *****
  dbs:
    - database1
    - database2
```
you can add multiple servers like this
```yaml
serverName1:
  host: 127.0.0.1
  port: 3306
  user: thabet
  pass: *****
  dbs:
    - database1
    - database2

serverName2:
  host: 127.0.0.2
  port: 3306
  user: thabet
  pass: *****
  dbs:
    - database1
```
## Run the app:
after configuring your servers.yml file, double click on the executable file you'll get app icon on the tray of your OS, if not open it in **:2345**
