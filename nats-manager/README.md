# nats-manager

nats-manager manages accounts for embedded Timeterm (frontend-embedded) devices, referred to as emdevs.

## Granted Permissions

Streams:
- Stream: `EMDEV-DISOWN-TOKEN`  
  Consumer: `EMDEV-{deviceId}`   
  Consumer topic filter: `EMDEV.{deviceId}.DISOWN-TOKEN`  
  Source topic: `EMDEV.*.DISOWN-TOKEN`  
  Makes for ACL entry: <kbd>pub</kbd> `$JS.API.CONSUMER.MSG.NEXT.EMDEV-DISOWN-TOKEN.EMDEV-{id}`

Topics:
- `EMDEV.{deviceId}.REBOOT`
  Makes for ACL entry: <kbd>sub</kbd> `EMDEV.{deviceId}.REBOOT`
