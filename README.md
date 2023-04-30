# tana-calendar-helper
A small service, written in Go, that helps Tana get your Calendar for the day. Provides an API you can call from Tana. 

Only works on Mac for now since it relies on Apple's Calendar.app and associated API to fetch the calendar data. If Calendar.app is configured to synchronize calendars from iCloud, Google and/or Microsoft Office, this will see all that data.

See also the Tana Template published at 

## Installation
Grab the latest release zip from github in the Releases area.

When you unzip the directory, you'll have three files 
```tana-calendar-helper
scripts/getcalendar.swift
scripts/calendar_auth.scpt
```

You launch the service on a command line like this:
`tana-calendar-helper`

The service will startup and begin listening for API calls on port 4000

You can also provide a different port number like this:
`tana-calendar-helper -port 4096`

Whenever you start the service, you ay get a permission box from Mac OS asking you to grant network access to the service. You need to Allow this.

You can test that the service is working by going to the following URL in your browser:
`http://localhost:4000/`  (or whatever port you are using)

If successful, you should get back a page of Usage instructions.

## Authorizing script to access Calendar.app
The first time you run this service and access the actual `/calendar` endpoint successfully, you will get another permissions dialog asking you to authorize access to your calendars. This is a one-time thing but if you don't allow it, you may be stuck and unable to retrieve data thereafter!

## Installing Command in Tana
Please see the Tana Template here:
![Tana Calendar Helper](https://app.tana.inc/?bundle=cVYW2gX8nY.WUDhKchZDK "Tana Calendar Helper config")

Or set up a command node like this:
![Getting Usage into Tana](assets/tana-helper-command.png?raw=true "Tana Command Node")

### Calendar API stuff

The Calendar API only works when run as a localhost service on a Mac since it relies on your Apple Calendar configuration to act as a "gateway" to your calendar services. This does allow it to reach iCloud, Google and Office365 calendars however.

The `/calendar` endpoint will by default return you a list of your meetings for today from a calendar named "Calendar".

You can change things with the following JSON payload. All fields are optional.

```
{
  "me": "self name", // your own name to avoid adding you as an attendee
  "meeting": "#tag", // the tag to use for meetings
  "person": "#tag", // tag for people / attendees
  "solo": true | false, // include meetings with just one person (yourself?)
  "one2one": "#tag", // tag for 1 to 1 meetings
  "calendar": "Calendar",
  "offset": -n | 0 | +n  // how many days before or after today to start from
  "range": >= 1 // how many daysto retrieve. Defaults to 1
}
```

See the [scripts/getcalendar.swift](scripts/getcalendar.swift) script for more details.