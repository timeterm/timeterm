# nats-manager

nats-manager manages accounts for embedded Timeterm (frontend-embedded) devices, referred to as emdevs to keep topics short.

## Granted Permissions

Streams:
- Stream: `EMDEV-DISOWN-TOKEN`  
  Consumer: `EMDEV-{deviceId}`   
  Consumer topic filter: `EMDEV.{deviceId}.DISOWN-TOKEN`  
  Source topic: `EMDEV.*.DISOWN-TOKEN`  
  Makes for ACL entries: 
  - <kbd>pub</kbd> `$JS.API.CONSUMER.MSG.NEXT.EMDEV-DISOWN-TOKEN.EMDEV-{deviceId}`  
    > Required for requesting new messages.
  - <kbd>pub</kbd> `$JS.ACK.EMDEV-DISOWN-TOKEN.EMDEV-{deviceId}.>`  
    > Required for ACKing messages.  
      Not manually created by the client but set by the NATS server as 
      reply subject in responses to requests to the topic above.

Topics:
- `EMDEV.{deviceId}.REBOOT`  
  Makes for ACL entry: <kbd>sub</kbd> `EMDEV.{deviceId}.REBOOT`
