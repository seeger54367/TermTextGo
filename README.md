# TermTextGo
A rewrite of [TermText](https://github.com/seeger54367/TermText) in the Go language

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
- Modify the info in main.go to add your Twilio number, and
  existing contacts if desired.
- In the backend.go file add paths to a file containing your Twilio
  Account ID and a file containing your account Auth token string. Both
  are found on the Twilio project dashboard.
- Run the following command to install:

```bash
  go install

```

# Keybindings
- 'k' for up in contact menu. 
- 'j' for down in contact menu. 
- 'ENTR' to send message to selected contact. 
- 'r' to reload messages for all contacts (may seem unresponsive for a couple seconds)
