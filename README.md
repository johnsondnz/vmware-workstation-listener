# vmware-workstation-listener
A go mqtt subscriber for processing power on and power off commands for virtual machines.  I use this with home assistant to trigger the remote starting of a VM.  I found this necessary as VMWare Workstation REST API is hard coded to bind to 127.0.0.1 with default port of 8697.

This program runs local to VMWare Workstation and listens to an MQTT topic for power on/shutdown commands.  I build this for my Microsoft Windows 11 Home PC.

# Home Assistant
```
switch:
    - platform: mqtt
        unique_id: eve_ng_switch
        name: EVE-NG WOL
        command_topic: home/servers
        payload_on: '{"id": "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "state": "on"}'
        payload_off: '{"id": "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "state": "shutdown"}'
        retain: false
```

# Usage
## .env file example
Create .env file for use along side .exe
```
AUTH_TOKEN=abcdefghijklmnopqrstuvwxyz
MQTT_SERVER=192.168.1.1
MQTT_PORT=1883
MQTT_USERNAME=mqtt
MQTT_PASSWORD=mqtt
MQTT_TOPIC=home/servers
```

- AUTH_TOKEN: This is the encoded base64 string used to auth with VMWare Workstation REST API
- Rest is self explanatory

## Steps
- build the executable with `env GOOS=target-OS GOARCH=target-architecture go build`
  - [Go Builds Article](https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04#step-4-building-executables-for-different-architectures)
- install `.env` files and `vmware-workstation-listener.exe` in the same directory on your machine
- Install exe as a scheduled task to start when computer starts, login as you user
  - Or install as a service
- Run

# Requires
VMWare Workstation 15+ and vmrest.exe to be started as a local service of scheduled task to start on boot.  [VMWare Docs](https://docs.vmware.com/en/VMware-Workstation-Pro/15.0/com.vmware.ws.using.doc/GUID-C3361DF5-A4C1-432E-850C-8F60D83E5E2B.html)
