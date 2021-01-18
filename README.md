# TermTextGo
A rewrite of TermText in the Go language

# Purpose
This app is a curses interface to send and receive SMS text messages
over the VoIP provider Twilio.

- Screenshot coming soon.
![termTextGo Screenshot](images/example1blurred.png)

# Setup

## Ground Work
Note that before this tool will work, you will need a Twilio project.
You will also need the numbers setup to send and receive SMS text
messages.


## After Twilio account setup
- Clone the repo.
- cd into repo
- Modify the info in interface.go to add your Twilio number, and
  existing contacts if desired.
- Run the following command to install:

```bash
  go install

```



