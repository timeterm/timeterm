# nats-manager

nats-manager manages accounts for embedded Timeterm (frontend-embedded) devices, referred to as emdevs to keep topics short.

## Granted Permissions

Streams:
- Stream: `EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG`  
  Consumer: `EMDEV-{deviceId}-EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG`   
  Consumer topic filter: `EMDEV.{deviceId}.RETRIEVE-NEW-NETWORKING-CONFIG`
  Source topic: `EMDEV.*.RETRIEVE-NEW-NETWORKING-CONFIG`
  
  Makes for ACL entries: 
  - <kbd>pub</kbd> `$JS.API.CONSUMER.MSG.NEXT.EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG.EMDEV-{deviceId}-EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG`
    > Required for requesting new messages.

  - <kbd>pub</kbd> `$JS.ACK.EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG.EMDEV-{deviceId}-EMDEV-RETRIEVE-NEW-NETWORKING-CONFIG.>`  
      
    > Required for ACKing messages.  
      Not manually created by the client but set by the NATS server as 
      reply subject in responses to requests to the topic above.

Topics:
- Topic: `EMDEV.{deviceId}.REBOOT`  
  Makes for ACL entry: <kbd>sub</kbd> `EMDEV.{deviceId}.REBOOT`

## Topics

nats-manager listens on the topic `NATS-MANAGER.PROVISION-NEW-DEVICE` for provisioning requests for new devices. 
