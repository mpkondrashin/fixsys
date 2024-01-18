1. Copy fixsys.exe to some folder on the endpoint
2. Open CMD as administrator and run it providing Apex One agent unload password
   as parameter. !!! This will reboot the machine !!!
3. Login to the same endpoint and and run it once more providing Apex One
   unload password as parameter.

Notes:
- Log will be written to file fixsys.log

- On each run, fixsys writes local file that indicates to this program what is
  the next task. So running it third time will not do anything.
