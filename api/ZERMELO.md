# Talking with Zermelo

## What does the webapp do?

### Appointment retrieval

`https://{institution}.zportal.nl/api/v3/liveschedule?student={student}&week={week}&fields=appointmentInstance,start,end,startTimeSlotName,endTimeSlotName,subjects,groups,locations,teachers,cancelled,changeDescription,schedulerRemark,content,appointmentType`  
<kbd>week</kbd> is formatted like so: `year` `week` (e.g. `202036`, where <kbd>year</kbd> = `2020`, <kbd>week</kbd> = `36`). Unsure if week is required to have a leading `0`.

The HTTP Header <kbd>[If-Modified-Since](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-Modified-Since)</kbd> is used to retrieve content that has been modified after earlier requests. The `liveschedule` endpoint is polled by the webapp.

## Datatypes

### Appointment

Properties:
* <kbd>actions</kbd>: <kbd>[][Action](#action)</kbd>
* <kbd>status</kbd>: <kbd>[][Status](#status)</kbd> :question:  
* <kbd>start</kbd>: <kbd>[time.Time]</kbd> (UNIX timestamp)
* <kbd>end</kbd>: <kbd>[time.Time]</kbd> (UNIX timestamp)
* <kbd>cancelled</kbd>: <kbd>bool</kbd>
* <kbd>appointmentInstance</kbd>: <kbd>int(64)</kbd>
* <kbd>startTimeSlotName</kbd>: <kbd>string</kbd>
* <kbd>endTimeSlotName</kbd>: <kbd>string</kbd>
* <kbd>subjects</kbd>: <kbd>[]string</kbd>
* <kbd>groups</kbd>: <kbd>[]string</kbd>
* <kbd>locations</kbd>: <kbd>[]string</kbd>
* <kbd>teachers</kbd>: <kbd>[]string</kbd>
* <kbd>changeDescription</kbd>: <kbd>string</kbd>
* <kbd>schedulerRemark</kbd>: <kbd>string</kbd>
* <kbd>content</kbd>: <kbd>string</kbd> :question: (null in example)
* <kbd>id</kbd>: <kbd>int(64)</kbd>

### Action

* <kbd>status</kbd>: <kbd>[][Status](#status)</kbd>
* <kbd>appointment</kbd>: <kbd>[][Appointment](#appointment)</kbd>
* <kbd>allowed</kbd>: <kbd>bool</kbd>
* <kbd>post</kbd>: <kbd>string</kbd> (relative URL)  
  In case of enrollment, can be used to enroll.

  Example: `/api/v3/liveschedule/enrollment?enroll=3313930&unenroll=`
  In this URL, `3313930` is <kbd>appointment.id</kbd>

### Status

Properties:
* <kbd>code</kbd>: <kbd>int</kbd>  
  Current status of the enrollment?

  Examples (en):  
  - `1002`: You will be unenrolled for all appointments at this time  
  - `2002`: Enrollment OK  
  - `4010`: This would cause a conflict in the schedule  

* <kbd>nl</kbd>: <kbd>string</kbd>
  Message localized in Dutch. See examples above.

* <kbd>en</kbd>: <kbd>string</kbd>
  Message localized in English. See examples above.

#### Notes

The properties <kbd>nl</kbd> and <kbd>en</kbd> may actually not be static, but ISO 639-1 language codes (think OpenAPI 3's `additionalProperties`).

[time.Time]: https://golang.org/pkg/time/#Time
