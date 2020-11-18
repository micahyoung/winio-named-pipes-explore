# demo of winio named-pipe IPC

## Running
Either:

On Windows:
```powershell
.\test.bat
```

or, with `DOCKER_HOST` set to WCOW daemon
```
docker build .
```

## Output

```
parent: 2020/11/18 13:13:24 opening pipe for writing \\.\pipe\cnb_exec_d
parent: 2020/11/18 13:13:24 running child subprocess \\.\pipe\cnb_exec_d
parent: 2020/11/18 13:13:24 listening on pipe \\.\pipe\cnb_exec_d
child: 2020/11/18 13:13:24 dialing pipe \\.\pipe\cnb_exec_d
child: 2020/11/18 13:13:24 writing message to pipe \\.\pipe\cnb_exec_d
parent: 2020/11/18 13:13:24 pipe content Hello World
parent: 2020/11/18 13:13:25 deferred closing of pipe \\.\pipe\cnb_exec_d
```
